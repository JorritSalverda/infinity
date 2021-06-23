package lib

import (
	"context"
	"fmt"
	"log"
	"os"
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
	manifestReader          ManifestReader
	commandRunner           CommandRunner
	verbose                 bool
	buildManifestFilename   string
	pulledImages            map[string]struct{}
	pulledImagesMutex       *MapMutex
	detachedContainers      map[*ManifestStage]string
	detachedContainersMutex *MapMutex
}

func NewBuilder(manifestReader ManifestReader, commandRunner CommandRunner, verbose bool, buildManifestFilename string) Builder {
	return &builder{
		manifestReader:          manifestReader,
		commandRunner:           commandRunner,
		buildManifestFilename:   buildManifestFilename,
		verbose:                 verbose,
		pulledImages:            make(map[string]struct{}),
		pulledImagesMutex:       NewMapMutex(),
		detachedContainers:      make(map[*ManifestStage]string),
		detachedContainersMutex: NewMapMutex(),
	}
}

func (b *builder) Validate(ctx context.Context) (manifest Manifest, err error) {
	log.Printf("Validating manifest %v", aurora.BrightBlue(b.buildManifestFilename))

	manifest, err = b.manifestReader.GetManifest(ctx, b.buildManifestFilename)
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
		terminateErr := b.terminateDetachedContainers(ctx)
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
		return b.containerRun(ctx, logger, stage)
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

func (b *builder) terminateDetachedContainers(ctx context.Context) (err error) {

	if len(b.detachedContainers) > 0 {
		log.Printf("Terminating %v detached stage containers\n\n", len(b.detachedContainers))

		semaphore := NewSemaphore(len(b.detachedContainers))
		errorChannel := make(chan error, len(b.detachedContainers))

		for stage, containerID := range b.detachedContainers {
			semaphore.Acquire()
			go func(ctx context.Context, stage ManifestStage, containerID string) {
				defer semaphore.Release()
				errorChannel <- b.terminateDetachedContainer(ctx, stage, containerID)
			}(ctx, *stage, containerID)
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

func (b *builder) terminateDetachedContainer(ctx context.Context, stage ManifestStage, containerID string) (err error) {

	logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

	errorChannel := make(chan error)

	go func() {
		// stop container
		dockerCommand := "docker"
		dockerStopArgs := []string{
			"stop",
			"--time=30",
			containerID,
		}
		if b.verbose {
			logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerStopArgs, " "))
		}
		errorChannel <- b.commandRunner.RunCommandWithLogger(ctx, logger, dockerCommand, dockerStopArgs)
	}()

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

	err = b.commandRunner.RunCommandWithLogger(ctx, logger, dockerCommand, dockerLogsArgs)
	if err != nil {
		return
	}

	return <-errorChannel
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
	err = b.commandRunner.RunCommandWithLogger(ctx, logger, dockerCommand, dockerPullArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed pulling in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("pulling image %v for stage %v failed: %w", stage.Image, stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Pulled in %v").String(), aurora.BrightGreen(elapsed.String()))

	b.pulledImages[stage.Image] = struct{}{}

	return nil
}

func (b *builder) containerRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	dockerCommand := "docker"
	dockerRunArgs := []string{
		"run",
		"--rm",
		fmt.Sprintf("--volume=%v:/work", pwd),
		"--workdir=/work",
	}
	for _, m := range stage.Mounts {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--volume=%v", m))
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
	if stage.Detach {
		dockerRunArgs = append(dockerRunArgs, "--detach")
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

		start := time.Now()
		containerIDBytes, err := b.commandRunner.RunCommandWithOutput(ctx, dockerCommand, dockerRunArgs)
		elapsed := time.Since(start)
		if err != nil {
			logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
			return err
		}
		logger.Printf(aurora.Gray(12, "Started in %v").String(), aurora.BrightGreen(elapsed.String()))

		containerID := string(containerIDBytes)
		containerID = strings.TrimSuffix(containerID, "\n")

		b.detachedContainersMutex.Lock(stage.Name)
		defer b.detachedContainersMutex.Unlock(stage.Name)

		b.detachedContainers[&stage] = containerID

		return nil
	}

	logger.Printf(aurora.Gray(12, "Executing commands").String())

	start := time.Now()
	err = b.commandRunner.RunCommandWithLogger(ctx, logger, dockerCommand, dockerRunArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("stage %v failed: %w", stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}

func (b *builder) metalRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	logger.Printf(aurora.Gray(12, "Executing commands in bare metal mode").String())

	start := time.Now()

	for _, c := range stage.Commands {
		logger.Printf(aurora.Gray(12, "> %v").String(), c)

		splitCommands := strings.Split(c, " ")
		err = b.commandRunner.RunCommandWithLogger(ctx, logger, splitCommands[0], splitCommands[1:])
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
