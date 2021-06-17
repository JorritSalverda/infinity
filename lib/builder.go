package lib

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Builder interface {
	Validate(ctx context.Context) (err error)
	Build(ctx context.Context) (err error)
}

type builder struct {
	manifestPath string
}

func NewBuilder() Builder {
	return &builder{
		manifestPath: ".infinity.yaml",
	}
}

func (b *builder) Validate(ctx context.Context) (err error) {
	log.Printf("Validating manifest %v...\n", b.manifestPath)

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
	log.Printf("Building manifest %v...\n", b.manifestPath)

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
	if _, err = os.Stat(b.manifestPath); os.IsNotExist(err) {
		return manifest, fmt.Errorf("Manifest %v does not exist, cannot continue", b.manifestPath)
	}

	// read manifest
	manifestBytes, err := ioutil.ReadFile(".infinity.yaml")
	if err != nil {
		return
	}

	// unmarshal bytes into manifest
	if err = yaml.UnmarshalStrict(manifestBytes, &manifest); err != nil {
		return manifest, fmt.Errorf("Manifest %v is not valid: %w", b.manifestPath, err)
	}

	return
}

func (b *builder) runManifest(ctx context.Context, manifest Manifest) (err error) {

	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	for _, stage := range manifest.Build.Stages {
		log.Printf("Running stage %v...\n", stage.Name)
		// docker run <image> <commands>
		commandsArg := strings.Join(stage.Commands, " ; ")

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
			fmt.Sprintf("set -e ; %v", commandsArg),
		}...)

		log.Printf("> %v %v\n", dockerCommand, strings.Join(dockerRunArgs, " "))
		err = b.runCommand(ctx, dockerCommand, dockerRunArgs)
		if err != nil {
			return fmt.Errorf("Stage %v failed: %w", stage.Name, err)
		}
	}

	return nil
}

func (b *builder) runCommand(ctx context.Context, command string, args []string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}
