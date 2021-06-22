package lib

type RunnerType string

const (
	RunnerTypeUnknown   RunnerType = ""
	RunnerTypeContainer RunnerType = "container"
	RunnerTypeMetal     RunnerType = "metal"
)

type runnerTypes []RunnerType

func (runnerTypes runnerTypes) ToStringArray() (result []string) {
	for _, r := range runnerTypes {
		result = append(result, string(r))
	}
	return
}

var SupportedRunnerTypes = runnerTypes{
	RunnerTypeContainer,
	RunnerTypeMetal,
}

func (runnerType RunnerType) IsSupported() bool {
	for _, r := range SupportedRunnerTypes {
		if runnerType == r {
			return true
		}
	}
	return false
}
