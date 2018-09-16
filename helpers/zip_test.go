package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipScenario(t *testing.T) {
	// Zip one directry recursively
	// extract zip file in a temp folder
	// check if there is a file on the deepest directry
	tempPath := filepath.Join(".", ".test")

	// If the tempPath exists, remove it
	if _, err := os.Stat(tempPath); err == nil {
		os.RemoveAll(tempPath)
	}

	os.MkdirAll(tempPath, os.ModePerm)
	zipFilePath := filepath.Join(tempPath, "circuit.zip")
	Zip(filepath.Join("..", "samples", "hello-world", "circuit"), zipFilePath)
	UnZip(zipFilePath, tempPath)
	manifestPath := filepath.Join(tempPath, "circuit", "manifest.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		assert.Fail(t, "Unzipped file  manifest can not found at :"+manifestPath)
	}

	// clean up

	os.RemoveAll(tempPath)
}
