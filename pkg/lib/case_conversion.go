package lib

import (
	"fmt"
	"regexp"
	"unicode"
)

// ToUpperSnakeCase turns any input string into an upper snake cased string
func ToUpperSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToUpper(runes[i]))
	}

	snake := string(out)

	// make sure nothing but alphanumeric characters and underscores are returned
	reg, err := regexp.Compile("[^A-Z0-9]+")
	if err != nil {
		panic(fmt.Errorf("Failed converting %v to upper snake case", in))
	}
	cleanSnake := reg.ReplaceAllString(snake, "_")

	return cleanSnake
}

// ToLowerSnakeCase turns any input string into an lower snake cased string
func ToLowerSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	snake := string(out)

	// make sure nothing but alphanumeric characters and underscores are returned
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		panic(fmt.Errorf("Failed converting %v to lower snake case", in))
	}
	cleanSnake := reg.ReplaceAllString(snake, "_")

	return cleanSnake
}
