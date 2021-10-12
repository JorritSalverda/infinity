package lib

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestEnvToMap(t *testing.T) {
	t.Run("ReturnsEnvironmentVariablesSplitByEqualSymbolInMap", func(t *testing.T) {

		runner := commandRunner{}
		env := []string{
			"PATH=abc",
		}

		// act
		envMap := runner.envToMap(env)

		assert.Equal(t, 1, len(envMap))
		assert.Equal(t, "abc", envMap["PATH"])
	})

	t.Run("ReturnsEnvironmentVariablesWithoutValueSplitByEqualSymbolInMap", func(t *testing.T) {

		runner := commandRunner{}
		env := []string{
			"PATH=",
			"USER",
		}

		// act
		envMap := runner.envToMap(env)

		assert.Equal(t, 2, len(envMap))
		assert.Equal(t, "", envMap["PATH"])
		assert.Equal(t, "", envMap["USER"])
	})
}

func TestEnvToArray(t *testing.T) {
	t.Run("ReturnsEnvironmentVariablesJoinedByEqualSymbolInArray", func(t *testing.T) {

		runner := commandRunner{}
		env := map[string]string{
			"PATH": "abc",
		}

		// act
		envArray := runner.envToArray(env)

		assert.Equal(t, 1, len(envArray))
		assert.Equal(t, "PATH=abc", envArray[0])
	})

	t.Run("ReturnsEnvironmentVariablesWithoutValueJoinedByEqualSymbolInArray", func(t *testing.T) {

		runner := commandRunner{}
		env := map[string]string{
			"PATH": "",
			"USER": "",
		}

		// act
		envArray := runner.envToArray(env)

		assert.Equal(t, 2, len(envArray))
		assert.Equal(t, "PATH=", envArray[0])
		assert.Equal(t, "USER=", envArray[1])
	})
}

func TestOverrideEnvvars(t *testing.T) {
	t.Run("ReturnsEnvIfThereAreNoExtraEnv", func(t *testing.T) {
		runner := commandRunner{}
		env := []string{
			"PATH=abc",
		}

		// act
		envArray := runner.overrideEnvvars(env)

		assert.Equal(t, 1, len(envArray))
		assert.Equal(t, "PATH=abc", envArray[0])
	})

	t.Run("JoinsEnvWithExtraEnvIfThereIsNoOverlap", func(t *testing.T) {
		runner := commandRunner{}
		env := []string{
			"PATH=abc",
		}
		extraEnv := []string{
			"USER=root",
		}

		// act
		envArray := runner.overrideEnvvars(env, extraEnv...)

		assert.Equal(t, 2, len(envArray))
		assert.Equal(t, "PATH=abc", envArray[0])
		assert.Equal(t, "USER=root", envArray[1])
	})

	t.Run("OverridesEnvWithExtraEnvIfEnvarNameIsEqual", func(t *testing.T) {
		runner := commandRunner{}
		env := []string{
			"PATH=abc",
			"USER=root",
		}
		extraEnv := []string{
			"USER=me",
		}

		// act
		envArray := runner.overrideEnvvars(env, extraEnv...)

		assert.Equal(t, 2, len(envArray))
		assert.Equal(t, "PATH=abc", envArray[0])
		assert.Equal(t, "USER=me", envArray[1])
	})
}
