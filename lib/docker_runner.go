package lib

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./docker_runner_mock.go -source=docker_runner.go
type DockerRunner interface {
	ContainerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error)
	ContainerStart(ctx context.Context, logger *log.Logger, stage ManifestStage, needsNetwork bool) (err error)
	ContainerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string, start time.Time) (err error)
	ContainerGetExitCode(ctx context.Context, logger *log.Logger, containerID string) (exitCode int, err error)
	ContainerRemove(ctx context.Context, logger *log.Logger, containerID string) (err error)
	ContainerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error)
	NetworkCreate(ctx context.Context) (err error)
	NetworkRemove(ctx context.Context) (err error)
	NeedsNetwork(stages []*ManifestStage) bool
	StopRunningContainers(ctx context.Context) (err error)
}

type dockerRunner struct {
	commandRunner          CommandRunner
	buildDirectory         string
	pulledImages           map[string]struct{}
	pulledImagesMutex      *MapMutex
	runningContainers      map[string]ManifestStage
	runningContainersMutex *MapMutex
	networkName            string
}

func NewDockerRunner(commandRunner CommandRunner, randomStringGenerator RandomStringGenerator, buildDirectory string) DockerRunner {

	networkName := fmt.Sprintf("infinity-%v", randomStringGenerator.GenerateRandomString(10))

	return &dockerRunner{
		commandRunner:          commandRunner,
		buildDirectory:         buildDirectory,
		pulledImages:           make(map[string]struct{}),
		pulledImagesMutex:      NewMapMutex(),
		runningContainers:      make(map[string]ManifestStage),
		runningContainersMutex: NewMapMutex(),
		networkName:            networkName,
	}
}

func (b *dockerRunner) ContainerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {

	b.pulledImagesMutex.Lock(stage.Image)
	defer b.pulledImagesMutex.Unlock(stage.Image)

	if _, ok := b.pulledImages[stage.Image]; ok {
		logger.Printf(aurora.Gray(12, "Already pulled image %v").String(), aurora.BrightBlue(stage.Image))
		return
	}

	dockerCommand := "docker"
	dockerPullArgs := []string{
		"pull",
		stage.Image,
	}

	logger.Printf(aurora.Gray(12, "Pulling image %v").String(), aurora.BrightBlue(stage.Image))

	start := time.Now()
	err = b.commandRunner.RunCommand(ctx, logger, "", dockerCommand, dockerPullArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed pulling in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("pulling image %v for stage %v failed: %w", stage.Image, stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Pulled in %v").String(), aurora.BrightGreen(elapsed.String()))

	b.pulledImages[stage.Image] = struct{}{}

	return nil
}

func (b *dockerRunner) ContainerStart(ctx context.Context, logger *log.Logger, stage ManifestStage, needsNetwork bool) (err error) {

	pwd, err := filepath.Abs(b.buildDirectory)
	if err != nil {
		return
	}

	dockerCommand := "docker"
	dockerRunArgs := []string{
		"run",
	}

	if stage.Detach {
		dockerRunArgs = append(dockerRunArgs, "--detach", fmt.Sprintf("--name=%v", stage.Name))
	} else {
		dockerRunArgs = append(dockerRunArgs, "--rm")
	}

	if needsNetwork {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--network=%v", b.networkName))
	}

	if stage.MountWorkingDirectory != nil && *stage.MountWorkingDirectory {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--volume=%v:%v", pwd, stage.WorkingDirectory))
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--workdir=%v", stage.WorkingDirectory))
	}

	for _, v := range stage.Volumes {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--volume=%v", v))
	}
	for _, d := range stage.Devices {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--device=%v", d))
	}

	// add parameters to envvars
	env := stage.Env
	for k, v := range stage.Parameters {
		env[ToUpperSnakeCase("INFINITY_PARAMETER_"+k)] = fmt.Sprintf("%v", v)
	}

	// loop envvars in sorted order
	envKeys := make([]string, 0, len(env))
	for k := range env {
		envKeys = append(envKeys, k)
	}
	sort.Strings(envKeys)
	for _, k := range envKeys {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--env=%v=%v", k, env[k]))
	}

	if stage.Privileged {
		dockerRunArgs = append(dockerRunArgs, "--privileged")
	}
	if len(stage.Commands) > 0 {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--entrypoint=%v", "/bin/sh"))
	}

	dockerRunArgs = append(dockerRunArgs, stage.Image)

	if len(stage.Commands) > 0 {
		commandsArg := []string{}
		if len(stage.Commands) > 0 {
			commandsArg = append(commandsArg, "set -e")
		}
		for _, c := range stage.Commands {
			commandsArg = append(commandsArg, fmt.Sprintf(`printf "\033[38;5;244m> %v\033[0m\n"`, c))
			commandsArg = append(commandsArg, c)
		}
		dockerRunArgs = append(dockerRunArgs, []string{
			"-c",
			strings.Join(commandsArg, " ; "),
		}...)
	}

	if stage.Detach {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "Starting detached stage").String())
		}
	}

	if stage.Detach {
		start := time.Now()
		containerIDBytes, detachedErr := b.commandRunner.RunCommandWithOutput(ctx, logger, "", dockerCommand, dockerRunArgs)
		elapsed := time.Since(start)

		if detachedErr != nil {
			logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
			return detachedErr
		}

		containerID := strings.TrimSuffix(string(containerIDBytes), "\n")
		b.addRunningContainer(stage, containerID)
		logger.Printf(aurora.Gray(12, "Started in %v").String(), aurora.BrightGreen(elapsed.String()))

		return
	}

	logger.Printf(aurora.Gray(12, "Executing commands").String())

	start := time.Now()
	err = b.commandRunner.RunCommand(ctx, logger, "", dockerCommand, dockerRunArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return
	}

	logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))

	return
}

func (b *dockerRunner) ContainerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string, start time.Time) (err error) {

	// follow logs
	dockerCommand := "docker"
	dockerLogsArgs := []string{
		"logs",
		"--follow",
		containerID,
	}

	err = b.commandRunner.RunCommand(ctx, logger, "", dockerCommand, dockerLogsArgs)
	elapsed := time.Since(start)
	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("stage %v failed: %w", stage.Name, err)
	}

	return nil
}

func (b *dockerRunner) ContainerGetExitCode(ctx context.Context, logger *log.Logger, containerID string) (exitCode int, err error) {
	// check exit code
	dockerCommand := "docker"
	dockerInspectArgs := []string{
		"inspect",
		"--format='{{.State.ExitCode}}'",
		containerID,
	}

	var output []byte
	output, err = b.commandRunner.RunCommandWithOutput(ctx, logger, "", dockerCommand, dockerInspectArgs)
	if err != nil {
		return
	}

	output = bytes.Trim(output, "'\n")

	return strconv.Atoi(string(output))
}

func (b *dockerRunner) ContainerRemove(ctx context.Context, logger *log.Logger, containerID string) (err error) {
	// tail logs
	dockerCommand := "docker"
	dockerRemoveArgs := []string{
		"rm",
		"--volumes",
		containerID,
	}

	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), logger, "", dockerCommand, dockerRemoveArgs)

	return
}

func (b *dockerRunner) ContainerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error) {
	dockerCommand := "docker"
	dockerStopArgs := []string{
		"stop",
		"--time=30",
		containerID,
	}

	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), logger, "", dockerCommand, dockerStopArgs)

	return
}

func (b *dockerRunner) NetworkCreate(ctx context.Context) (err error) {
	dockerCommand := "docker"
	dockerNetworkCreateArgs := []string{
		"network",
		"create",
		b.networkName,
	}

	log.Printf(aurora.Gray(12, "Creating network %v").String(), aurora.BrightBlue(b.networkName))

	start := time.Now()
	_, err = b.commandRunner.RunCommandWithOutput(ctx, nil, "", dockerCommand, dockerNetworkCreateArgs)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf(aurora.Gray(12, "Failed in %v\n").String(), aurora.BrightRed(elapsed.String()))
		return err
	}
	log.Printf(aurora.Gray(12, "Completed in %v\n").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}

func (b *dockerRunner) NetworkRemove(ctx context.Context) (err error) {
	dockerCommand := "docker"
	dockerNetworkRemoveArgs := []string{
		"network",
		"rm",
		b.networkName,
	}

	log.Printf(aurora.Gray(12, "Removing network %v").String(), aurora.BrightBlue(b.networkName))

	start := time.Now()
	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), nil, "", dockerCommand, dockerNetworkRemoveArgs)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf(aurora.Gray(12, "Failed in %v\n").String(), aurora.BrightRed(elapsed.String()))
		return err
	}
	log.Printf(aurora.Gray(12, "Completed in %v\n").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}

func (b *dockerRunner) NeedsNetwork(stages []*ManifestStage) bool {
	for _, s := range stages {
		if s.Detach {
			return true
		}
		if b.NeedsNetwork(s.Stages) {
			return true
		}
	}

	return false
}

func (b *dockerRunner) StopRunningContainers(ctx context.Context) (err error) {

	if len(b.runningContainers) > 0 {
		log.Printf("Stopping %v running stage containers\n\n", len(b.runningContainers))

		semaphore := NewSemaphore(len(b.runningContainers))
		errorChannel := make(chan error, len(b.runningContainers))

		for containerID, stage := range b.runningContainers {
			semaphore.Acquire()
			go func(ctx context.Context, stage ManifestStage, containerID string) {
				defer semaphore.Release()

				logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

				defer func() {
					removeErr := b.ContainerRemove(ctx, logger, containerID)
					if err == nil {
						err = removeErr
					}
				}()

				stopErrorChannel := make(chan error)
				go func() {
					stopErrorChannel <- b.ContainerStop(ctx, logger, stage, containerID)
				}()

				start := time.Now()

				// tail logs
				err = b.ContainerLogs(ctx, logger, stage, containerID, start)
				if err != nil {
					errorChannel <- err
				}

				// check exit code
				exitCode, err := b.ContainerGetExitCode(ctx, logger, containerID)
				if err != nil {
					errorChannel <- err
				}

				if exitCode > 0 {
					errorChannel <- fmt.Errorf("stage %v failed with exit code %v", stage.Name, exitCode)
				}

				// wait until stop finishes
				err = <-stopErrorChannel
				if err != nil {
					errorChannel <- err
				}
			}(ctx, stage, containerID)
		}

		semaphore.Wait()

		log.Println("")

		close(errorChannel)
		for err = range errorChannel {
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (b *dockerRunner) addRunningContainer(stage ManifestStage, containerID string) {
	// add container id to running containers map
	b.runningContainersMutex.Lock(stage.Name)
	defer b.runningContainersMutex.Unlock(stage.Name)
	b.runningContainers[containerID] = stage
}
