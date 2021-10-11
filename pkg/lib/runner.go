package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -package=lib -destination ./runner_mock.go -source=runner.go
type Runner interface {
	Validate(ctx context.Context) (manifest Manifest, err error)
	Run(ctx context.Context, target string) (err error)
}

type runner struct {
	manifestReader        ManifestReader
	dockerRunner          DockerRunner
	hostRunner            HostRunner
	forcePull             bool
	buildDirectory        string
	buildManifestFilename string
}

func NewRunner(manifestReader ManifestReader, dockerRunner DockerRunner, hostRunner HostRunner, forcePull bool, buildDirectory, buildManifestFilename string) Runner {
	return &runner{
		manifestReader:        manifestReader,
		dockerRunner:          dockerRunner,
		hostRunner:            hostRunner,
		forcePull:             forcePull,
		buildDirectory:        buildDirectory,
		buildManifestFilename: buildManifestFilename,
	}
}

func (b *runner) Validate(ctx context.Context) (manifest Manifest, err error) {
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

func (b *runner) Run(ctx context.Context, target string) (err error) {
	manifest, err := b.Validate(ctx)
	if err != nil {
		return
	}

	log.Printf("Running manifest %v target %v", aurora.BrightBlue(b.buildManifestFilename), aurora.BrightBlue(target))

	if err = b.handleFunc(ctx, nil, func() error {
		return b.runManifest(ctx, manifest, target)
	}); err != nil {
		if errors.Is(err, ErrCanceled) {
			return nil
		}
		return
	}

	return nil
}

func (b *runner) getManifestTarget(ctx context.Context, manifest Manifest, target string) (manifestTarget *ManifestTarget, err error) {
	for _, t := range manifest.Targets {
		if t.Name == target {
			return t, nil
		}
	}

	return nil, fmt.Errorf("Target %v is not defined in manifest", target)
}

func (b *runner) runManifest(ctx context.Context, manifest Manifest, target string) (err error) {
	log.Println("")

	// get target
	manifestTarget, err := b.getManifestTarget(ctx, manifest, target)
	if err != nil {
		return
	}

	// set color codes for coloring stage logs
	b.setColorCode(manifestTarget.Stages)

	needsNetwork := b.dockerRunner.NeedsNetwork(manifestTarget.Stages)

	if needsNetwork {
		logger := log.New(os.Stdout, aurora.Gray(12, "[infinity] ").String(), 0)
		if err = b.handleFunc(ctx, logger, func() error {
			return b.dockerRunner.NetworkCreate(ctx, logger)
		}); err != nil {
			if errors.Is(err, ErrCanceled) {
				return nil
			}
			return
		}
		log.Println("")

		defer func() {
			terminateErr := b.dockerRunner.StopRunningContainers(ctx)
			if err == nil {
				err = terminateErr
			}

			if needsNetwork {
				if terminateErr = b.handleFunc(ctx, logger, func() error {
					return b.dockerRunner.NetworkRemove(ctx, logger)
				}); err == nil {
					if !errors.Is(terminateErr, ErrCanceled) {
						err = terminateErr
					}
				}
				log.Println("")
			}
		}()
	}

	// get metadata as envvars
	env := map[string]string{}
	env["INFINITY_METADATA_NAME"] = manifest.Metadata.Name
	env["INFINITY_METADATA_TYPE"] = string(manifest.Metadata.ApplicationType)
	env["INFINITY_METADATA_LANGUAGE"] = string(manifest.Metadata.Language)

	//  add/overwrite global environment variables
	for k, v := range manifest.Env {
		env[k] = v
	}

	// add/overwrite target environment variables
	for k, v := range manifestTarget.Env {
		env[k] = v
	}

	for _, stage := range manifestTarget.Stages {
		err = b.runStage(ctx, *stage, env, needsNetwork)
		log.Println("")
		if err != nil {
			return
		}
	}

	return nil
}

func (b *runner) getColorCode(stageIndex int) uint8 {

	availableColors := []uint8{11, 12, 13, 8, 6}

	return availableColors[stageIndex%len(availableColors)]
}

func (b *runner) setColorCode(stages []*ManifestStage) {
	for i, st := range stages {
		st.colorCode = b.getColorCode(i)
		b.setColorCode(st.Stages)
	}
}

func (b *runner) runStage(ctx context.Context, stage ManifestStage, env map[string]string, needsNetwork bool, prefixes ...string) (err error) {

	prefixes = append(prefixes, stage.Name)
	prefix := strings.Join(prefixes, "] [")

	logger := log.New(os.Stdout, aurora.Index(stage.colorCode, fmt.Sprintf("[%v] ", prefix)).String(), 0)

	if len(stage.Stages) > 0 {
		return b.runParallelStages(ctx, stage, env, needsNetwork)
	}

	// add and override with stage environment variables
	for k, v := range stage.Env {
		env[k] = v
	}

	// add parameters to envvars
	for k, v := range stage.Parameters {
		env[ToUpperSnakeCase("INFINITY_PARAMETER_"+k)] = fmt.Sprintf("%v", v)
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
			if err = b.handleFunc(ctx, logger, func() error {
				return b.dockerRunner.ContainerPull(ctx, logger, stage)
			}); err != nil {
				if errors.Is(err, ErrCanceled) {
					return nil
				}
				return
			}
		}

		if err = b.handleFunc(ctx, logger, func() error {
			return b.dockerRunner.ContainerStart(ctx, logger, stage, env, needsNetwork)
		}); err != nil {
			if errors.Is(err, ErrCanceled) {
				return nil
			}
			return
		}

		return nil

	case RunnerTypeHost:
		return b.hostRunner.RunStage(ctx, logger, stage, env)
	}

	return fmt.Errorf("runner %v is not supported", stage.RunnerType)
}

func (b *runner) runParallelStages(ctx context.Context, stage ManifestStage, env map[string]string, needsNetwork bool) (err error) {
	g, ctx := errgroup.WithContext(ctx)
	for _, s := range stage.Stages {
		s := s
		g.Go(func() error { return b.runStage(ctx, *s, env, needsNetwork, stage.Name) })
	}

	return g.Wait()
}

var (
	ErrCanceled = fmt.Errorf("This function got canceled")
)

func (b *runner) handleFunc(ctx context.Context, logger *log.Logger, funcToRun func() error) error {

	start := time.Now()
	err := funcToRun()
	elapsed := time.Since(start)

	select {
	case <-ctx.Done():
		if logger != nil {
			logger.Printf(aurora.Gray(12, "Canceled in %v").String(), aurora.BrightCyan(elapsed.String()))
		} else {
			log.Printf(aurora.Gray(12, "Canceled in %v").String(), aurora.BrightCyan(elapsed.String()))
		}
		return ErrCanceled
	default:
	}

	if err != nil {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		} else {
			log.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		}
		return err
	}

	if logger != nil {
		logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))
	} else {
		log.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))
	}

	return nil
}
