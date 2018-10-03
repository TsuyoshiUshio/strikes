package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	createPackageBlockBlob(token, manifest, packageDirBase)

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

func createPackageBlockBlob(token *repository.RepositoryAccessToken, manifest *config.Manifest, packageDirBase string) {

	packageDirPath := filepath.Join(packageDirBase, "pacakge")
	files, err := ioutil.ReadDir(packageDirPath)
	if err != nil {
		log.Fatalf("Can not read zip file for the target directory: %v\n", err)
	}
	fileNames := mapFileInfo(files, func(file os.FileInfo) string {
		return file.Name()
	})
	zipFileNames := helpers.Filter(fileNames, func(s string) bool {
		if strings.HasSuffix(s, ".zip") {
			return true
		} else {
			return false
		}
	})

	for _, zipFileName := range zipFileNames {
		packageBlobName := getPackageBlobName(manifest, zipFileName)
		packageFilePath := filepath.Join(packageDirBase, "package", zipFileName)
		packageBlockBlobURL := helpers.NewBlockBlobWithSASQueryParameter(token.StorageAccountName, token.ContainerName, packageBlobName, token.SASQueryParameter)
		packageBlockBlobURL.Upload(packageFilePath)
	}
}

func mapFileInfo(files []os.FileInfo, f func(file os.FileInfo) string) []string {
	result := make([]string, len(files))
	for _, file := range files {
		result = append(result, f(file))
	}
	return result
}

func getPackageBlobName(manifest *config.Manifest, zipFileName string) string {
	return manifest.Name + "/" + manifest.Version + "/package/" + zipFileName
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
