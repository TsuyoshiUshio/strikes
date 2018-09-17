package command

import (
	"log"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/urfave/cli"
)

// 	defaultResourceGroup := resources.DEFAULT_RESOURCE_GROUP_NAME + "-" + location

type InstallCommand struct {
}

func (s *InstallCommand) Install(c *cli.Context) error {
	// Get the package Name from the parameter
	packageName := c.Args().Get(0)

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
	_, err = config.NewManifestFromFile(targetDirPath) // TODO after developing Provider, _ should be
	if err != nil {
		log.Fatalf("Can not read manifest file from the download contents. :%v\n", err)
		return err
	}

	// Execute deployment using Provider.
	// Update the PowerPlant
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
