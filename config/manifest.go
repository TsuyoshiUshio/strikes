package config

import (
	"io/ioutil"
	"log"

	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Name        string `yaml:"name" validate:"required"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	ProjectPage string `yaml:"projectPage" validate:"url"`
	ProjectRepo string `yaml:"projectRepo" validate:"url"`

	// Release
	Version      string   `yaml:"version"`
	ProviderType string   `yaml:"providerType" validate:"providerType"`
	ReleaseNote  string   `yaml:"releaseNote"`
	Visibility   string   `yaml:"visibility" validate:"visibility"`
	StartScript  string   `yaml:"startScript"`
	ZipFileNames []string `yaml:"zipFileNames"`
}

func NewManifestFromFile(path string) (*Manifest, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read Manifest file: %v\n", err)
		return nil, err
	}
	manifest := Manifest{}
	log.Printf("[DEBUG] Read manifest: %s\n", string(d))
	err = yaml.Unmarshal(d, &manifest)
	if err != nil {
		log.Fatalf("Cannot unmarshall the Manifest file: %v\n", err)
		return nil, err
	}
	return &manifest, nil
}

func validateVisibility(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return (value == "public" || value == "preview" || value == "private")
}

func validateProviderType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return (value == "Terraform" || value == "ARM")
}

func (m *Manifest) Validate() error {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("providerType", validateProviderType)
	validate.RegisterValidation("visibility", validateVisibility)
	return validate.Struct(m)
}
