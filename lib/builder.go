package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -package=lib -destination ./builder_mock.go -source=builder.go
type Builder interface {
	Validate(ctx context.Context) (manifest Manifest, err error)
	Build(ctx context.Context) (err error)
}

type builder struct {
	manifestReader        ManifestReader
	dockerRunner          DockerRunner
	hostRunner            HostRunner
	forcePull             bool
	buildDirectory        string
	buildManifestFilename string
}

func NewBuilder(manifestReader ManifestReader, dockerRunner DockerRunner, hostRunner HostRunner, forcePull bool, buildDirectory, buildManifestFilename string) Builder {

	return &builder{
		manifestReader:        manifestReader,
		dockerRunner:          dockerRunner,
		hostRunner:            hostRunner,
		forcePull:             forcePull,
		buildDirectory:        buildDirectory,
		buildManifestFilename: buildManifestFilename,
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

	needsNetwork := b.dockerRunner.NeedsNetwork(manifest.Build.Stages)

	if needsNetwork {
		err = b.dockerRunner.NetworkCreate(ctx)
		if err != nil {
			return
		}

		defer func() {
			terminateErr := b.dockerRunner.StopRunningContainers(ctx)
			if err == nil {
				err = terminateErr
			}

			if needsNetwork {
				terminateErr = b.dockerRunner.NetworkRemove(ctx)
				if err == nil {
					err = terminateErr
				}
			}
		}()
	}

	for _, stage := range manifest.Build.Stages {
		err = b.runStage(ctx, *stage, needsNetwork)
		log.Println("")
		if err != nil {
			return
		}
	}

	return nil
}

func (b *builder) runStage(ctx context.Context, stage ManifestStage, needsNetwork bool, prefixes ...string) (err error) {

	prefixes = append(prefixes, stage.Name)
	prefix := strings.Join(prefixes, "] [")

	logger := log.New(os.Stdout, aurora.Gray(12, fmt.Sprintf("[%v] ", prefix)).String(), 0)

	if len(stage.Stages) > 0 {
		return b.runParallelStages(ctx, stage, needsNetwork)
	}

	switch stage.RunnerType {
	case RunnerTypeContainer:
		var isPulled bool
		if !b.forcePull {
			isPulled, err = b.dockerRunner.ContainerImageIsPulled(ctx, logger, stage)
			if err != nil {
				return
			}
		}

		if !isPulled {
			err = b.dockerRunner.ContainerPull(ctx, logger, stage)
			if err != nil {
				return
			}
		}

		err = b.dockerRunner.ContainerStart(ctx, logger, stage, needsNetwork)
		if err != nil {
			return
		}

		return nil

	case RunnerTypeHost:
		return b.hostRunner.RunStage(ctx, logger, stage)
	}

	return fmt.Errorf("runner %v is not supported", stage.RunnerType)
}

func (b *builder) runParallelStages(ctx context.Context, stage ManifestStage, needsNetwork bool) (err error) {
	g, ctx := errgroup.WithContext(ctx)
	for _, s := range stage.Stages {
		s := s
		g.Go(func() error { return b.runStage(ctx, *s, needsNetwork, stage.Name) })
	}

	return g.Wait()
}
