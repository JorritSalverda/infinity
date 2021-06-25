package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./builder_mock.go -source=builder.go
type Builder interface {
	Validate(ctx context.Context) (manifest Manifest, err error)
	Build(ctx context.Context) (err error)
}

type builder struct {
	manifestReader         ManifestReader
	commandRunner          CommandRunner
	verbose                bool
	buildDirectory         string
	buildManifestFilename  string
	pulledImages           map[string]struct{}
	pulledImagesMutex      *MapMutex
	runningContainers      map[string]ManifestStage
	runningContainersMutex *MapMutex
}

func NewBuilder(manifestReader ManifestReader, commandRunner CommandRunner, verbose bool, buildDirectory, buildManifestFilename string) Builder {
	return &builder{
		manifestReader:         manifestReader,
		commandRunner:          commandRunner,
		buildDirectory:         buildDirectory,
		buildManifestFilename:  buildManifestFilename,
		verbose:                verbose,
		pulledImages:           make(map[string]struct{}),
		pulledImagesMutex:      NewMapMutex(),
		runningContainers:      make(map[string]ManifestStage),
		runningContainersMutex: NewMapMutex(),
	}
}

func (b *builder) Validate(ctx context.Context) (manifest Manifest, err error) {
	log.Printf("Validating manifest %v", aurora.BrightBlue(b.buildManifestFilename))

	manifestPath := filepath.Join(b.buildDirectory, b.buildManifestFilename)

	manifest, err = b.manifestReader.GetManifest(ctx, manifestPath)
	if err != nil {
		return
	}

	warnings, errors := manifest.Validate()
	if len(warnings) > 0 {
		log.Println(aurora.BrightYellow("Manifest has warnings:"))
		for _, w := range warnings {
			log.Println(aurora.BrightYellow(w))
		}
	}
	if len(errors) > 0 {
		log.Println(aurora.BrightRed("Manifest has errors:"))
		for _, e := range errors {
			log.Println(aurora.BrightRed(e))
		}
		return manifest, fmt.Errorf("manifest failed validation")
	}

	log.Println("Manifest is valid!")

	return
}

func (b *builder) Build(ctx context.Context) (err error) {
	manifest, err := b.Validate(ctx)
	if err != nil {
		return
	}

	log.Printf("Building manifest %v", aurora.BrightBlue(b.buildManifestFilename))

	start := time.Now()
	err = b.runManifest(ctx, manifest)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("Build failed in %v\n", aurora.BrightRed(elapsed.String()))
		return
	}

	log.Printf("Build succeeded %v\n", aurora.BrightGreen(elapsed.String()))

	return nil
}

func (b *builder) runManifest(ctx context.Context, manifest Manifest) (err error) {
	log.Println("")

	defer func() {
		terminateErr := b.stopRunningContainers(ctx)
		if err == nil {
			err = terminateErr
		}
	}()

	for _, stage := range manifest.Build.Stages {
		err = b.runStage(ctx, *stage)
		log.Println("")
		if err != nil {
			return
		}
	}

	return nil
}

func (b *builder) runStage(ctx context.Context, stage ManifestStage) (err error) {
	logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

	if len(stage.Stages) > 0 {
		return b.runParallelStages(ctx, stage)
	}

	switch stage.RunnerType {
	case RunnerTypeContainer:
		// docker pull <image>
		err = b.containerPull(ctx, logger, stage)
		if err != nil {
			return
		}

		// docker run <image> <commands>
		var containerID string
		containerID, err = b.containerStart(ctx, logger, stage)
		if err != nil {
			return
		}

		if stage.Detach {
			// tailing logs happens when all stages are done
			return nil
		}

		return b.containerLogs(ctx, logger, stage, containerID)

	case RunnerTypeMetal:
		return b.metalRun(ctx, logger, stage)
	}

	return fmt.Errorf("runner %v is not supported", stage.RunnerType)
}

func (b *builder) runParallelStages(ctx context.Context, stage ManifestStage) (err error) {
	semaphore := NewSemaphore(len(stage.Stages))
	errorChannel := make(chan error, len(stage.Stages))

	for _, s := range stage.Stages {
		semaphore.Acquire()
		go func(ctx context.Context, s ManifestStage) {
			defer semaphore.Release()
			errorChannel <- b.runStage(ctx, s)
		}(ctx, *s)
	}

	semaphore.Wait()

	close(errorChannel)
	for err = range errorChannel {
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *builder) stopRunningContainers(ctx context.Context) (err error) {

	if len(b.runningContainers) > 0 {
		log.Printf("Terminating %v running stage containers\n\n", len(b.runningContainers))

		semaphore := NewSemaphore(len(b.runningContainers))
		errorChannel := make(chan error, len(b.runningContainers))

		for containerID, stage := range b.runningContainers {
			semaphore.Acquire()
			go func(ctx context.Context, stage ManifestStage, containerID string) {
				defer semaphore.Release()

				logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

				stopErrorChannel := make(chan error)
				go func() {
					stopErrorChannel <- b.containerStop(ctx, logger, stage, containerID)
				}()

				err = b.containerLogs(ctx, logger, stage, containerID)
				if err != nil {
					errorChannel <- err
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

func (b *builder) containerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {

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
	if b.verbose {
		logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerPullArgs, " "))
	}

	logger.Printf(aurora.Gray(12, "Pulling image %v").String(), aurora.BrightBlue(stage.Image))

	start := time.Now()
	err = b.commandRunner.RunCommandWithLogger(ctx, logger, "", dockerCommand, dockerPullArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed pulling in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("pulling image %v for stage %v failed: %w", stage.Image, stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Pulled in %v").String(), aurora.BrightGreen(elapsed.String()))

	b.pulledImages[stage.Image] = struct{}{}

	return nil
}

func (b *builder) containerStart(ctx context.Context, logger *log.Logger, stage ManifestStage) (containerID string, err error) {

	pwd, err := filepath.Abs(b.buildDirectory)
	if err != nil {
		return
	}

	dockerCommand := "docker"
	dockerRunArgs := []string{
		"run",
		"--rm",
		"--detach",
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

	if b.verbose {
		logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerRunArgs, " "))
	}

	if stage.Detach {
		logger.Printf(aurora.Gray(12, "Starting detached stage").String())
	}

	start := time.Now()
	containerIDBytes, err := b.commandRunner.RunCommandWithOutput(ctx, "", dockerCommand, dockerRunArgs)
	elapsed := time.Since(start)
	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return
	}

	containerID = strings.TrimSuffix(string(containerIDBytes), "\n")

	b.addRunningContainer(stage, containerID)

	if stage.Detach {
		logger.Printf(aurora.Gray(12, "Started in %v").String(), aurora.BrightGreen(elapsed.String()))
	}

	return
}

func (b *builder) containerLogs(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error) {

	// tail logs
	dockerCommand := "docker"
	dockerLogsArgs := []string{
		"logs",
		"--follow",
		containerID,
	}
	if b.verbose {
		logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerLogsArgs, " "))
	}

	logger.Printf(aurora.Gray(12, "Executing commands").String())

	start := time.Now()
	err = b.commandRunner.RunCommandWithLogger(ctx, logger, "", dockerCommand, dockerLogsArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("stage %v failed: %w", stage.Name, err)
	}

	if !stage.Detach {
		logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))
	}

	b.removeRunningContainer(stage, containerID)

	return nil
}

func (b *builder) containerStop(ctx context.Context, logger *log.Logger, stage ManifestStage, containerID string) (err error) {
	dockerCommand := "docker"
	dockerStopArgs := []string{
		"stop",
		"--time=30",
		containerID,
	}
	if b.verbose {
		logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerStopArgs, " "))
	}

	return b.commandRunner.RunCommandWithLogger(context.Background(), logger, "", dockerCommand, dockerStopArgs)
}

func (b *builder) addRunningContainer(stage ManifestStage, containerID string) {
	// add container id to running containers map
	b.runningContainersMutex.Lock(stage.Name)
	defer b.runningContainersMutex.Unlock(stage.Name)
	b.runningContainers[containerID] = stage
}

func (b *builder) removeRunningContainer(stage ManifestStage, containerID string) {
	// remove container from map
	b.runningContainersMutex.Lock(stage.Name)
	defer b.runningContainersMutex.Unlock(stage.Name)
	delete(b.runningContainers, containerID)
}

func (b *builder) metalRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	logger.Printf(aurora.Gray(12, "Executing commands in bare metal mode").String())

	start := time.Now()

	for _, c := range stage.Commands {
		logger.Printf(aurora.Gray(12, "> %v").String(), c)

		splitCommands := strings.Split(c, " ")
		err = b.commandRunner.RunCommandWithLogger(ctx, logger, b.buildDirectory, splitCommands[0], splitCommands[1:])
		if err != nil {
			break
		}
	}

	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("stage %v failed: %w", stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}
