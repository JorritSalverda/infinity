package lib

import (
	"context"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/golang/mock/gomock"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedsIfInfinityManifestIsValid", func(t *testing.T) {
		builder := NewBuilder(NewManifestReader(), NewDockerRunner(NewCommandRunner(false), NewRandomStringGenerator(), ""), NewHostRunner(NewCommandRunner(false), ""), false, "", ".infinity-test.yaml")

		// act
		_, err := builder.Validate(context.Background())

		assert.Nil(t, err)
	})
}

func TestBuild(t *testing.T) {
	t.Run("CallsContainerStartForEachStages", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
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
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerImageIsPulled(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).Times(2)

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerPullForEachStage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
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
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerImageIsPulled(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).AnyTimes()

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerStartForEachParallelStage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
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
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerImageIsPulled(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).Times(2)

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CallsContainerPullForEachParallelStage", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
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
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		dockerRunner.EXPECT().ContainerImageIsPulled(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).Times(2)
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(false)).AnyTimes()

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("CreatesNetworkIfAnyStagesNeedNetwork", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:       "stage-1",
						Background: true,
						Image:      "alpine:3.13",
						Commands:   []string{"sleep 1"},
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
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(true).Times(1)
		dockerRunner.EXPECT().NetworkCreate(gomock.Any(), gomock.Any()).Times(1)
		dockerRunner.EXPECT().ContainerImageIsPulled(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerPull(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		dockerRunner.EXPECT().ContainerStart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(true)).AnyTimes()
		dockerRunner.EXPECT().StopRunningContainers(gomock.Any()).Times(1)
		dockerRunner.EXPECT().NetworkRemove(gomock.Any(), gomock.Any()).Times(1)

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})

	t.Run("RunsHostRunForEachStageWithHostRunner", func(t *testing.T) {

		ctrl := gomock.NewController(t)

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name:       "stage-1",
						RunnerType: RunnerTypeHost,
						Commands:   []string{"sleep 1"},
					},
				},
			},
		}
		manifest.SetDefault()

		manifestReader := NewMockManifestReader(ctrl)
		dockerRunner := NewMockDockerRunner(ctrl)
		hostRunner := NewMockHostRunner(ctrl)

		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		dockerRunner.EXPECT().NeedsNetwork(gomock.Eq(manifest.Build.Stages)).Return(false).Times(1)
		hostRunner.EXPECT().RunStage(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		builder := NewBuilder(manifestReader, dockerRunner, hostRunner, false, "", ".infinity.yaml")

		// act
		err := builder.Build(context.Background())

		assert.Nil(t, err)
	})
}

func TestCancellation(t *testing.T) {
	t.Run("FirstFailingParallelStageWithHostRunnerCancelsOtherStages", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name: "parallel",
						Stages: []*ManifestStage{
							{
								Name:       "fails",
								RunnerType: RunnerTypeHost,
								Commands:   []string{"sleep 1s", "exit 1"},
							},
							{
								Name:       "gets-canceled",
								RunnerType: RunnerTypeHost,
								Commands:   []string{"sleep 25s"},
							},
						},
					},
				},
			},
		}
		manifest.SetDefault()
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		builder := NewBuilder(manifestReader, NewDockerRunner(NewCommandRunner(false), NewRandomStringGenerator(), ""), NewHostRunner(NewCommandRunner(false), ""), false, "", ".infinity.yaml")

		// act
		start := time.Now()
		err := builder.Build(ctx)
		elapsed := time.Since(start)

		assert.NotNil(t, err)
		assert.Equal(t, "stage fails failed: exec: \"exit\": executable file not found in $PATH", err.Error())
		assert.True(t, elapsed.Seconds() < 10)
	})

	t.Run("FirstFailingParallelStageWithContainerRunnerCancelsOtherStages", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		manifest := Manifest{
			Metadata: ManifestMetadata{
				ApplicationType: ApplicationTypeAPI,
				Language:        LanguageGo,
				Name:            "test-app",
			},
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						Name: "parallel",
						Stages: []*ManifestStage{
							{
								Name:       "fails",
								RunnerType: RunnerTypeContainer,
								Image:      "alpine:3.13",
								Commands:   []string{"sleep 1s", "exit 1"},
							},
							{
								Name:       "gets-canceled",
								RunnerType: RunnerTypeContainer,
								Image:      "alpine:3.13",
								Commands:   []string{"exec sleep 25s"},
							},
						},
					},
				},
			},
		}
		manifest.SetDefault()
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		manifestReader := NewMockManifestReader(ctrl)
		manifestReader.EXPECT().GetManifest(gomock.Any(), gomock.Eq(".infinity.yaml")).Return(manifest, nil)
		builder := NewBuilder(manifestReader, NewDockerRunner(NewCommandRunner(false), NewRandomStringGenerator(), ""), NewHostRunner(NewCommandRunner(false), ""), false, "", ".infinity.yaml")

		// act
		start := time.Now()
		err := builder.Build(ctx)
		elapsed := time.Since(start)

		assert.NotNil(t, err)
		assert.True(t, elapsed.Seconds() < 10)
	})
}
