package storage

import (
	"context"
	"log"
	"os"

	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/TsuyoshiUshio/strikes/helpers"
)

func CreateStorageAccountIfNotExists(authorizer autorest.Authorizer, name string, resourceGroup string, location string) (*[]storageAccount.AccountKey, error) {
	authInfo, err := helpers.ReadJson(os.Getenv(helpers.AZURE_AUTH_LOCATION))
	if err != nil {
		return nil, err
	}
	client := storageAccount.NewAccountsClient((*authInfo)["subscriptionId"])
	client.Authorizer = authorizer
	ctx := context.Background()
	result, err := client.CheckNameAvailability(
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
	account, err := client.Create(
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

	err = account.WaitForCompletion(ctx, client.Client)
	if err != nil {
		log.Fatal("Can not get the storage account response")
		return nil, err
	}

	keysResult, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		log.Fatal("Can not fetch list keys: resourceGroup: %s StorageAccount: %s \n", resourceGroup, name)
		return nil, err
	}

	return keysResult.Keys, nil
	// createTableIfNotExists(DEFAULT_STORAGE_TABLE_NAME, name, *accessKeys[0].Value)
}

func CreateTableIfNotExists(tableName, storageAccountName, accessKey string) error {
	client, err := newStorageTableClient(storageAccountName, accessKey)
	if err != nil {
		log.Fatal("Storage Table Client %s can not create: %v\n", storageAccountName, err)
		return err
	}

	tableService := client.GetTableService()
	table := tableService.GetTableReference(tableName)

	err = table.Create(30, storage.EmptyPayload, nil)
	if err != nil {
		log.Fatal("%s: %v", "Table Creating error", err)
		return err
	}
	return nil
}

func newStorageTableClient(name, key string) (*storage.Client, error) {
	client, err := storage.NewBasicClient(name, key)
	if err != nil {
		return nil, err
	}
	return &client, err
}
