package lib

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestSupportedApplicationTypes(t *testing.T) {
	t.Run("ReturnsAllEnumValuesExceptForUnknown", func(t *testing.T) {

		// act
		supportedApplicationTypes := SupportedApplicationTypes.ToStringArray()

		assert.Equal(t, 5, len(supportedApplicationTypes))
	})
}

func TestIsSupportedApplicationType(t *testing.T) {
	t.Run("ReturnsFalseForUnknownApplicationType", func(t *testing.T) {

		unknownApplicationType := ApplicationType("unknown")

		// act
		isSupported := unknownApplicationType.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsFalseForApplicationTypeUnknown", func(t *testing.T) {

		// act
		isSupported := ApplicationTypeUnknown.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsTrueForAllSupportedApplicationTypes", func(t *testing.T) {
		for _, applicationType := range SupportedApplicationTypes {
			// act
			isSupported := applicationType.IsSupported()
			assert.True(t, isSupported)
		}
	})
}
