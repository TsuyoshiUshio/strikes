package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"github.com/Azure/go-autorest/autorest"
	"github.com/TsuyoshiUshio/strikes/helpers"
)

const CONFIG_DIR = ".strikes"

type Config struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
}

const CONFIG_FILE_NAME = "config"

type PowerPlantConfig struct {
	ResourceGroup      string `json:"resourceGroup"`
	StorageAccountName string `json:"storageAccountName"`
	StorageAccountKey  string `json:"storageAccountKey"`
}

const POWER_PLANT_CONFIG_FILE_NAME = "powerplant"

type ICOnfigHelper interface {
	GetConfigDir() (string, error)
	GetConfigFilePath() (string, error)
	GetPowerPlantConfigFilePath() (string, error)
	GetPowerPlantConfig() (*PowerPlantConfig, error)
	GetAuthorizer() (autorest.Authorizer, error)
}
type ConfigHelper struct {
}

func NewConfigHelper() *ConfigHelper {
	return &ConfigHelper{}
}

func (h *ConfigHelper) GetConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return filepath.Join(usr.HomeDir, CONFIG_DIR), nil
}

func (h *ConfigHelper) GetConfigFilePath() (string, error) {
	configDir, err := h.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, CONFIG_FILE_NAME), nil
}

func (h *ConfigHelper) GetPowerPlantConfigFilePath() (string, error) {
	configDir, err := h.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, POWER_PLANT_CONFIG_FILE_NAME), nil
}

func (h *ConfigHelper) GetPowerPlantConfig() (*PowerPlantConfig, error) {
	filePath, err := h.GetPowerPlantConfigFilePath()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(filePath)
	var config PowerPlantConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Cannot unmarshal the config file %s \n", filePath)
		return nil, err
	}
	return &config, nil
}

func (h *ConfigHelper) GetAuthorizer() (autorest.Authorizer, error) {
	configFilePath, err := h.GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	authorizer, err := helpers.NewAuthorizer(configFilePath)
	if err != nil {
		return nil, err
	}
	return authorizer, nil
}
