package lib

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
)

func TestValidate(t *testing.T) {
	t.Run("SucceedsIfInfinityManifestIsValid", func(t *testing.T) {
		builder := NewBuilder()

		// act
		err := builder.Validate(context.Background())

		assert.Nil(t, err)
	})
}
