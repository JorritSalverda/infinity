package lib

import (
	"io/ioutil"
	"testing"

	"github.com/alecthomas/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalManifest(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		manifestData, err := ioutil.ReadFile(".infinity.yaml")
		assert.Nil(t, err)
		var manifest Manifest

		// act
		err = yaml.UnmarshalStrict(manifestData, &manifest)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(manifest.Build.Stages))
	})
}
