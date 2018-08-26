package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Azure/go-autorest/autorest"
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
	err = createDefaultStorageAccountWithTable(resourceGroup, location, false)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFileAndDirectory() error {

	configDir, err := getConfigDir()
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
	configFilePath, err := getConfigFilePath()
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

func getConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return filepath.Join(usr.HomeDir, config.CONFIG_DIR), nil
}

func getConfigFilePath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, config.CONFIG_FILE_NAME), nil
}

func getPowerPlantConfigFilePath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, config.POWER_PLANT_CONFIG_FILE_NAME), nil
}

func getAuthorizer() (autorest.Authorizer, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	authorizer, err := helpers.NewAuthorizer(configFilePath)
	if err != nil {
		return nil, err
	}
	return authorizer, nil
}

func createDefaultResourceGroup(location string) (string, error) {
	authorizer, err := getAuthorizer()
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

	// Create Storage Account
	storageAccountName := storage.DEFAULT_STORAGE_ACCOUNT_NAME + helpers.RandomStrings(8)
	// Store the powerplant configuration file

	authorizer, err := getAuthorizer()
	if err != nil {
		return err
	}
	accountKeys, err := storage.CreateStorageAccountIfNotExists(authorizer, storageAccountName, resourceGroup, location)
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
	powerPlantConfigFilePath, err := getPowerPlantConfigFilePath()
	if err != nil {
		return err
	}
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
