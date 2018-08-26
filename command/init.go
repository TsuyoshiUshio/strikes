package command

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/urfave/cli"
)

func Initialize(c *cli.Context) error {
	// Check if there is .strikes on your home directory
	fmt.Println("Initializing Strikes...")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(usr.HomeDir, ".strikes")
	fmt.Printf("create %s if not exists.\n", configDir)
	err = helpers.CreateDirIfNotExist(configDir)
	// Move config file to ~/.strikes/config
	configFilePath := filepath.Join(configDir, "config")
	err = os.Rename(".config", configFilePath)
	if err != nil {
		fmt.Printf("%s file not found.\n", configFilePath)
		fmt.Println("Get ServicePrincipal to execute this command: az ad sp create-for-rbac -n \"Strikes\" --sdk-auth > .config")
		fmt.Println("Then execute strikes init again.")
		return nil
	}
	// Create a configuration file
	// Create a ResourceGroup is not exists
	// Default ResourceGroupName: strikes-storage-japaneast
	// Create a Storage Account if not exists
	// Default Storage AccountName: sastorikes{random English or Numeric}
	return err
}
