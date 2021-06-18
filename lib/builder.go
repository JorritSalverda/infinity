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
	pulledImages          map[string]struct{}
}

func NewBuilder(verbose bool, buildManifestFilename string) Builder {
	return &builder{
		buildManifestFilename: buildManifestFilename,
		verbose:               verbose,
		pulledImages:          make(map[string]struct{}),
	}
}

func (b *builder) Validate(ctx context.Context) (err error) {
	log.Printf("Validating manifest %v\n", aurora.BrightBlue(b.buildManifestFilename))

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
	log.Printf("Building manifest %v\n", aurora.BrightBlue(b.buildManifestFilename))

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
	logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", stage.Name)).String(), 0)

	// docker pull <image>
	err = b.dockerPull(ctx, logger, stage)
	if err != nil {
		return
	}

	// docker run <image> <commands>
	err = b.dockerRun(ctx, logger, stage)
	if err != nil {
		return
	}

	return nil
}

func (b *builder) dockerPull(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {

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
	err = b.runCommandWithLogger(ctx, logger, dockerCommand, dockerPullArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed pulling in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("Pulling image %v for stage %v failed: %w", stage.Image, stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Pulled in %v").String(), aurora.BrightGreen(elapsed.String()))

	b.pulledImages[stage.Image] = struct{}{}

	return nil
}

func (b *builder) dockerRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	commandsArg := []string{
		"set -e",
	}
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
		strings.Join(commandsArg, " ; "),
	}...)

	if b.verbose {
		logger.Printf(aurora.Gray(12, "> %v %v").String(), dockerCommand, strings.Join(dockerRunArgs, " "))
	}

	logger.Printf(aurora.Gray(12, "Executing commands").String())

	start := time.Now()
	err = b.runCommandWithLogger(ctx, logger, dockerCommand, dockerRunArgs)
	elapsed := time.Since(start)

	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("Stage %v failed: %w", stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))

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
