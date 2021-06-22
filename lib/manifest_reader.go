package lib

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//go:generate mockgen -package=lib -destination ./manifest_reader_mock.go -source=manifest_reader.go
type ManifestReader interface {
	GetManifest(ctx context.Context, buildManifestFilename string) (manifest Manifest, err error)
}

type manifestReader struct {
}

func NewManifestReader() ManifestReader {
	return &manifestReader{}
}

func (b *manifestReader) GetManifest(ctx context.Context, buildManifestFilename string) (manifest Manifest, err error) {
	// check if manifest exists
	if _, err = os.Stat(buildManifestFilename); os.IsNotExist(err) {
		return manifest, fmt.Errorf("manifest %v does not exist, cannot continue", buildManifestFilename)
	}

	// read manifest
	manifestBytes, err := ioutil.ReadFile(buildManifestFilename)
	if err != nil {
		return
	}

	// unmarshal bytes into manifest
	if err = yaml.UnmarshalStrict(manifestBytes, &manifest); err != nil {
		return manifest, fmt.Errorf("manifest %v is invalid: %w", buildManifestFilename, err)
	}

	manifest.SetDefault()

	return
}
