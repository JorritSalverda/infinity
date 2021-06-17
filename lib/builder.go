package lib

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

type Builder interface {
	Validate(ctx context.Context) (err error)
	Build(ctx context.Context) (err error)
}

type builder struct {
	verbose               bool
	buildManifestFilename string
}

func NewBuilder(verbose bool, buildManifestFilename string) Builder {
	return &builder{
		buildManifestFilename: buildManifestFilename,
		verbose:               verbose,
	}
}

func (b *builder) Validate(ctx context.Context) (err error) {
	log.Printf("Validating manifest %v...\n", b.buildManifestFilename)

	manifest, err := b.getManifest(ctx)
	if err != nil {
		return
	}

	err = manifest.Validate()
	if err != nil {
		return
	}

	log.Println("Manifest is valid!")

	return nil
}

func (b *builder) Build(ctx context.Context) (err error) {
	log.Printf("Building manifest %v...\n", b.buildManifestFilename)

	manifest, err := b.getManifest(ctx)
	if err != nil {
		return
	}

	err = manifest.Validate()
	if err != nil {
		return
	}

	err = b.runManifest(ctx, manifest)
	if err != nil {
		return
	}

	log.Println("Build succeeded!")

	return nil
}

func (b *builder) getManifest(ctx context.Context) (manifest Manifest, err error) {
	// check if manifest exists
	if _, err = os.Stat(b.buildManifestFilename); os.IsNotExist(err) {
		return manifest, fmt.Errorf("Manifest %v does not exist, cannot continue", b.buildManifestFilename)
	}

	// read manifest
	manifestBytes, err := ioutil.ReadFile(b.buildManifestFilename)
	if err != nil {
		return
	}

	// unmarshal bytes into manifest
	if err = yaml.UnmarshalStrict(manifestBytes, &manifest); err != nil {
		return manifest, fmt.Errorf("Manifest %v is not valid: %w", b.buildManifestFilename, err)
	}

	return
}

func (b *builder) runManifest(ctx context.Context, manifest Manifest) (err error) {
	log.Println("")
	for _, stage := range manifest.Build.Stages {
		err = b.runStage(ctx, stage)
		log.Println("")
		if err != nil {
			return
		}
	}

	return nil
}

func (b *builder) runStage(ctx context.Context, stage ManifestStage) (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	stageLogger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

	// docker run <image> <commands>

	commandsArg := []string{}
	for _, c := range stage.Commands {
		commandsArg = append(commandsArg, fmt.Sprintf(`echo -e "\x1b[38;5;244m> %v\x1b[0m"`, c))
		commandsArg = append(commandsArg, c)
	}

	dockerCommand := "docker"
	dockerRunArgs := []string{
		"run",
		"--rm",
		fmt.Sprintf("--volume=%v:/work", pwd),
		"--workdir=/work",
		"--entrypoint=/bin/sh",
	}
	for _, m := range stage.Mounts {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--volume=%v", m))
	}
	for k, v := range stage.Env {
		dockerRunArgs = append(dockerRunArgs, fmt.Sprintf("--env=%v=%v", k, v))
	}
	if stage.Privileged {
		dockerRunArgs = append(dockerRunArgs, "--privileged")
	}
	dockerRunArgs = append(dockerRunArgs, []string{
		stage.Image,
		"-c",
		fmt.Sprintf("set -e ; %v", strings.Join(commandsArg, " ; ")),
	}...)

	if b.verbose {
		stageLogger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerRunArgs, " "))
	}

	start := time.Now()
	err = b.runCommandWithLogger(ctx, stageLogger, dockerCommand, dockerRunArgs)
	elapsed := time.Since(start)

	if err != nil {
		stageLogger.Printf(aurora.Gray(12, "failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("Stage %v failed: %w", stage.Name, err)
	}
	stageLogger.Printf(aurora.Gray(12, "completed in %v").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}

func (b *builder) runCommandWithLogger(ctx context.Context, logger *log.Logger, command string, args []string) (err error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// start the command after having set up the pipe
	if err = cmd.Start(); err != nil {
		return
	}

	multi := io.MultiReader(stdout, stderr)

	// read command's output line by line
	in := bufio.NewScanner(multi)

	for in.Scan() {
		logger.Printf(in.Text())
	}

	if err = in.Err(); err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		return
	}

	return nil
}
