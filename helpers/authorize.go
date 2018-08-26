package helpers

import (
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const AZURE_AUTH_LOCATION = "AZURE_AUTH_LOCATION"

func NewAuthorizer(configFilePath string) (autorest.Authorizer, error) {
	os.Setenv(AZURE_AUTH_LOCATION, configFilePath)
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Fails to get Authorizer from the %s .", configFilePath)
		return nil, err
	}
	return authorizer, nil
}
