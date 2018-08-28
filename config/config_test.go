package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHomeDir(t *testing.T) {
	homeDir, _ := getHomeDir()
	usr, _ := user.Current()

	assert.Equal(t, usr.HomeDir, homeDir)
}

func TestGetPowerPlantConfig(t *testing.T) {
	content := `{"resourceGroup":"strikes-storage-japaneast","storageAccountName":"sastrikesdijjeqcx","storageAccountKey":"SomeKeys=="}`
	file, err := ioutil.TempFile(".", "strikes")
	defer os.Remove(file.Name())
	if err != nil {
		assert.Failf(t, "Can not open temp file", err.Error())
		return
	}
	_, err = file.WriteString(content)
	context := ConfigContext{
		PowerPlantConfigFilePath: file.Name(),
	}
	config, err := context.GetPowerPlantConfig()
	assert.Equal(t, "strikes-storage-japaneast", config.ResourceGroup)
	assert.Equal(t, "sastrikesdijjeqcx", config.StorageAccountName)
	assert.Equal(t, "SomeKeys==", config.StorageAccountKey)

}
