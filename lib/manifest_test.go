package lib

import (
	"io/ioutil"
	"testing"

	"github.com/alecthomas/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalManifest(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		manifestData, err := ioutil.ReadFile(".infinity-test.yaml")
		assert.Nil(t, err)
		var manifest Manifest

		// act
		err = yaml.UnmarshalStrict(manifestData, &manifest)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(manifest.Build.Stages))

		assert.Equal(t, "test", manifest.Build.Stages[0].Name)
		assert.Equal(t, "golang:1.16-alpine", manifest.Build.Stages[0].Image)
		assert.Equal(t, 1, len(manifest.Build.Stages[0].Commands))
		assert.Equal(t, "go test -short ./...", manifest.Build.Stages[0].Commands[0])

		assert.Equal(t, "build", manifest.Build.Stages[1].Name)
		assert.Equal(t, "golang:1.16-alpine", manifest.Build.Stages[1].Image)
		assert.Equal(t, "0", manifest.Build.Stages[1].Env["CGO_ENABLED"])
		assert.Equal(t, 1, len(manifest.Build.Stages[1].Commands))
		assert.Equal(t, "go build -a -installsuffix cgo .", manifest.Build.Stages[1].Commands[0])
	})
}

func TestSetDefaultForManifest(t *testing.T) {
	t.Run("CallsSetDefaultOnBuildStages", func(t *testing.T) {
		manifest := Manifest{
			Build: ManifestBuild{
				Stages: []*ManifestStage{
					{
						RunnerType: RunnerTypeUnknown,
					},
				},
			},
		}

		// act
		manifest.SetDefault()

		assert.Equal(t, RunnerTypeContainer, manifest.Build.Stages[0].RunnerType)
	})
}

func TestValidateForManifest(t *testing.T) {
	t.Run("ReturnsNoErrorIfManifestIsValid", func(t *testing.T) {
		manifest := getValidManifest()

		// act
		_, errors := manifest.Validate()

		assert.Equal(t, 0, len(errors))
	})

	t.Run("ReturnsErrorIfApplicationTypeIsUnknown", func(t *testing.T) {
		manifest := getValidManifest()
		manifest.ApplicationType = ApplicationTypeUnknown

		// act
		_, errors := manifest.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "application is unknown; set to a supported application type with 'application: library|cli|firmware|api|web|controller'", errors[0].Error())
	})

	t.Run("ReturnsErrorIfLanguageIsUnknown", func(t *testing.T) {
		manifest := getValidManifest()
		manifest.Language = LanguageUnknown

		// act
		_, errors := manifest.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "language is unknown; set to a supported language with 'language: go|c|c++|java|csharp|python|node|rust|kotlin|swift|scala'", errors[0].Error())
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {
		manifest := getValidManifest()
		manifest.Name = ""

		// act
		_, errors := manifest.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "application has no name; please set 'name: <name>'", errors[0].Error())
	})

	t.Run("CallsValidateOnAllBuildStages", func(t *testing.T) {
		manifest := getValidManifest()
		manifest.Build.Stages[0].Name = ""

		// act
		_, errors := manifest.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[?] stage has no name; please set 'name: <name>'", errors[0].Error())
	})
}

func TestSetDefaultForManifestStage(t *testing.T) {
	t.Run("DefaultsRunnerTypeToContainerIfUnknown", func(t *testing.T) {
		stage := ManifestStage{
			RunnerType: RunnerTypeUnknown,
		}

		// act
		stage.SetDefault()

		assert.Equal(t, RunnerTypeContainer, stage.RunnerType)
	})

	t.Run("KeepsRunnerTypeIfNotUnknown", func(t *testing.T) {
		stage := ManifestStage{
			RunnerType: RunnerTypeHost,
		}

		// act
		stage.SetDefault()

		assert.Equal(t, RunnerTypeHost, stage.RunnerType)
	})

	t.Run("DefaultsWorkingDirectoryToWorkIfEmpty", func(t *testing.T) {
		stage := ManifestStage{
			WorkingDirectory: "",
		}

		// act
		stage.SetDefault()

		assert.Equal(t, "/work", stage.WorkingDirectory)
	})

	t.Run("KeepsWorkingDirectoryIfNotEmpty", func(t *testing.T) {
		stage := ManifestStage{
			WorkingDirectory: "/go/src/github.com/JorritSalverda/infinity",
		}

		// act
		stage.SetDefault()

		assert.Equal(t, "/go/src/github.com/JorritSalverda/infinity", stage.WorkingDirectory)
	})

	t.Run("DefaultsEnvToEmptyMapIfNil", func(t *testing.T) {
		stage := ManifestStage{
			Env: nil,
		}

		// act
		stage.SetDefault()

		assert.NotNil(t, stage.Env)
	})

	t.Run("CallsSetDefaultOnNestedStages", func(t *testing.T) {
		stage := ManifestStage{
			Stages: []*ManifestStage{
				{
					RunnerType: RunnerTypeUnknown,
				},
			},
		}

		// act
		stage.SetDefault()

		assert.Equal(t, RunnerTypeContainer, stage.Stages[0].RunnerType)
	})
}

func TestValidateForManifestStage(t *testing.T) {
	t.Run("ReturnsNoErrorIfStageIsValid", func(t *testing.T) {
		stage := getValidManifestStage()

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 0, len(errors))
	})

	t.Run("ReturnsErrorIfNameIsEmpty", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.Name = ""

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[?] stage has no name; please set 'name: <name>'", errors[0].Error())
	})

	t.Run("ReturnsErrorIfWorkingDirectoryIsEmpty", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.WorkingDirectory = ""

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[stage-1] work has no value; please set 'work: <working directory>'", errors[0].Error())
	})

	t.Run("ReturnsErrorIfRunnerTypeIsUnknown", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.RunnerType = RunnerTypeUnknown

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[stage-1] unknown runner; please set 'runner: container|host'", errors[0].Error())
	})

	t.Run("ReturnsNoErrorIfRunnerTypeIsUnknownAndStageHasNestedStages", func(t *testing.T) {
		innerStage := getValidManifestStage()
		stage := getValidManifestStage()
		stage.RunnerType = RunnerTypeUnknown
		stage.Stages = []*ManifestStage{
			&innerStage,
		}

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 0, len(errors))
	})

	t.Run("ReturnsErrorIfImageIsEmpty", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.Image = ""

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[stage-1] stage has no image; please set 'image: <image>'", errors[0].Error())
	})

	t.Run("ReturnsNoErrorIfImageIsEmptyAndStageHasNestedStages", func(t *testing.T) {
		innerStage := getValidManifestStage()
		stage := getValidManifestStage()
		stage.Image = ""
		stage.Stages = []*ManifestStage{
			&innerStage,
		}

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 0, len(errors))
	})

	t.Run("ReturnsErrorIfImageIsSetWhenRunnerTypeIsHost", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.RunnerType = RunnerTypeHost
		stage.Image = "jsalverda/arduino-cli:0.18.3"

		// act
		_, errors := stage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[stage-1] stage has image which is not supported in combination with 'runner: host'; please do not set 'image: <image>'", errors[0].Error())
	})

	t.Run("ReturnsWarningIfNoCommandsAreSet", func(t *testing.T) {
		stage := getValidManifestStage()
		stage.Commands = []string{}

		// act
		warnings, _ := stage.Validate()

		assert.Equal(t, 1, len(warnings))
		assert.Equal(t, "[stage-1] stage has no commands; you might want to define at least one command through 'commands'", warnings[0])
	})

	t.Run("CallsValidateOnNestedStages", func(t *testing.T) {

		innerStage := getValidManifestStage()
		innerStage.Name = ""

		outerStage := ManifestStage{
			Name: "outer-stage-1",
			Stages: []*ManifestStage{
				&innerStage,
			},
		}

		// act
		_, errors := outerStage.Validate()

		assert.Equal(t, 1, len(errors))
		assert.Equal(t, "[outer-stage-1] [?] stage has no name; please set 'name: <name>'", errors[0].Error())
	})
}

func getValidManifest() Manifest {
	stage := getValidManifestStage()

	manifest := Manifest{
		ApplicationType: ApplicationTypeAPI,
		Language:        LanguageGo,
		Name:            ":myapp",
		Build: ManifestBuild{
			Stages: []*ManifestStage{
				&stage,
			},
		},
	}
	manifest.SetDefault()

	return manifest
}

func getValidManifestStage() ManifestStage {
	stage := ManifestStage{
		Name:  "stage-1",
		Image: "jsalverda/arduino-cli:0.18.3",
		Commands: []string{
			"arduino-cli board list",
		},
	}
	stage.SetDefault()

	return stage
}
