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
