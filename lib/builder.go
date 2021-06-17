package lib

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Builder interface {
	Validate(ctx context.Context) (err error)
	Build(ctx context.Context) (err error)
}

type builder struct {
}

func NewBuilder() Builder {
	return &builder{}
}

func (b *builder) Validate(ctx context.Context) (err error) {
	manifest, err := b.getManifest(ctx)
	if err != nil {
		return
	}

	err = manifest.Validate()
	if err != nil {
		return
	}

	return nil
}

func (b *builder) Build(ctx context.Context) (err error) {
	manifest, err := b.getManifest(ctx)
	if err != nil {
		return
	}

	err = manifest.Validate()
	if err != nil {
		return
	}

	err = b.runManifest(ctx, manifest)
	if err != nil {
		return
	}

	return nil
}

func (b *builder) getManifest(ctx context.Context) (manifest Manifest, err error) {
	manifestPath := ".infinity.yaml"

	// check if manifest exists
	if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
		return manifest, fmt.Errorf("Manifest %v does not exist, cannot continue", manifestPath)
	}

	// read manifest
	manifestBytes, err := ioutil.ReadFile(".infinity.yaml")
	if err != nil {
		return
	}

	// unmarshal bytes into manifest
	if err = yaml.UnmarshalStrict(manifestBytes, &manifest); err != nil {
		return manifest, fmt.Errorf("Manifest %v is not valid: %w", manifestPath, err)
	}

	return
}

func (b *builder) runManifest(ctx context.Context, manifest Manifest) (err error) {
	return nil
}
