package integaration

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/TsuyoshiUshio/strikes/command"
	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/services/resources"
	"github.com/TsuyoshiUshio/strikes/services/storage"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

const STRIKES_CONFIG_ENVIRONMENT_VARIABLE = "STRIKES_CONFIG"

func getHomeDir() (string, error) {
	return getConfigDir(), nil
}

func getConfigDir() string {
	return filepath.Join(".", ".strikes")
}

// NOTE: This secnario test should execute one by one.
func TestInitializeCommand_NormalCase(t *testing.T) {
	// NOTE: configure Environment Variables STRIKES_CONFIG
	// e.g. export STRIKES_CONFIG=`cat ~/.strikes/config`

	// Write config file from environment variables
	content := os.Getenv(STRIKES_CONFIG_ENVIRONMENT_VARIABLE)
	dotConfigDir := filepath.Join(".", ".config")
	err := ioutil.WriteFile(dotConfigDir, []byte(content), 0644)
	if err != nil {
		assert.Failf(t, "Can not open temp file", err.Error())
		return
	}

	configDir := getConfigDir()
	powerPlantConfigFilePath := filepath.Join(configDir, "powerplant")

	// Delete configDir if exists
	helpers.DeleteDirIfExists(configDir)

	// Create ConfigContext object
	configContext := config.ConfigContext{
		ConfigDir:                configDir,
		PowerPlantConfigFilePath: powerPlantConfigFilePath,
		GetHomeDir:               getHomeDir,
	}

	set := flag.NewFlagSet("test", 0)
	set.String("l", "japaneast", "doc")
	c := cli.NewContext(nil, set, nil)
	// Execute Initialize Command.
	err = command.InitializeWithConfigContextAndResrouceGroup(
		c,
		&configContext,
		"strikes-integration-test")

	powerplangConfig, err := configContext.GetPowerPlantConfig()
	theFirstKey := powerplangConfig.StorageAccountKey
	// This must be skipped.
	err = command.InitializeWithConfigContextAndResrouceGroup(
		c,
		&configContext,
		"strikes-integration-test")

	result, err := storage.CheckTableExists(storage.DEFAULT_STORAGE_TABLE_NAME, powerplangConfig.StorageAccountName, powerplangConfig.StorageAccountKey)
	assert.Equal(t, true, result, "Storage Table doesn't exist")
	assert.Nil(t, err)
	powerplangConfig, err = configContext.GetPowerPlantConfig()
	// It should be the same.
	assert.Equal(t, theFirstKey, powerplangConfig.StorageAccountKey)
	// Check if there is config file and it creates Storage Account with Table on Azure

	// Force update scenario
	set.Bool("f", true, "doc")
	c = cli.NewContext(nil, set, nil)
	err = command.InitializeWithConfigContextAndResrouceGroup(
		c,
		&configContext,
		"strikes-integration-test")
	powerplangConfig, err = configContext.GetPowerPlantConfig()
	// New key should be generated.
	assert.NotEqual(t, theFirstKey, powerplangConfig.StorageAccountKey)

	// Clean up section
	// Remove the resource Group
	authorizer, err := configContext.GetAuthorizer()
	assert.Nil(t, err)

	resourceClient, err := resources.NewResrouceGroupClient(authorizer)
	assert.Nil(t, err)
	resourceClient.Delete(powerplangConfig.ResourceGroup)

	// Remove the config directory
	fmt.Printf("Removing.... %s \n", configContext.ConfigDir)
	err = os.RemoveAll(configContext.ConfigDir)

	assert.Nil(t, err)
}
