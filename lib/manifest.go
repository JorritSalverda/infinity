package lib

import (
	"fmt"
)

type Manifest struct {
	ApplicationType ApplicationType `yaml:"type,omitempty" json:"type,omitempty"`
	Language        Language        `yaml:"language,omitempty" json:"language,omitempty"`
	Build           ManifestBuild   `yaml:"build,omitempty" json:"build,omitempty"`
}

func (m *Manifest) SetDefault() {
	m.Build.SetDefault()
}

func (m *Manifest) Validate() (warnings []string, errors []error) {
	if m.ApplicationType == ApplicationTypeUnknown {
		errors = append(errors, fmt.Errorf("application is unknown; set to a supported application type with 'application: <application>'"))
	}
	if m.Language == LanguageUnknown {
		errors = append(errors, fmt.Errorf("language is unknown; set to a supported language with 'language: <language>'"))
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
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	RunnerType RunnerType        `yaml:"runner,omitempty" json:"runner,omitempty"`
	Image      string            `yaml:"image,omitempty" json:"image,omitempty"`
	Detach     bool              `yaml:"detach,omitempty" json:"detach,omitempty"`
	Privileged bool              `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	Mounts     []string          `yaml:"mounts,omitempty" json:"mounts,omitempty"`
	Devices    []string          `yaml:"devices,omitempty" json:"devices,omitempty"`
	Env        map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Commands   []string          `yaml:"commands,omitempty" json:"commands,omitempty"`
	Stages     []*ManifestStage  `yaml:"stages,omitempty" json:"stages,omitempty"`
}

func (s *ManifestStage) SetDefault() {
	if s.RunnerType == RunnerUnknown {
		s.RunnerType = RunnerContainer
	}
	for _, st := range s.Stages {
		st.SetDefault()
	}
}

func (s *ManifestStage) Validate() (warnings []string, errors []error) {
	if s.Name == "" {
		errors = append(errors, fmt.Errorf("stage has no name; please set 'name: <name>'"))
	}
	if s.RunnerType == RunnerUnknown {
		errors = append(errors, fmt.Errorf("unknown runner; please set 'runner: container|metal'"))
	}

	switch s.RunnerType {
	case RunnerContainer:
		if len(s.Stages) == 0 && s.Image == "" {
			errors = append(errors, fmt.Errorf("stage has no image; please set 'image: <image>'"))
		}
	case RunnerMetal:
		if len(s.Stages) == 0 && s.Image != "" {
			errors = append(errors, fmt.Errorf("stage has image which is not supported in combination with 'runner: metal'; please do not set 'image: <image>'"))
		}
	}

	if len(s.Stages) == 0 && len(s.Commands) == 0 {
		errors = append(errors, fmt.Errorf("stage has no commands; define at least stage through 'commands'"))
	}

	for _, st := range s.Stages {
		w, e := st.Validate()
		warnings = append(warnings, w...)
		errors = append(errors, e...)
	}

	return
}
