package command

import (
	"log"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/providers"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/TsuyoshiUshio/strikes/services/storage"
	"github.com/urfave/cli"
)

// 	defaultResourceGroup := resources.DEFAULT_RESOURCE_GROUP_NAME + "-" + location

type InstallCommand struct {
}

func (s *InstallCommand) Install(c *cli.Context) error {
	// Get the package Name from the parameter
	packageName := c.Args().Get(0)

	// Get the instance name from the parameter
	instanceName := c.Args().Get(1)

	// Get Metadata from Backend API
	p, err := repository.GetPackage(packageName)
	if err != nil {
		log.Fatalf("Can not load the package name: %v\n", err)
		return err
	}

	setUpStrikesTemp()
	// Download Circuit
	zipFilePath := filepath.Join(STRIKES_TEMP, "circuit.zip")
	targetDirPath := filepath.Join(STRIKES_TEMP, "circuit")
	err = helpers.DownloadFile(zipFilePath, p.GetCircuitZipURL())
	if err != nil {
		log.Fatalf("Can not download cricuit zip file.: %v\n", err)
		return err
	}

	helpers.UnZip(zipFilePath, targetDirPath)
	manifest, err := config.NewManifestFromFile(targetDirPath) // TODO after developing Provider, _ should be
	if err != nil {
		log.Fatalf("Can not read manifest file from the download contents. :%v\n", err)
		return err
	}

	// Execute deployment using Provider.
	provider := providers.NewTerraformProvider(manifest, targetDirPath) //targetDirPath is here or adding one deep directory
	result := provider.CreateResource(c.Args().Tail())                  // The first one is the package name.

	// Update the PowerPlant
	instance := storage.StrikesInstance{
		PackageID:         p.Id,
		Name:              instanceName,
		ResourceGroup:     result.GetResourceGroup(),
		PackageName:       p.Name,
		PackageVersion:    p.LatestVersion(), // TODO: The version should be changed by the parameter
		PackageParameters: result.GetConfigrationsJosn(),
	}
	err = storage.InsertOrUpdate(&instance)
	if err != nil {
		log.Fatalf("Can not insert strikes instance to the PowerPlant. %v", err)
	}

	return nil
}

const STRIKES_TEMP = ".strikesTemp"

func setUpStrikesTemp() {
	err := helpers.DeleteDirIfExists(STRIKES_TEMP)
	if err != nil {
		log.Fatalf("Can not delete .strikesTemp. as set up : %v\n", err)
		return
	}
	err = helpers.CreateDirIfNotExist(STRIKES_TEMP)
	if err != nil {
		log.Fatalf("Can not create .strikesTemp. as set up: %v\n", err)
		return
	}
}

func cleanUpStrikesTemp() {
	err := helpers.DeleteDirIfExists(STRIKES_TEMP)
	if err != nil {
		log.Fatalf("Can not delete .strikesTemp. as clean up : %v\n", err)
		return
	}
}
