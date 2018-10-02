package command

import (
	"fmt"
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

	packageDirBasePath := c.String("p")
	packageDirBase, _ := filepath.Abs(packageDirBasePath)
	log.Printf("[DEBUG] packageDirBase: %s\n", packageDirBase)
	// Load manifest file
	manifest, err := loadManifestFile(packageDirBase)

	// Create a temp directory
	tempPath := filepath.Join(".", ".test")
	createTempDir(tempPath)
	defer os.RemoveAll(tempPath)

	// Compress source dir to zip file.
	zipFilePath := filepath.Join(tempPath, "circuit.zip")
	sourceDir := filepath.Join(packageDirBase, "circuit")
	log.Printf("[DEUBG] zipFilePath: %s\nsourceDir: %s\n", zipFilePath, sourceDir)

	helpers.Zip(sourceDir, zipFilePath)

	// Get AccessToken from BackendAPI
	token, err := repository.GetRepositoryAccessToken()
	if err != nil {
		log.Fatalf("Can not get repository access token %v\n", err)
	}

	createCircuitBlockBlob(token, manifest, zipFilePath)

	createPackageBlockBlob(token, manifest, packageDirBase, zipFilePath)

	cretePackageToBackendAPI(manifest)

	fmt.Printf("Package %s successfully pushed.\n", manifest.Name)

	return nil
}

func loadManifestFile(packageDirBase string) (*config.Manifest, error) {
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
	return manifest, nil
}

func createTempDir(tempPath string) {
	// If the tempPath exists, remove it
	if _, err := os.Stat(tempPath); err == nil {
		os.RemoveAll(tempPath)
	}

	// TODO you can replace ioutil.TempFile.
	os.MkdirAll(tempPath, os.ModePerm)
}

func createCircuitBlockBlob(token *repository.RepositoryAccessToken, manifest *config.Manifest, zipFilePath string) {
	circuitBlobName := getCircuitBlobName(manifest)

	// upload both zips to blob storage
	circuitBlockBlobURL := helpers.NewBlockBlobWithSASQueryParameter(token.StorageAccountName, token.ContainerName, circuitBlobName, token.SASQueryParameter)
	circuitBlockBlobURL.Upload(zipFilePath)
}

func getCircuitBlobName(manifest *config.Manifest) string {
	return manifest.Name + "/" + manifest.Version + "/circuit/" + "circuit.zip"
}

func createPackageBlockBlob(token *repository.RepositoryAccessToken, manifest *config.Manifest, packageDirBase string, zipFilePath string) {
	packageBlobName := getPackageBlobName(manifest)

	// TODO the zip file name could be change. You need to search the directory and iterate the upload.
	packageFilePath := filepath.Join(packageDirBase, "package", "hello.zip")
	packageBlockBlobURL := helpers.NewBlockBlobWithSASQueryParameter(token.StorageAccountName, token.ContainerName, packageBlobName, token.SASQueryParameter)
	packageBlockBlobURL.Upload(packageFilePath)
}

func getPackageBlobName(manifest *config.Manifest) string {
	return manifest.Name + "/" + manifest.Version + "/package/" + "package.zip"
}

func cretePackageToBackendAPI(manifest *config.Manifest) {

	pkg := repository.NewPackageWithCurrentTime(
		manifest.Name,
		manifest.Description,
		manifest.Author,
		manifest.ProjectPage,
		manifest.ProjectRepo,
		manifest.Version,
		manifest.ReleaseNote,
		manifest.ProviderType)
	_, err := pkg.Create()
	if err != nil {
		log.Fatalf("Can not create package: %v\n", err)
	}
}
