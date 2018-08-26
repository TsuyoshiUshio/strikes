package command

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/resources"
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
	err = createDefaultResourceGroup(location)
	if err != nil {
		return err
	}
	// Default ResourceGroupName: strikes-storage-japaneast
	// Create a Storage Account if not exists
	// Default Storage AccountName: sastorikes{random English or Numeric}
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
	return filepath.Join(usr.HomeDir, ".strikes"), nil
}

func getConfigFilePath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config"), nil
}

func createDefaultResourceGroup(location string) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	authorizer, err := helpers.NewAuthorizer(configFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("Creating ResourceGroup %s...\n", resources.DEFAULT_RESOURCE_GROUP_NAME+"-"+location)
	err = resources.CreateDefaultResourceGroup(authorizer, location)
	if err != nil {
		return err
	}
	return nil
}
