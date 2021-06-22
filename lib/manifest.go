package lib

import (
	"fmt"
)

type Manifest struct {
	Build ManifestBuild `yaml:"build,omitempty" json:"build,omitempty"`
}

func (m *Manifest) SetDefault() {
	m.Build.SetDefault()
}

func (m *Manifest) Validate() (warnings []string, errors []error) {
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
		errors = append(errors, fmt.Errorf("Manifest has no stages; define at least stage through 'build.stages'"))
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
	Image      string            `yaml:"image,omitempty" json:"image,omitempty"`
	BareMetal  bool              `yaml:"bareMetal,omitempty" json:"bareMetal,omitempty"`
	Shell      string            `yaml:"shell,omitempty" json:"shell,omitempty"`
	Privileged bool              `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	Mounts     []string          `yaml:"mounts,omitempty" json:"mounts,omitempty"`
	Devices    []string          `yaml:"devices,omitempty" json:"devices,omitempty"`
	Env        map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Commands   []string          `yaml:"commands,omitempty" json:"commands,omitempty"`
	Stages     []*ManifestStage  `yaml:"stages,omitempty" json:"stages,omitempty"`
}

func (s *ManifestStage) SetDefault() {
	if s.Shell == "" {
		s.Shell = "/bin/sh"
	}

	for _, st := range s.Stages {
		st.SetDefault()
	}
}

func (s *ManifestStage) Validate() (warnings []string, errors []error) {
	if s.Name == "" {
		errors = append(errors, fmt.Errorf("Stage has no name; please set 'name: <name>'"))
	}
	if len(s.Stages) == 0 && !s.BareMetal && s.Image == "" {
		errors = append(errors, fmt.Errorf("Stage has no image; please set 'image: <image>'"))
	}
	if len(s.Stages) == 0 && s.BareMetal && s.Image != "" {
		errors = append(errors, fmt.Errorf("Stage has image while bareMetal is set to true; please do not set 'image: <image>'"))
	}
	if len(s.Stages) == 0 && s.Shell == "" {
		errors = append(errors, fmt.Errorf("Stage has no shell; please set 'shell: <shell>'"))
	}
	if len(s.Stages) == 0 && len(s.Commands) == 0 {
		errors = append(errors, fmt.Errorf("Stage has no commands; define at least stage through 'commands'"))
	}

	for _, st := range s.Stages {
		w, e := st.Validate()
		warnings = append(warnings, w...)
		errors = append(errors, e...)
	}

	return
}
