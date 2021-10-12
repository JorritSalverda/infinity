package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	gomock "github.com/golang/mock/gomock"
)

func TestContainerStart(t *testing.T) {
	t.Run("PassesSnakeCasedEnvironmentVariableForEachStageParameter", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		stage := ManifestStage{
			Name:     "stage-1",
			Image:    "alpine:3.13",
			Commands: []string{"sleep 1"},
			Parameters: map[string]interface{}{
				"vulnerabilityThreshold": "CRITICAL",
				"containerName":          "mycontainer",
			},
		}
		stage.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf '\033[38;5;244m> %s\033[0m\n' 'sleep 1' ; sleep 1`})).Return([]byte("abcd\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"inspect", "--format='{{.State.ExitCode}}'", "abcd"})).Return([]byte("0\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"wait", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"rm", "--volumes", "abcd"})).Times(1)
		logger := log.New(os.Stdout, "", 0)

		runner := NewDockerRunner(commandRunner, randomStringGenerator, "")

		// act
		err = runner.ContainerStart(context.Background(), logger, stage, map[string]string{"INFINITY_PARAMETER_VULNERABILITY_THRESHOLD": "CRITICAL", "INFINITY_PARAMETER_CONTAINER_NAME": "mycontainer"}, false)

		assert.Nil(t, err)
	})

	t.Run("EscapeSingleQuotesInPrintfCommand", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		stage := ManifestStage{
			Name:     "stage-1",
			Image:    "alpine:3.13",
			Commands: []string{"echo '<xml />'"},
			Parameters: map[string]interface{}{
				"vulnerabilityThreshold": "CRITICAL",
				"containerName":          "mycontainer",
			},
		}
		stage.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf '\033[38;5;244m> %s\033[0m\n' 'echo \'<xml />\'' ; echo '<xml />'`})).Return([]byte("abcd\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"inspect", "--format='{{.State.ExitCode}}'", "abcd"})).Return([]byte("0\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"wait", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"rm", "--volumes", "abcd"})).Times(1)
		logger := log.New(os.Stdout, "", 0)

		runner := NewDockerRunner(commandRunner, randomStringGenerator, "")

		// act
		err = runner.ContainerStart(context.Background(), logger, stage, map[string]string{"INFINITY_PARAMETER_VULNERABILITY_THRESHOLD": "CRITICAL", "INFINITY_PARAMETER_CONTAINER_NAME": "mycontainer"}, false)

		assert.Nil(t, err)
	})

	t.Run("EscapeBackslashInPrintfCommand", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		stage := ManifestStage{
			Name:     "stage-1",
			Image:    "alpine:3.13",
			Commands: []string{`echo "a\nb"`},
			Parameters: map[string]interface{}{
				"vulnerabilityThreshold": "CRITICAL",
				"containerName":          "mycontainer",
			},
		}
		stage.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf '\033[38;5;244m> %s\033[0m\n' 'echo "a\\nb"' ; echo "a\nb"`})).Return([]byte("abcd\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"inspect", "--format='{{.State.ExitCode}}'", "abcd"})).Return([]byte("0\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"wait", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"rm", "--volumes", "abcd"})).Times(1)
		logger := log.New(os.Stdout, "", 0)

		runner := NewDockerRunner(commandRunner, randomStringGenerator, "")

		// act
		err = runner.ContainerStart(context.Background(), logger, stage, map[string]string{"INFINITY_PARAMETER_VULNERABILITY_THRESHOLD": "CRITICAL", "INFINITY_PARAMETER_CONTAINER_NAME": "mycontainer"}, false)

		assert.Nil(t, err)
	})
}
