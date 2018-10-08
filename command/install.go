package command

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/providers"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/TsuyoshiUshio/strikes/services/storage"
	"github.com/rs/xid"
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

	fmt.Printf("packageName: %s\n instanceName: %s\n", packageName, instanceName)

	// if the packageName directory exists, it will be local install. Not install from repository.

	p, err := retrivePackageFromRepository(packageName)
	if err != nil {
		return err
	}
	targetDirPath := filepath.Join(STRIKES_TEMP, "circuit")

	manifestFilePath := filepath.Join(targetDirPath, "manifest.yaml")
	manifest, err := config.NewManifestFromFile(manifestFilePath) // TODO after developing Provider, _ should be
	if err != nil {
		log.Fatalf("Can not read manifest file from the download contents. :%v\n", err)
		return err
	}

	// Execute deployment using Provider.
	provider := providers.NewTerraformProvider(manifest, targetDirPath) //targetDirPath is here or adding one deep directory
	result := provider.CreateResource(c.Args().Tail(), instanceName)    // The first one is the package name.

	// Update the PowerPlant
	instance := storage.StrikesInstance{
		InstanceID:        xid.New().String(), // Automatically generated xid sortable. more detail https://github.com/rs/xid
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

func retrivePackageFromRepository(packageName string) (*repository.Package, error) {
	// Get Metadata from Backend API
	p, err := repository.GetPackage(packageName)
	if err != nil {
		log.Printf("[DEBUG] GetPackage Error: %v,\n", err)
		log.Fatalf("Can not find package: %s \n", packageName)
		return nil, err
	}

	setUpStrikesTemp()
	// Download Circuits
	zipFilePath := filepath.Join(STRIKES_TEMP, "circuit.zip")
	err = helpers.DownloadFile(zipFilePath, p.GetCircuitZipURL())
	if err != nil {
		log.Fatalf("Can not download cricuit zip file.: %v\n", err)
		return nil, err
	}

	err = helpers.UnZip(zipFilePath, STRIKES_TEMP)
	if err != nil {
		log.Printf("[DEBUG] Extract Zip Error.: %v\n", err)
		log.Fatalf("Can not extract the Zip file.: %v\n", zipFilePath)
		return nil, err
	}
	return p, nil
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
