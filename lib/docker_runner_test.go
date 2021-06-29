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
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf '\033[38;5;244m> %s\033[0m\n' "sleep 1" ; sleep 1`})).Times(1)
		logger := log.New(os.Stdout, "", 0)

		runner := NewDockerRunner(commandRunner, randomStringGenerator, "")

		// act
		err = runner.ContainerStart(context.Background(), logger, stage, false)

		assert.Nil(t, err)
	})
}
