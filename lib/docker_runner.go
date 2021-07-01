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

	"github.com/logrusorgru/aurora"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -package=lib -destination ./docker_runner_mock.go -source=docker_runner.go
type DockerRunner interface {
	ContainerImageIsPulled(ctx context.Context, logger *log.Logger, stage ManifestStage) (isPulled bool, err error)
	ContainerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error)
	ContainerStart(ctx context.Context, logger *log.Logger, stage ManifestStage, needsNetwork bool) (err error)
	ContainerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error)
	ContainerGetExitCode(ctx context.Context, logger *log.Logger, containerID string) (exitCode int, err error)
	ContainerWait(ctx context.Context, logger *log.Logger, containerID string) (err error)
	ContainerRemove(ctx context.Context, logger *log.Logger, containerID string) (err error)
	ContainerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string, timeoutSeconds int) (err error)
	NetworkCreate(ctx context.Context, logger *log.Logger) (err error)
	NetworkRemove(ctx context.Context, logger *log.Logger) (err error)
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

func (b *dockerRunner) ContainerImageIsPulled(ctx context.Context, logger *log.Logger, stage ManifestStage) (isPulled bool, err error) {

	b.pulledImagesMutex.Lock(stage.Image)
	defer b.pulledImagesMutex.Unlock(stage.Image)

	if _, ok := b.pulledImages[stage.Image]; ok {
		logger.Printf(aurora.Gray(12, "Already pulled image %v").String(), aurora.BrightBlue(stage.Image))
		return true, nil
	}

	dockerCommand := "docker"
	dockerPullArgs := []string{
		"images",
		"--format='{{.Repository}}:{{.Tag}}'",
		stage.Image,
	}

	output, err := b.commandRunner.RunCommandWithOutput(ctx, logger, "", dockerCommand, dockerPullArgs)
	if err != nil {
		return false, err
	}

	output = bytes.Trim(output, "'\n")

	if string(output) == stage.Image {
		return true, nil
	}

	return false, nil
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

	err = b.commandRunner.RunCommand(ctx, logger, "", dockerCommand, dockerPullArgs)
	if err != nil {
		return fmt.Errorf("pulling image %v for stage %v failed: %w", stage.Image, stage.Name, err)
	}

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
		"--detach",
	}

	if stage.Background {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--name=%v", stage.Name))
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
		commandsArg := []string{"set -e"}
		for _, c := range stage.Commands {
			commandsArg = append(commandsArg, fmt.Sprintf(`printf '\033[38;5;244m> %%s\033[0m\n' '%v'`, c))
			commandsArg = append(commandsArg, c)
		}
		dockerRunArgs = append(dockerRunArgs, []string{
			"-c",
			strings.Join(commandsArg, " ; "),
		}...)
	}

	if stage.Background {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "Starting stage in background").String())
		}
	} else {
		logger.Printf(aurora.Gray(12, "Starting stage").String())
	}

	containerIDBytes, err := b.commandRunner.RunCommandWithOutput(context.Background(), logger, "", dockerCommand, dockerRunArgs)
	if err != nil {
		return
	}

	containerID := strings.TrimSuffix(string(containerIDBytes), "\n")
	b.addRunningContainer(stage, containerID)

	if stage.Background {
		return
	}

	// ensure container gets removed at the end
	defer func() {
		waitErr := b.ContainerWait(ctx, logger, containerID)
		if err == nil {
			err = waitErr
		}
		removeErr := b.ContainerRemove(ctx, logger, containerID)
		if err == nil {
			err = removeErr
		}
		b.removeRunningContainer(stage, containerID)
	}()

	// stop container on cancellation
	waitDone := make(chan struct{})
	defer close(waitDone)
	go func() {
		select {
		case <-ctx.Done():
			_ = b.ContainerStop(ctx, logger, stage, containerID, 5)
		case <-waitDone:
		}
	}()

	// tail logs
	err = b.ContainerLogs(ctx, logger, stage, containerID)
	if err != nil {
		return
	}

	// check exit code
	exitCode, err := b.ContainerGetExitCode(ctx, logger, containerID)
	if err != nil {
		return
	}
	if exitCode > 0 {
		return fmt.Errorf("stage %v failed with exit code %v", stage.Name, exitCode)
	}

	return
}

func (b *dockerRunner) ContainerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error) {

	// follow logs
	dockerCommand := "docker"
	dockerLogsArgs := []string{
		"logs",
		"--follow",
		containerID,
	}

	err = b.commandRunner.RunCommand(ctx, logger, "", dockerCommand, dockerLogsArgs)
	if err != nil {
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

func (b *dockerRunner) ContainerWait(ctx context.Context, logger *log.Logger, containerID string) (err error) {

	// tail logs
	dockerCommand := "docker"
	dockerWaitArgs := []string{
		"wait",
		containerID,
	}

	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), logger, "", dockerCommand, dockerWaitArgs)

	return
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

func (b *dockerRunner) ContainerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string, timeoutSeconds int) (err error) {

	dockerCommand := "docker"
	dockerStopArgs := []string{
		"stop",
		fmt.Sprintf("--time=%v", timeoutSeconds),
		containerID,
	}

	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), logger, "", dockerCommand, dockerStopArgs)

	return
}

func (b *dockerRunner) NetworkCreate(ctx context.Context, logger *log.Logger) (err error) {
	dockerCommand := "docker"
	dockerNetworkCreateArgs := []string{
		"network",
		"create",
		b.networkName,
	}

	logger.Printf(aurora.Gray(12, "Creating network %v").String(), aurora.BrightBlue(b.networkName))

	_, err = b.commandRunner.RunCommandWithOutput(ctx, nil, "", dockerCommand, dockerNetworkCreateArgs)
	if err != nil {
		return err
	}

	return nil
}

func (b *dockerRunner) NetworkRemove(ctx context.Context, logger *log.Logger) (err error) {
	dockerCommand := "docker"
	dockerNetworkRemoveArgs := []string{
		"network",
		"rm",
		b.networkName,
	}

	logger.Printf(aurora.Gray(12, "Removing network %v").String(), aurora.BrightBlue(b.networkName))

	_, err = b.commandRunner.RunCommandWithOutput(context.Background(), nil, "", dockerCommand, dockerNetworkRemoveArgs)
	if err != nil {
		return err
	}

	return nil
}

func (b *dockerRunner) NeedsNetwork(stages []*ManifestStage) bool {
	for _, s := range stages {
		if s.Background {
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

		g, ctx := errgroup.WithContext(ctx)
		for containerID, stage := range b.runningContainers {
			stage := stage
			containerID := containerID
			g.Go(func() (err error) {
				logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

				// ensure container gets removed at the end
				defer func() {
					waitErr := b.ContainerWait(ctx, logger, containerID)
					if err == nil {
						err = waitErr
					}
					removeErr := b.ContainerRemove(ctx, logger, containerID)
					if err == nil {
						err = removeErr
					}
					b.removeRunningContainer(stage, containerID)
				}()

				stopErrorChannel := make(chan error)
				go func() {
					stopErrorChannel <- b.ContainerStop(ctx, logger, stage, containerID, 30)
				}()

				// tail logs
				err = b.ContainerLogs(ctx, logger, stage, containerID)
				if err != nil {
					return err
				}

				// check exit code
				exitCode, err := b.ContainerGetExitCode(ctx, logger, containerID)
				if err != nil {
					return err
				}

				if exitCode > 0 {
					return fmt.Errorf("stage %v failed with exit code %v", stage.Name, exitCode)
				}

				// wait until stop finishes
				err = <-stopErrorChannel
				if err != nil {
					return err
				}

				return nil
			})
		}

		// wait for all containers to be stopped
		if err = g.Wait(); err != nil {
			return err
		}

		log.Println("")
	}

	return nil
}

func (b *dockerRunner) addRunningContainer(stage ManifestStage, containerID string) {
	// add container id to running containers map
	b.runningContainersMutex.Lock(stage.Name)
	defer b.runningContainersMutex.Unlock(stage.Name)
	b.runningContainers[containerID] = stage
}

func (b *dockerRunner) removeRunningContainer(stage ManifestStage, containerID string) {
	// remove container from map
	b.runningContainersMutex.Lock(stage.Name)
	defer b.runningContainersMutex.Unlock(stage.Name)
	delete(b.runningContainers, containerID)
}
