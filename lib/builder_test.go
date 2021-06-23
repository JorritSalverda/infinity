package lib

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/golang/mock/gomock"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedsIfInfinityManifestIsValid", func(t *testing.T) {
		builder := NewBuilder(NewManifestReader(), NewCommandRunner(), false, ".infinity-test.yaml")

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

		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--volume=/Users/jorrit/work/personal/infinity/lib:/work", "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Times(2)

		builder := NewBuilder(manifestReader, commandRunner, false, ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

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

		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).Times(1)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--volume=/Users/jorrit/work/personal/infinity/lib:/work", "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).AnyTimes()

		builder := NewBuilder(manifestReader, commandRunner, false, ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

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

		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--volume=/Users/jorrit/work/personal/infinity/lib:/work", "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Times(2)

		builder := NewBuilder(manifestReader, commandRunner, false, ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

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

		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).Times(1)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--volume=/Users/jorrit/work/personal/infinity/lib:/work", "--workdir=/work", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).AnyTimes()

		builder := NewBuilder(manifestReader, commandRunner, false, ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

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

		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		commandRunner := NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"pull", "alpine:3.13"})).AnyTimes()
		commandRunner.EXPECT().RunCommandWithLogger(gomock.Any(), gomock.Any(), gomock.Eq("docker"), gomock.Eq([]string{"run", "--rm", "--volume=/Users/jorrit/work/personal/infinity/lib:/work", "--workdir=/work", "--env=INFINITY_PARAMETER_CONTAINER_NAME=mycontainer", "--env=INFINITY_PARAMETER_VULNERABILITY_THRESHOLD=CRITICAL", "--entrypoint=/bin/sh", "alpine:3.13", "-c", `set -e ; printf "\033[38;5;244m> sleep 1\033[0m\n" ; sleep 1`})).Times(1)

		builder := NewBuilder(manifestReader, commandRunner, false, ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

}
