package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	ProjectPage string `yaml:"projectPage"`
	ProjectRepo string `yaml:"projectRepo"`

	// Release
	Version      string `yaml:"version"`
	ProviderType string `yaml:"providerType"`
	ReleaseNote  string `yaml:"releaseNote"`
	Visibility   string `yaml:"visibility"`
	StartScript  string `yaml:"startScript"`
}

func NewManifestFromFile(path string) (*Manifest, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read Manifest file: %v\n", err)
		return nil, err
	}
	manifest := Manifest{}
	err = yaml.Unmarshal(d, &manifest)
	if err != nil {
		log.Fatalf("Cannot unmarshall the Manifest file: %v\n", err)
		return nil, err
	}
	return &manifest, nil
}
