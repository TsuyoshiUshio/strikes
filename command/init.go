package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/resources"
	"github.com/TsuyoshiUshio/strikes/services/storage"
	"github.com/urfave/cli"
)

func Initialize(c *cli.Context) error {
	// Check if there is .strikes on your home directory
	fmt.Println("Initializing Strikes...")
	// Create .strikes and copy .config file to the directory.
	location := c.String("l")

	err := createConfigFileAndDirectory()
	if err != nil {
		return err
	}

	// Create a ResourceGroup is not exists
	resourceGroup, err := createDefaultResourceGroup(location)
	if err != nil {
		return err
	}
	// Create a Storage Account if not exists
	// Default Storage AccountName: sastorikes{random English or Numeric}
	force := c.Bool("f")
	err = createDefaultStorageAccountWithTable(resourceGroup, location, force)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFileAndDirectory() error {

	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	// create config dir
	fmt.Printf("create %s if not exists.\n", configDir)
	err = helpers.CreateDirIfNotExist(configDir)
	if err != nil {
		return err
	}
	// Move config file to ~/.strikes/config
	configFilePath, err := config.GetConfigFilePath()
	if err != nil {
		return err
	}
	err = os.Rename(".config", configFilePath)
	if err != nil {
		fmt.Printf("%s file not found.\n", configFilePath)
		fmt.Println("Get ServicePrincipal to execute this command: az ad sp create-for-rbac -n \"Strikes\" --sdk-auth > .config")
		fmt.Println("Then execute strikes init again.")
		return nil
	}
	return nil
}

func createDefaultResourceGroup(location string) (string, error) {
	authorizer, err := config.GetAuthorizer()
	if err != nil {
		return "", err
	}
	defaultResourceGroup := resources.DEFAULT_RESOURCE_GROUP_NAME + "-" + location
	fmt.Printf("Creating ResourceGroup %s...\n", defaultResourceGroup)
	resourceGroup, err := resources.CreateDefaultResourceGroup(authorizer, location)
	if err != nil {
		return "", err
	}
	return resourceGroup, err
}

// CreateStorageAccountIfNotExists(authorizer autorest.Authorizer, name string, resourceGroup string,location string)

func createDefaultStorageAccountWithTable(resourceGroup string, location string, force bool) error {
	// Read powerplant configration file, check if there is existing storage account.
	powerPlantConfigFilePath, err := config.GetPowerPlantConfigFilePath()
	if err != nil {
		return err
	}

	authorizer, err := config.GetAuthorizer()
	if err != nil {
		return err
	}
	storageAccountClient, err := storage.NewStorageAccountClient(&authorizer)
	if err != nil {
		return err
	}

	if helpers.Exists(powerPlantConfigFilePath) {
		if force {
			// Read config file
			config, err := config.GetPowerPlantConfig()
			if err != nil {
				return err
			}
			// Remove Storage Account
			storageAccountClient.DeleteIfExists(config.StorageAccountName, config.ResourceGroup)
			// Remove the config file
			fmt.Printf("Current PowerPlant configration and storage account: %s (ResourceGroup %s) has been removed\n", config.StorageAccountName, config.ResourceGroup)
			err = os.Remove(powerPlantConfigFilePath)
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

	err = ioutil.WriteFile(powerPlantConfigFilePath, powerPlantConfigBody, 0644)
	if err != nil {
		fmt.Printf("Can not write file: %s \n", powerPlantConfigFilePath)
		return err
	}

	// Create Table Storage
	err = storage.CreateTableIfNotExists(storage.DEFAULT_STORAGE_TABLE_NAME, storageAccountName, accountKey)
	if err != nil {
		return err
	}
	return nil
}
