package lib

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora"
)

//go:generate mockgen -package=lib -destination ./command_runner_mock.go -source=command_runner.go
type CommandRunner interface {
	RunCommand(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) (err error)
	RunCommandWithOutput(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) (output []byte, err error)
}

type commandRunner struct {
	verbose bool
}

func NewCommandRunner(verbose bool) CommandRunner {
	return &commandRunner{
		verbose: verbose,
	}
}

func (c *commandRunner) RunCommand(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) (err error) {
	if c.verbose {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		} else {
			log.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		}
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = c.overrideEnvvars(os.Environ(), env...)
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

func (c *commandRunner) RunCommandWithOutput(ctx context.Context, logger *log.Logger, dir, command string, args []string, env ...string) (output []byte, err error) {
	if c.verbose {
		if logger != nil {
			logger.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		} else {
			log.Printf(aurora.Gray(12, "> %v %v").String(), command, strings.Join(args, " "))
		}
	}

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = c.overrideEnvvars(os.Environ(), env...)
	cmd.Dir = dir

	return cmd.CombinedOutput()
}

func (c *commandRunner) overrideEnvvars(env []string, extraEnv ...string) (combinedEnv []string) {
	// convert to map
	envMap := c.envToMap(env)

	// add/overwrite keys with extra values
	for k, v := range c.envToMap(extraEnv) {
		envMap[k] = v
	}

	// convert to array
	combinedEnv = c.envToArray(envMap)

	return
}

func (c *commandRunner) envToMap(envArray []string) (envMap map[string]string) {
	envMap = make(map[string]string)
	for _, e := range envArray {
		envSplit := strings.Split(e, "=")
		if len(envSplit) == 1 {
			envMap[envSplit[0]] = ""
			continue
		}
		if len(envSplit) > 1 {
			envMap[envSplit[0]] = strings.Join(envSplit[1:], "=")
		}
	}

	return
}

func (c *commandRunner) envToArray(envMap map[string]string) (envArray []string) {
	// loop envvars in sorted order
	envKeys := make([]string, 0, len(envMap))
	for k := range envMap {
		envKeys = append(envKeys, k)
	}
	sort.Strings(envKeys)
	for _, k := range envKeys {
		envArray = append(envArray, fmt.Sprintf("%v=%v", k, envMap[k]))
	}

	return
}
