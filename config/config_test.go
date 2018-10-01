package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/bouk/monkey"
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

func TestGetConfig(t *testing.T) {
	path := filepath.Join(".", "test-fixture", "config-basic", "config")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	ActualFileName := ""
	fakeReadFile := func(fileName string) ([]byte, error) {
		ActualFileName = fileName
		return file, nil
	}
	monkey.Patch(ioutil.ReadFile, fakeReadFile)
	defer monkey.UnpatchAll()
	getHomeDir := func() (string, error) {
		return "baz", nil
	}
	context := ConfigContext{
		ConfigDir:                "foo",
		PowerPlantConfigFilePath: "bar",
		GetHomeDir:               getHomeDir,
	}
	config, err := context.GetConfig()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "foo", config.ClientID)
	assert.Equal(t, "bar", config.ClientSecret)
	assert.Equal(t, "baz", config.SubscriptionID)
	assert.Equal(t, "qux", config.TenantID)
}
