package resources

import (
	"os"

	resourceGroup "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/TsuyoshiUshio/strikes/helpers"
)

func NewResrouceGroupClient(authorizer autorest.Authorizer) (*ResourceGroupClient, error) {
	authInfo, err := helpers.ReadJson(os.Getenv(helpers.AZURE_AUTH_LOCATION))
	if err != nil {
		return nil, err
	}
	client := resourceGroup.NewGroupsClient((*authInfo)["subscriptionId"])
	client.Authorizer = authorizer
	resourceGroupClient := ResourceGroupClient{
		Client: &client,
	}
	return &resourceGroupClient, nil
}
