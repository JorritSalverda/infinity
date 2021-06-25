package lib

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/golang/mock/gomock"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedsIfInfinityManifestIsValid", func(t *testing.T) {
		builder := NewBuilder(NewManifestReader(), NewCommandRunner(false), NewRandomStringGenerator(), "", ".infinity-test.yaml")

		// act
		_, err := builder.Validate(context.Background())

		assert.Nil(t, err)
	})
}

func TestBuild(t *testing.T) {
	t.Run("RunsDockerRunForEachStage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:     "stage-1",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
					{
						Name:     "stage-2",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).Times(2)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(2)

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("RunsDockerPullOnceForEachImage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:     "stage-1",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
					{
						Name:     "stage-2",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).AnyTimes()
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).AnyTimes()

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("RunsDockerRunForEachParallelStage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name: "parallel",
						Stages: []*ManifestStage{
							{
								Name:     "stage-1",
								Image:    "alpine:3.13",
								Commands: []string{"sleep 1"},
							},
							{
								Name:     "stage-2",
								Image:    "alpine:3.13",
								Commands: []string{"sleep 1"},
							},
						},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).Times(2)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(2)

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("RunsDockerPullOnceForEachParallelImage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name: "parallel",
						Stages: []*ManifestStage{
							{
								Name:     "stage-1",
								Image:    "alpine:3.13",
								Commands: []string{"sleep 1"},
							},
							{
								Name:     "stage-2",
								Image:    "alpine:3.13",
								Commands: []string{"sleep 1"},
							},
						},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).AnyTimes()
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).AnyTimes()

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("PassesSnakeCasedEnvironmentVariableForEachStageParameter", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:     "stage-1",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
						Parameters: map[string]interface{}{
							"vulnerabilityThreshold": "CRITICAL",
							"containerName":          "mycontainer",
						},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(1)

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CreatesNetworkIfAnyStagesHaveDetachedTrue", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:     "stage-1",
						Detach:   true,
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
					{
						Name:     "stage-2",
						Image:    "alpine:3.13",
						Commands: []string{"sleep 1"},
					},
				},
			},
		}
		manifest.SetDefault()

		pwd, err := os.Getwd()
		assert.Nil(t, err)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		randomStringGenerator := NewMockRandomStringGenerator(ctrl)
		randomStringGenerator.EXPECT().GenerateRandomString(10).Return("abcdefghij").Times(1)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Nil(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"network", "create", "infinity-abcdefghij"})).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Nil(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"network", "rm", "infinity-abcdefghij"})).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", "--network=infinity-abcdefghij", "--name=stage-1", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("abcd\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommandWithOutput(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--detach", "--network=infinity-abcdefghij", fmt.Sprintf("--volume=%v:/work", pwd), "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Return([]byte("efgh\n"), nil).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "efgh"})).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"logs", "--follow", "abcd"})).Times(1)
		commandRunner.EXPECT().RunCommand(gomock.Any(), gomock.Any(), gomock.Eq(""), gomock.Eq("docker"), gomock.Eq([]string{"stop", "--time=30", "abcd"})).Times(1)

		builder := NewBuilder(manifestReader, commandRunner, randomStringGenerator, "", ".infinity.yaml")

		// act
		err = builder.Build(context.Background())

		assert.Nil(t, err)
	})

}
