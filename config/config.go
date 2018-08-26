package config

type Config struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
}

type PowerPlantConfig struct {
	ResourceGroup                  string
	StorageAccountName             string
	StorageAccountConnectionString string
}
