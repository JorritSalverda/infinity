package lib

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
)

//go:generate mockgen -package=lib -destination ./command_runner_mock.go -source=command_runner.go
type CommandRunner interface {
	RunCommandWithLogger(ctx context.Context, logger *log.Logger, command string, args []string) (err error)
	RunCommandWithOutput(ctx context.Context, command string, args []string) (output []byte, err error)
}

type commandRunner struct {
}

func NewCommandRunner() CommandRunner {
	return &commandRunner{}
}

func (c *commandRunner) RunCommandWithLogger(ctx context.Context, logger *log.Logger, command string, args []string) (err error) {
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

	// start container
	if err = cmd.Start(); err != nil {
		return
	}

	// tail logs with custom logger
	multi := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(multi)
	for scanner.Scan() {
		logger.Printf(scanner.Text())
	}

	// wait until the container is done
	if err = cmd.Wait(); err != nil {
		return
	}

	return nil
}

func (c *commandRunner) RunCommandWithOutput(ctx context.Context, command string, args []string) (output []byte, err error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = os.Environ()

	return cmd.CombinedOutput()
}
