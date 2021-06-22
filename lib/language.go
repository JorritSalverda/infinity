package lib

type Language string

const (
	LanguageUnknown   Language = ""
	LanguageGo        Language = "go"
	LanguageC         Language = "c"
	LanguageCPlusPlus Language = "c++"
	LanguageJava      Language = "java"
	LanguageCSharp    Language = "csharp"
	LanguagePython    Language = "python"
	LanguageNode      Language = "node"
)

type languages []Language

func (languages languages) ToStringArray() (result []string) {
	for _, l := range languages {
		result = append(result, string(l))
	}
	return
}

var SupportedLanguages = languages{
	LanguageGo,
	LanguageC,
	LanguageCPlusPlus,
	LanguageJava,
	LanguageCSharp,
	LanguagePython,
	LanguageNode,
}

func (language Language) IsSupported() bool {
	for _, l := range SupportedLanguages {
		if language == l {
			return true
		}
	}
	return false
}
