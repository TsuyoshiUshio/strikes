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

type IConfigContext interface {
	GetConfigDir() (string, error)
	GetConfigFilePath() string
	GetPowerPlantConfigFilePath() string
	GetPowerPlantConfig() (*PowerPlantConfig, error)
	GetAuthorizer() (autorest.Authorizer, error)
}

type IGetHomeDir func() (string, error)

type ConfigContext struct {
	ConfigDir                string
	PowerPlantConfigFilePath string
	GetHomeDir               IGetHomeDir
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return usr.HomeDir, nil
}

func NewConfigContext() (*ConfigContext, error) {
	return NewConfigContextWithGetHomeDir(getHomeDir)
}

func NewConfigContextWithGetHomeDir(getHomeDirFunc IGetHomeDir) (*ConfigContext, error) {
	context := ConfigContext{}
	context.GetHomeDir = getHomeDirFunc

	configDir, err := context.GetConfigDir()
	if err != nil {
		return nil, err
	}

	context.ConfigDir = configDir
	context.PowerPlantConfigFilePath = context.GetPowerPlantConfigFilePath()
	return &context, nil
}

func (c *ConfigContext) GetConfigDir() (string, error) {
	homeDir, err := c.GetHomeDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return filepath.Join(homeDir, CONFIG_DIR), nil
}

func (c *ConfigContext) GetConfigFilePath() string {
	return filepath.Join(c.ConfigDir, CONFIG_FILE_NAME)
}

func (c *ConfigContext) GetPowerPlantConfigFilePath() string {
	return filepath.Join(c.ConfigDir, POWER_PLANT_CONFIG_FILE_NAME)
}

func (c *ConfigContext) GetPowerPlantConfig() (*PowerPlantConfig, error) {
	content, err := ioutil.ReadFile(c.PowerPlantConfigFilePath)
	var config PowerPlantConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Cannot unmarshal the config file %s \n", c.PowerPlantConfigFilePath)
		return nil, err
	}
	return &config, nil
}

func (c *ConfigContext) GetAuthorizer() (autorest.Authorizer, error) {
	authorizer, err := helpers.NewAuthorizer(c.GetConfigFilePath())
	if err != nil {
		return nil, err
	}
	return authorizer, nil
}
