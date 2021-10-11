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
	RunStage(ctx context.Context, logger *log.Logger, stage ManifestStage, env map[string]string) (err error)
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

func (b *hostRunner) RunStage(ctx context.Context, logger *log.Logger, stage ManifestStage, env map[string]string) (err error) {
	logger.Printf(aurora.Gray(12, "Starting stage on host").String())

	// loop envvars in sorted order
	envArray := make([]string, 0, len(env))
	for k, v := range env {
		envArray = append(envArray, fmt.Sprintf("%v=%v", k, v))
	}

	start := time.Now()

	for _, c := range stage.Commands {
		logger.Printf(aurora.Gray(12, "> %v").String(), c)

		splitCommands := strings.Split(c, " ")
		err = b.commandRunner.RunCommandWithEnv(ctx, logger, b.buildDirectory, splitCommands[0], splitCommands[1:], envArray)
		if err != nil {
			break
		}
	}

	elapsed := time.Since(start)

	select {
	case <-ctx.Done():
		logger.Printf(aurora.Gray(12, "Canceled in %v").String(), aurora.BrightCyan(elapsed.String()))
		return nil
	default:
	}
	if err != nil {
		logger.Printf(aurora.Gray(12, "Failed in %v").String(), aurora.BrightRed(elapsed.String()))
		return fmt.Errorf("stage %v failed: %w", stage.Name, err)
	}
	logger.Printf(aurora.Gray(12, "Completed in %v").String(), aurora.BrightGreen(elapsed.String()))

	return nil
}
