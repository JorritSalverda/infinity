package lib

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./command_runner_mock.go -source=command_runner.go
type CommandRunner interface {
	RunCommand(ctx context.Context, logger *log.Logger, dir, command string, args []string) (err error)
	RunCommandWithEnv(ctx context.Context, logger *log.Logger, dir, command string, args []string, env []string) (err error)
	RunCommandWithOutput(ctx context.Context, logger *log.Logger, dir, command string, args []string) (output []byte, err error)
	RunCommandWithOutputAndEnv(ctx context.Context, logger *log.Logger, dir, command string, args []string, env []string) (output []byte, err error)
}

type commandRunner struct {
	verbose bool
}

func NewCommandRunner(verbose bool) CommandRunner {
	return &commandRunner{
		verbose: verbose,
	}
}

func (c *commandRunner) RunCommand(ctx context.Context, logger *log.Logger, dir, command string, args []string) (err error) {
	return c.RunCommandWithEnv(ctx, logger, dir, command, args, os.Environ())
}

func (c *commandRunner) RunCommandWithEnv(ctx context.Context, logger *log.Logger, dir, command string, args []string, env []string) (err error) {
	if c.verbose {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		} else {
			log.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		}
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = env
	cmd.Dir = dir

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
		if logger != nil {
			logger.Printf(scanner.Text())
		} else {
			log.Print(scanner.Text())
		}
	}

	// wait until the container is done
	if err = cmd.Wait(); err != nil {
		return
	}

	return nil
}

func (c *commandRunner) RunCommandWithOutput(ctx context.Context, logger *log.Logger, dir, command string, args []string) (output []byte, err error) {
	return c.RunCommandWithOutputAndEnv(ctx, logger, dir, command, args, os.Environ())
}

func (c *commandRunner) RunCommandWithOutputAndEnv(ctx context.Context, logger *log.Logger, dir, command string, args []string, env []string) (output []byte, err error) {
	if c.verbose {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		} else {
			log.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		}
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = env
	cmd.Dir = dir

	return cmd.CombinedOutput()
}
