package lib

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./host_runner_mock.go -source=host_runner.go
type HostRunner interface {
	RunStage(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error)
}

type hostRunner struct {
	commandRunner  CommandRunner
	buildDirectory string
}

func NewHostRunner(commandRunner CommandRunner, buildDirectory string) HostRunner {

	return &hostRunner{
		commandRunner:  commandRunner,
		buildDirectory: buildDirectory,
	}
}

func (b *hostRunner) RunStage(ctx context.Context, logger *log.Logger, stage ManifestStage) (err error) {
	logger.Printf(aurora.Gray(12, "Executing commands on host").String())

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
