package storage

import (
	"log"
	"os"
	"time"

	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
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
	table, err := GetTable(tableName, storageAccountName, accessKey)
	if err != nil {
		log.Fatalf("Can not get table reference tableName: %s storageAccount: %s, %v", tableName, storageAccountName, err)
		return err
	}

	err = table.Create(30, storage.EmptyPayload, nil)
	if err != nil {
		log.Fatalf("%s: %v", "Table Creating error", err)
		return err
	}
	return nil
}

func GetContainer(storageAccountName, accessKey, containerName string) (*storage.Container, error) {
	client, err := newStorageAccountClient(storageAccountName, accessKey)
	if err != nil {
		return nil, err
	}
	blobClient := client.GetBlobService()
	return blobClient.GetContainerReference(containerName), nil
}

func CreateContainerIfNotExists(storageAccountName, accessKey, containerName string) error {
	container, err := GetContainer(storageAccountName, accessKey, containerName)
	options := storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypeBlob,
	}
	_, err = container.CreateIfNotExists(&options)
	if err != nil {
		return err
	}
	return nil
}

func FetchSASQueryParameters(accountName, accountKey, containerName string) string {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(err)
	}
	// for the go lang
	sasQueryParams := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		ContainerName: containerName,
		BlobName:      "",
		Permissions:   azblob.ContainerSASPermissions{Add: true, Read: true, Write: true}.String(),
	}.NewSASQueryParameters(credential)
	qb := sasQueryParams.Encode()
	return qb
}

func GetTable(tableName, storageAccountName, accessKey string) (*storage.Table, error) {
	client, err := newStorageAccountClient(storageAccountName, accessKey)
	if err != nil {
		log.Fatalf("Storage Table Client %s can not create: %v\n", storageAccountName, err)
		return nil, err
	}

	tableService := client.GetTableService()
	return tableService.GetTableReference(tableName), nil
}

func CheckTableExists(tableName, storageAccountName, accessKey string) (bool, error) {
	table, err := GetTable(tableName, storageAccountName, accessKey)
	if err != nil {
		return false, nil
	}
	err = table.Get(30, storage.FullMetadata)
	if err != nil {
		return false, err
	}
	if tableName == table.Name {
		return true, nil
	} else {
		return false, nil
	}
}

func newStorageAccountClient(name, key string) (*storage.Client, error) {
	client, err := storage.NewBasicClient(name, key)
	if err != nil {
		return nil, err
	}
	return &client, err
}
