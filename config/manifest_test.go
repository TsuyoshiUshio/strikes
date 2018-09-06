package config

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestReadManifest(t *testing.T) {
	manifest, _ := NewManifestFromFile("./test-fixture/manifest-basic/manifest.yaml")
	assert.Equal(t, "bar", manifest.Name)
	assert.Equal(t, "Explanation of bar", manifest.Description)
	assert.Equal(t, "Foo Bar", manifest.Author)
	assert.Equal(t, "https://github.com/foo/bar", manifest.ProjectPage)
	assert.Equal(t, "https://foo.bar.com", manifest.ProjectRepo)
	assert.Equal(t, "1.0.0", manifest.Version)
	assert.Equal(t, "Terraform", manifest.ProviderType)
	assert.Equal(t, "Explanation of this release", manifest.ReleaseNote)
	assert.Equal(t, "public", manifest.Visibility)
	assert.Equal(t, "terraform.tf", manifest.StartScript)
}
func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
func TestReadManifestWithoutFile(t *testing.T) {
	fakeExit := func(int) {
		// do nothing
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	//var err error.Error
	output := captureOutput(func() {
		_, _ = NewManifestFromFile("./test-fixture/manifest-basic/foo.yaml")
	})

	assert.Regexp(t, "Cannot read Manifest file: open ./test-fixture/manifest-basic/foo.yaml: no such file or directory", output)

}

func TestReadManifestWithMissingColumn(t *testing.T) {
	fakeExit := func(int) {
		// do nothing
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	//var err error.Error
	output := captureOutput(func() {
		_, _ = NewManifestFromFile("./test-fixture/manifest-wrong-yaml/manifest.yaml")
	})
	assert.Regexp(t, "Cannot unmarshall the Manifest file: yaml: mapping values are not allowed in this context\n", output)

}

func validateValidationError(t *testing.T, tag string, field string, validationError *validator.FieldError) {
	assert.Equal(t, tag, (*validationError).Tag())
	assert.Equal(t, field, (*validationError).Field())
}

func TestValidateManifestWithSuccess(t *testing.T) {
	manifest, _ := NewManifestFromFile("./test-fixture/manifest-basic/manifest.yaml")
	err := manifest.Validate()
	assert.Nil(t, err, "err should be nil.")
}

func TestValidateManifestWithFailure(t *testing.T) {
	manifest, _ := NewManifestFromFile("./test-fixture/manifest-validation-fail/manifest.yaml")
	err := manifest.Validate()
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, 5, len(validationErrors), "Number of the error should be 5.")
	validateValidationError(t, "required", "Name", &(validationErrors[0]))
	validateValidationError(t, "url", "ProjectPage", &(validationErrors[1]))
	validateValidationError(t, "url", "ProjectRepo", &(validationErrors[2]))
	validateValidationError(t, "providerType", "ProviderType", &validationErrors[3])
	validateValidationError(t, "visibility", "Visibility", &(validationErrors[4]))
}
