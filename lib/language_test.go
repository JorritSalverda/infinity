package lib

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestSupportedLanguages(t *testing.T) {
	t.Run("ReturnsAllEnumValuesExceptForUnknown", func(t *testing.T) {

		// act
		supportedLanguages := SupportedLanguages.ToStringArray()

		assert.Equal(t, 7, len(supportedLanguages))
	})
}

func TestIsSupportedLanguage(t *testing.T) {
	t.Run("ReturnsFalseForUnknownLanguage", func(t *testing.T) {

		unknownLanguage := Language("unknown")

		// act
		isSupported := unknownLanguage.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsFalseForLanguageUnknown", func(t *testing.T) {

		// act
		isSupported := LanguageUnknown.IsSupported()

		assert.False(t, isSupported)
	})

	t.Run("ReturnsTrueForAllSupportedLanguages", func(t *testing.T) {
		for _, language := range SupportedLanguages {
			// act
			isSupported := language.IsSupported()
			assert.True(t, isSupported)
		}
	})
}
