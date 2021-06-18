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

func (m *Manifest) Validate() error {
	err := m.Build.Validate()
	if err != nil {
		return err
	}

	return nil
}

type ManifestBuild struct {
	Stages []*ManifestStage `yaml:"stages,omitempty" json:"stages,omitempty"`
}

func (b *ManifestBuild) SetDefault() {
	for _, s := range b.Stages {
		s.SetDefault()
	}
}

func (b *ManifestBuild) Validate() error {
	if len(b.Stages) == 0 {
		return fmt.Errorf("Manifest has no stages; define at least stage through 'build.stages'")
	}

	for _, s := range b.Stages {
		err := s.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type ManifestStage struct {
	Name       string            `yaml:"name,omitempty" json:"name,omitempty"`
	Image      string            `yaml:"image,omitempty" json:"image,omitempty"`
	Shell      string            `yaml:"shell,omitempty" json:"shell,omitempty"`
	Privileged bool              `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	Mounts     []string          `yaml:"mounts,omitempty" json:"mounts,omitempty"`
	Devices    []string          `yaml:"devices,omitempty" json:"devices,omitempty"`
	Env        map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Commands   []string          `yaml:"commands,omitempty" json:"commands,omitempty"`
}

func (s *ManifestStage) SetDefault() {
	if s.Shell == "" {
		s.Shell = "/bin/sh"
	}
}

func (s *ManifestStage) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Stage has no name; please set 'name: <name>'")
	}
	if s.Image == "" {
		return fmt.Errorf("Stage has no image; please set 'image: <image>'")
	}
	if s.Shell == "" {
		return fmt.Errorf("Stage has no shell; please set 'shell: <shell>'")
	}
	if len(s.Commands) == 0 {
		return fmt.Errorf("Stage has no commands; define at least stage through 'commands'")
	}

	return nil
}
