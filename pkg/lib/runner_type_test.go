package lib

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestSupportedRunnerTypes(t *testing.T) {
	t.Run("ReturnsAllEnumValuesExceptForUnknown", func(t *testing.T) {

		// act
		supportedRunnerTypes := SupportedRunnerTypes.ToStringArray()

		assert.Equal(t, 2, len(supportedRunnerTypes))
	})
}

func TestIsSupportedRunnerType(t *testing.T) {
	t.Run("ReturnsFalseForUnknownRunnerType", func(t *testing.T) {

		unknownRunnerType := RunnerType("unknown")

		// act
		isSupported := unknownRunnerType.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsFalseForRunnerTypeUnknown", func(t *testing.T) {

		// act
		isSupported := RunnerTypeUnknown.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsTrueForAllSupportedRunnerTypes", func(t *testing.T) {
		for _, runnerType := range SupportedRunnerTypes {
			// act
			isSupported := runnerType.IsSupported()
			assert.True(t, isSupported)
		}
	})
}
