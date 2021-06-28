package lib

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./metal_runner_mock.go -source=metal_runner.go
type MetalRunner interface {
	MetalRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error)
}

type metalRunner struct {
	commandRunner  CommandRunner
	buildDirectory string
}

func NewMetalRunner(commandRunner CommandRunner, buildDirectory string) MetalRunner {

	return &metalRunner{
		commandRunner:  commandRunner,
		buildDirectory: buildDirectory,
	}
}

func (b *metalRunner) MetalRun(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	logger.Printf(aurora.Gray(12, "Executing commands in bare metal mode").String())

	start := time.Now()

	for _, c := range stage.Commands {
		logger.Printf(aurora.Gray(12, "> %v").String(), c)

		splitCommands := strings.Split(c, " ")
		err = b.commandRunner.RunCommand(ctx, logger, b.buildDirectory, splitCommands[0], splitCommands[1:])
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
