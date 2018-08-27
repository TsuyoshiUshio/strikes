package storage

import (
	"log"
	"os"

	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/TsuyoshiUshio/strikes/helpers"
)

func NewStorageAccountClient(authorizer *autorest.Authorizer) (*StorageAccountClient, error) {
	authInfo, err := helpers.ReadJson(os.Getenv(helpers.AZURE_AUTH_LOCATION))
	if err != nil {
		return nil, err
	}
	storageAccountClient := storageAccount.NewAccountsClient((*authInfo)["subscriptionId"])
	storageAccountClient.Authorizer = *authorizer
	return &StorageAccountClient{
		Client:            &storageAccountClient,
		AutoRestClient:    &(storageAccountClient.Client),
		WaitForCompletion: &waitForCompletionImpl{},
	}, nil
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
