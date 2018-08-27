package storage

import (
	"context"
	"log"

	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

const DEFAULT_STORAGE_TABLE_NAME = "powerplantstatus"
const DEFAULT_STORAGE_ACCOUNT_NAME = "sastrikes"

type IAccountsClient interface {
	CheckNameAvailability(ctx context.Context, accountName storageAccount.AccountCheckNameAvailabilityParameters) (result storageAccount.CheckNameAvailabilityResult, err error)
	Create(ctx context.Context, resourceGroupName string, accountName string, parameters storageAccount.AccountCreateParameters) (result storageAccount.AccountsCreateFuture, err error)
	GetProperties(ctx context.Context, resourceGroupName string, accountName string) (result storageAccount.Account, err error)
	Delete(ctx context.Context, resourceGroupName string, accountName string) (result autorest.Response, err error)
	ListKeys(ctx context.Context, resourceGroupName string, accountName string) (result storageAccount.AccountListKeysResult, err error)
}

type IWaitForCompletion interface {
	Wait(account storageAccount.AccountsCreateFuture, ctx context.Context, autoRestClient *autorest.Client) error
}

type waitForCompletionImpl struct{}

func (w *waitForCompletionImpl) Wait(account storageAccount.AccountsCreateFuture, ctx context.Context, autoRestClient *autorest.Client) error {
	return account.WaitForCompletion(ctx, *autoRestClient)
}

type StorageAccountClient struct {
	Client            IAccountsClient
	AutoRestClient    *autorest.Client
	WaitForCompletion IWaitForCompletion
}

func (c *StorageAccountClient) CreateStorageAccountIfNotExists(name string, resourceGroup string, location string) (*[]storageAccount.AccountKey, error) {
	ctx := context.Background()
	result, err := c.Client.CheckNameAvailability(
		ctx,
		storageAccount.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(name),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		},
	)

	if err != nil {
		log.Fatalf("storage account %s creation failed. error: %v \n", name, err)
		return nil, err
	}
	if *result.NameAvailable != true {
		log.Fatalf("storage account %s is not available. error: %v \n", name, err)
		return nil, err
	}
	account, err := c.Client.Create(
		ctx,
		resourceGroup,
		name,
		storageAccount.AccountCreateParameters{
			Sku: &storageAccount.Sku{
				Name: storageAccount.StandardLRS,
			},
			Location: to.StringPtr(location),
			AccountPropertiesCreateParameters: &storageAccount.AccountPropertiesCreateParameters{},
		},
	)

	if err != nil {
		log.Fatal("Storage Creation fail: %s\n", name)
		return nil, err
	}

	err = c.WaitForCompletion.Wait(account, ctx, c.AutoRestClient)
	if err != nil {
		log.Fatal("Can not get the storage account response")
		return nil, err
	}

	keysResult, err := c.Client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Fatal("Can not fetch list keys: resourceGroup: %s StorageAccount: %s \n", resourceGroup, name)
		return nil, err
	}

	return keysResult.Keys, nil
}

func (c *StorageAccountClient) DeleteIfExists(accountName string, resourceGroupName string) error {
	ctx := context.Background()
	_, err := c.Client.GetProperties(ctx, resourceGroupName, accountName)
	// TODO this implementation is not ideal. We can find much safer way to check if it exists.
	if err == nil {
		_, err := c.Client.Delete(ctx, resourceGroupName, accountName)
		if err != nil {
			log.Fatalf("Can not delete the resource group: %s storage account: %s Error: %v\n", resourceGroupName, accountName, err)
			return err
		}
	}
	return nil
}
