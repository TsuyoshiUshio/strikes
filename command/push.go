package command

import (
	"log"
	"os"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/urfave/cli"
)

type PushCommand struct {
}

func (p *PushCommand) Push(c *cli.Context) error {

	packageDirBase := c.String("p")
	// Read the manifest file
	packageDir := filepath.Join(packageDirBase, "circuit", "manifest.yaml")
	// Read manifest file
	manifest, err := config.NewManifestFromFile(packageDir)
	if err != nil {
		log.Fatal(err)
	}
	// Validate the manifest file
	err = manifest.Validate()
	if err != nil {
		log.Fatal(err)
	}
	// create zip file for circuit
	tempPath := filepath.Join(".", ".test")

	// If the tempPath exists, remove it
	if _, err := os.Stat(tempPath); err == nil {
		os.RemoveAll(tempPath)
	}

	// TODO you can replace ioutil.TempFile.
	os.MkdirAll(tempPath, os.ModePerm)
	zipFilePath := filepath.Join(tempPath, "circuit.zip")
	defer os.RemoveAll(tempPath)

	sourceDir := filepath.Join(packageDirBase, "circuit")

	helpers.Zip(sourceDir, zipFilePath)

	// retrive SAS key to the backend.
	// TODO after creating the backend, call it.
	// Backend Service return whole infromation for requires to the upload. (download is fine.)
	token, err := repository.GetRepositoryAccessToken()
	if err != nil {
		log.Fatalf("Can not get repository access token %v\n", err)
	}
	circuitBlobName := manifest.Name + "/" + manifest.Version + "/circuit/" + "circuit.zip"

	// upload both zips to blob storage
	circuitBlockBlobURL := helpers.NewBlockBlobWithSASQueryParameter(token.StorageAccountName, token.ContainerName, circuitBlobName, token.SASQueryParameter)
	circuitBlockBlobURL.Upload(zipFilePath)

	packageBlobName := manifest.Name + "/" + manifest.Version + "/package/" + "package.zip"

	// TODO the zip file name could be change. You need to search the directory and iterate the upload.
	packageFilePath := filepath.Join(packageDirBase, "package", "hello.zip")
	packageBlockBlobURL := helpers.NewBlockBlobWithSASQueryParameter(token.StorageAccountName, token.ContainerName, packageBlobName, token.SASQueryParameter)
	packageBlockBlobURL.Upload(packageFilePath)
	// rest api: create (or update) record

	// successfully uploaded

	return nil
}
