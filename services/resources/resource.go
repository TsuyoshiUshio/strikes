package resources

import (
	"context"
	"log"
	"os"

	resourceGroup "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/TsuyoshiUshio/strikes/helpers"
)

func CreateDefaultResourceGroup(authorizer autorest.Authorizer, location string) error {
	authInfo, err := helpers.ReadJson(os.Getenv(helpers.AZURE_AUTH_LOCATION))
	if err != nil {
		return err
	}
	client := resourceGroup.NewGroupsClient((*authInfo)["subscriptionId"])
	client.Authorizer = authorizer
	ctx := context.Background()

	resourceGroupName := DEFAULT_RESOURCE_GROUP_NAME + "-" + location
	tags := make(map[string]*string)
	tags["AppName"] = to.StringPtr("strikes")
	parameters := resourceGroup.Group{
		Name:     to.StringPtr(resourceGroupName),
		Location: to.StringPtr(location),
		Tags:     tags,
	}
	_, err = client.CreateOrUpdate(ctx, resourceGroupName, parameters)
	if err != nil {
		log.Fatalf("ResourceGroup creation failed. %s\n", resourceGroupName)
		return err
	}
	return nil
}
