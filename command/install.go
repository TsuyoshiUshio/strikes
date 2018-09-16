package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/TsuyoshiUshio/strikes/services/resources"
	"github.com/TsuyoshiUshio/strikes/services/storage"
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
		return nil, err
	}

	// Download Circuit

	// Generate Unique name for the deployment
	// Configure Circuit by the possible parameters
	// Execute terraform with the parameters
	// Update the Powerplant
}

func InitializeWithConfigContextAndResrouceGroup(c *cli.Context, configContext *config.ConfigContext, resourceGroupName string) error {
	// Check if there is .strikes on your home directory
	fmt.Println("Initializing Strikes...")
	// Create .strikes and copy .config file to the directory.
	location := c.String("l")

	err := createConfigFileAndDirectory(configContext)
	if err != nil {
		return err
	}

	// Create a ResourceGroup is not exists
	resourceGroup, err := createResourceGroup(configContext, resourceGroupName, location)
	if err != nil {
		return err
	}
	// Create a Storage Account if not exists
	// Default Storage AccountName: sastorikes{random English or Numeric}
	force := c.Bool("f")
	err = createDefaultStorageAccountWithTable(configContext, resourceGroup, location, force)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFileAndDirectory(configContext *config.ConfigContext) error {
	// create config dir
	fmt.Printf("create %s if not exists.\n", configContext.ConfigDir)
	err := helpers.CreateDirIfNotExist(configContext.ConfigDir)
	if err != nil {
		return err
	}
	// Move config file to ~/.strikes/config
	configFilePath := configContext.GetConfigFilePath()
	err = os.Rename(".config", configFilePath)
	if err != nil {
		fmt.Printf("%s file not found.\n", configFilePath)
		fmt.Println("Get ServicePrincipal to execute this command: az ad sp create-for-rbac -n \"Strikes\" --sdk-auth > .config")
		fmt.Println("Then execute strikes init again.")
		return nil
	}
	return nil
}

func createResourceGroup(configContext *config.ConfigContext, resourceGroupName string, location string) (string, error) {
	authorizer, err := configContext.GetAuthorizer()
	if err != nil {
		return "", err
	}

	fmt.Printf("Creating ResourceGroup %s...\n", resourceGroupName)
	client, err := resources.NewResrouceGroupClient(authorizer)
	if err != nil {
		return "", err
	}
	resourceGroup, err := client.CreateResourceGroup(resourceGroupName, location)
	if err != nil {
		return "", err
	}
	return resourceGroup, err
}

func createDefaultStorageAccountWithTable(configContext *config.ConfigContext, resourceGroup string, location string, force bool) error {
	// Read powerplant configration file, check if there is existing storage account.
	authorizer, err := configContext.GetAuthorizer()
	if err != nil {
		return err
	}
	storageAccountClient, err := storage.NewStorageAccountClient(&authorizer)
	if err != nil {
		return err
	}

	if helpers.Exists(configContext.PowerPlantConfigFilePath) {
		if force {
			// Read config file
			config, err := configContext.GetPowerPlantConfig()
			if err != nil {
				return err
			}
			// Remove Storage Account
			storageAccountClient.DeleteIfExists(config.StorageAccountName, config.ResourceGroup)
			// Remove the config file
			fmt.Printf("Current PowerPlant configration and storage account: %s (ResourceGroup %s) has been removed\n", config.StorageAccountName, config.ResourceGroup)
			err = os.Remove(configContext.PowerPlantConfigFilePath)
		} else {
			// Do nothing
			fmt.Printf("PowerPlant strage account is already exists. For more details, see {Home}/.strikes/powerplant.\n")
			return nil
		}
	}

	// Create Storage Account
	storageAccountName := storage.DEFAULT_STORAGE_ACCOUNT_NAME + helpers.RandomStrings(8)
	// Store the powerplant configuration file

	accountKeys, err := storageAccountClient.CreateStorageAccountIfNotExists(storageAccountName, resourceGroup, location)
	if err != nil {
		return err
	}
	accountKey := *(*accountKeys)[0].Value
	fmt.Printf("AccountKey: %s\n", accountKey)

	// Write Config File
	powerPlantConfig := &config.PowerPlantConfig{
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		StorageAccountKey:  accountKey,
	}

	powerPlantConfigBody, _ := json.Marshal(powerPlantConfig)

	err = ioutil.WriteFile(configContext.PowerPlantConfigFilePath, powerPlantConfigBody, 0644)
	if err != nil {
		fmt.Printf("Can not write file: %s \n", configContext.PowerPlantConfigFilePath)
		return err
	}

	// Create Table Storage
	err = storage.CreateTableIfNotExists(storage.DEFAULT_STORAGE_TABLE_NAME, storageAccountName, accountKey)
	if err != nil {
		return err
	}
	return nil
}
