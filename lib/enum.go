package lib

type ApplicationType string

const (
	ApplicationTypeUnknown  ApplicationType = ""
	ApplicationTypeLibrary  ApplicationType = "library"
	ApplicationTypeCLI      ApplicationType = "cli"
	ApplicationTypeFirmware ApplicationType = "firmware"
	ApplicationTypeAPI      ApplicationType = "api"
	ApplicationTypeWeb      ApplicationType = "web"
)

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

type RunnerType string

const (
	RunnerUnknown   RunnerType = ""
	RunnerContainer RunnerType = "container"
	RunnerMetal     RunnerType = "metal"
)
