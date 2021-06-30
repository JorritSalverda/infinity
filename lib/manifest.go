package lib

import (
	"fmt"
	"strings"
)

type Manifest struct {
	ApplicationType ApplicationType `yaml:"type,omitempty" json:"type,omitempty"`
	Language        Language        `yaml:"language,omitempty" json:"language,omitempty"`
	Name            string          `yaml:"name,omitempty" json:"name,omitempty"`
	Build           ManifestBuild   `yaml:"build,omitempty" json:"build,omitempty"`
}

func (m *Manifest) SetDefault() {
	m.Build.SetDefault()
}

func (m *Manifest) Validate() (warnings []string, errors []error) {
	if !m.ApplicationType.IsSupported() {
		errors = append(errors, fmt.Errorf("application is unknown; set to a supported application type with 'application: %v'", strings.Join(SupportedApplicationTypes.ToStringArray(), "|")))
	}
	if !m.Language.IsSupported() {
		errors = append(errors, fmt.Errorf("language is unknown; set to a supported language with 'language: %v'", strings.Join(SupportedLanguages.ToStringArray(), "|")))
	}
	if m.Name == "" {
		errors = append(errors, fmt.Errorf("application has no name; please set 'name: <name>'"))
	}

	w, e := m.Build.Validate()
	warnings = append(warnings, w...)
	errors = append(errors, e...)

	return
}

type ManifestBuild struct {
	Stages []*ManifestStage `yaml:"stages,omitempty" json:"stages,omitempty"`
}

func (b *ManifestBuild) SetDefault() {
	for _, s := range b.Stages {
		s.SetDefault()
	}
}

func (b *ManifestBuild) Validate() (warnings []string, errors []error) {
	if len(b.Stages) == 0 {
		errors = append(errors, fmt.Errorf("manifest has no stages; define at least stage through 'build.stages'"))
	}

	for _, s := range b.Stages {
		w, e := s.Validate()
		warnings = append(warnings, w...)
		errors = append(errors, e...)
	}

	return
}

type ManifestStage struct {
	Name                  string                 `yaml:"name,omitempty" json:"name,omitempty"`
	RunnerType            RunnerType             `yaml:"runner,omitempty" json:"runner,omitempty"`
	Image                 string                 `yaml:"image,omitempty" json:"image,omitempty"`
	Background            bool                   `yaml:"background,omitempty" json:"background,omitempty"`
	Privileged            bool                   `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	MountWorkingDirectory *bool                  `yaml:"mount,omitempty" json:"mount,omitempty"`
	WorkingDirectory      string                 `yaml:"work,omitempty" json:"work,omitempty"`
	Volumes               []string               `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Devices               []string               `yaml:"devices,omitempty" json:"devices,omitempty"`
	Env                   map[string]string      `yaml:"env,omitempty" json:"env,omitempty"`
	Commands              []string               `yaml:"commands,omitempty" json:"commands,omitempty"`
	Stages                []*ManifestStage       `yaml:"stages,omitempty" json:"stages,omitempty"`
	Parameters            map[string]interface{} `yaml:",inline"`
}

func (s *ManifestStage) SetDefault() {
	if s.RunnerType == RunnerTypeUnknown {
		s.RunnerType = RunnerTypeContainer
	}
	if s.MountWorkingDirectory == nil {
		defaultValue := true
		s.MountWorkingDirectory = &defaultValue
	}
	if s.WorkingDirectory == "" {
		s.WorkingDirectory = "/work"
	}
	if s.Env == nil {
		s.Env = make(map[string]string)
	}
	for _, st := range s.Stages {
		st.SetDefault()
	}
}

func (s *ManifestStage) Validate(prefixes ...string) (warnings []string, errors []error) {

	if s.Name != "" {
		prefixes = append(prefixes, s.Name)
	} else {
		prefixes = append(prefixes, "?")
	}

	prefix := strings.Join(prefixes, "] [")

	if s.Name == "" {
		errors = append(errors, fmt.Errorf("[%v] stage has no name; please set 'name: <name>'", prefix))
	}
	if len(s.Stages) == 0 {
		if s.MountWorkingDirectory == nil {
			errors = append(errors, fmt.Errorf("[%v] mount has no value; please set 'mountWork: true|false'", prefix))
		}
		if s.WorkingDirectory == "" {
			errors = append(errors, fmt.Errorf("[%v] work has no value; please set 'work: <working directory>'", prefix))
		}
		if s.RunnerType == RunnerTypeUnknown {
			errors = append(errors, fmt.Errorf("[%v] unknown runner; please set 'runner: %v'", prefix, strings.Join(SupportedRunnerTypes.ToStringArray(), "|")))
		}
		if len(s.Commands) == 0 && !s.Background {
			warnings = append(warnings, fmt.Sprintf("[%v] stage has no commands; you might want to define at least one command through 'commands'", prefix))
		}

		switch s.RunnerType {
		case RunnerTypeContainer:
			if s.Image == "" {
				errors = append(errors, fmt.Errorf("[%v] stage has no image; please set 'image: <image>'", prefix))
			}
		case RunnerTypeHost:
			if s.Image != "" {
				errors = append(errors, fmt.Errorf("[%v] stage has image which is not supported in combination with 'runner: host'; please do not set 'image: <image>'", prefix))
			}
		}
	}

	for _, st := range s.Stages {

		w, e := st.Validate(prefixes...)
		warnings = append(warnings, w...)
		errors = append(errors, e...)
	}

	return
}
