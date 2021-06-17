package lib

import "fmt"

type Manifest struct {
	Build ManifestBuild `yaml:"build,omitempty" json:"build,omitempty"`
}

type ManifestBuild struct {
	Stages []ManifestStage `yaml:"stages,omitempty" json:"stages,omitempty"`
}

type ManifestStage struct {
	Name     string            `yaml:"name,omitempty" json:"name,omitempty"`
	Image    string            `yaml:"image,omitempty" json:"image,omitempty"`
	Env      map[string]string `yaml:"env,omitempty" json:"env,omitempty"`
	Commands []string          `yaml:"commands,omitempty" json:"commands,omitempty"`
}

func (m *Manifest) Validate() error {
	if len(m.Build.Stages) == 0 {
		return fmt.Errorf("Manifest has no stages; define at least stage through build.stages")
	}

	return nil
}
