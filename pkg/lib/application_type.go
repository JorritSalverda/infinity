package lib

type ApplicationType string

const (
	ApplicationTypeUnknown    ApplicationType = ""
	ApplicationTypeLibrary    ApplicationType = "library"
	ApplicationTypeCLI        ApplicationType = "cli"
	ApplicationTypeFirmware   ApplicationType = "firmware"
	ApplicationTypeAPI        ApplicationType = "api"
	ApplicationTypeWeb        ApplicationType = "web"
	ApplicationTypeController ApplicationType = "controller"
)

type applicationTypes []ApplicationType

func (applicationTypes applicationTypes) ToStringArray() (result []string) {
	for _, a := range applicationTypes {
		result = append(result, string(a))
	}
	return
}

var SupportedApplicationTypes = applicationTypes{
	ApplicationTypeLibrary,
	ApplicationTypeCLI,
	ApplicationTypeFirmware,
	ApplicationTypeAPI,
	ApplicationTypeWeb,
	ApplicationTypeController,
}

func (applicationType ApplicationType) IsSupported() bool {
	for _, a := range SupportedApplicationTypes {
		if applicationType == a {
			return true
		}
	}
	return false
}
