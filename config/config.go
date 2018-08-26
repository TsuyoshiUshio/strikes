package config

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
