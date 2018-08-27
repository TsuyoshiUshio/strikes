package resources

import (
	"context"
	"log"

	resourceGroup "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

const DEFAULT_RESOURCE_GROUP_NAME = "strikes-storage"

type IGroupsClient interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, parameters resourceGroup.Group) (resourceGroup.Group, error)
}

type ResourceGroupClient struct {
	Client *resourceGroup.GroupsClient
}

func (c *ResourceGroupClient) CreateDefaultResourceGroup(location string) (string, error) {

	ctx := context.Background()

	resourceGroupName := DEFAULT_RESOURCE_GROUP_NAME + "-" + location
	tags := make(map[string]*string)
	tags["AppName"] = to.StringPtr("strikes")
	parameters := resourceGroup.Group{
		Name:     to.StringPtr(resourceGroupName),
		Location: to.StringPtr(location),
		Tags:     tags,
	}
	_, err := c.Client.CreateOrUpdate(ctx, resourceGroupName, parameters)
	if err != nil {
		log.Fatalf("ResourceGroup creation failed. %s\n", resourceGroupName)
		return "", err
	}
	return resourceGroupName, nil
}
