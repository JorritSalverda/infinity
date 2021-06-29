package lib

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/golang/mock/gomock"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedsIfInfinityManifestIsValid", func(t *testing.T) {
		builder := NewBuilder(NewManifestReader(), NewDockerRunner(NewCommandRunner(false), NewRandomStringGenerator(), ""), NewMetalRunner(NewCommandRunner(false), ""), "", ".infinity-test.yaml")

		// act
		_, err := builder.Validate(context.Background())

		assert.Nil(t, err)
	})
}

func TestBuild(t *testing.T) {
	t.Run("CallsContainerStartForEachStages", func(t *testing.T) {

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
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).Times(2)

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerPullForEachStage", func(t *testing.T) {

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
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).AnyTimes()

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerStartForEachParallelStage", func(t *testing.T) {

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
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).Times(2)

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerPullForEachParallelStage", func(t *testing.T) {

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
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).AnyTimes()

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CreatesNetworkIfAnyStagesNeedNetwork", func(t *testing.T) {

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

		manifestReader := NewMockManifestReader(ctrl)
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(true).Times(1)
		dockerRunner.EXPECT().NetworkCreate(gomock.Any()).Times(1)
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(true)).AnyTimes()
		dockerRunner.EXPECT().StopRunningContainers(gomock.Any()).Times(1)
		dockerRunner.EXPECT().NetworkRemove(gomock.Any()).Times(1)

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("RunsMetalRunForEachStageWithMetalRunner", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			ApplicationType: ApplicationTypeAPI,
			Language:        LanguageGo,
			Name:            "test-app",
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:       "stage-1",
						RunnerType: RunnerTypeMetal,
						Commands:   []string{"sleep 1"},
					},
				},
			},
		}
		manifest.SetDefault()

		manifestReader := NewMockManifestReader(ctrl)
		dockerRunner := NewMockDockerRunner(ctrl)
		metalRunner := NewMockMetalRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		metalRunner.EXPECT().MetalRun(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		builder := NewBuilder(manifestReader, dockerRunner, metalRunner, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})
}
