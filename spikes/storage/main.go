package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	resourceGroup "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/azure-sdk-for-go/storage"
)

func readJson(fileName string) (*map[string]string, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can not open file %s: %v", fileName, err)
	}
	defer jsonFile.Close()
	jsonByte, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	err = json.Unmarshal(jsonByte, &result)
	return &result, err
}

func main() {
	//	TODO currently, it requres Resource group.
	//  Please create .config file from az ad sp create-for-rbac -n "Strikes" --sdk-auth
	//  Other problem is dep doesn't seem work for azure-sdk-for-go
	os.Setenv("AZURE_AUTH_LOCATION", "./.config")
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Failed to get ")
	}
	authInfo, err := readJson(os.Getenv("AZURE_AUTH_LOCATION"))

	// create or update a resource group.
	resourceGroupClient := resourceGroup.NewGroupsClient((*authInfo)["subscriptionId"])
	resourceGroupClient.Authorizer = authorizer
	ctx := context.Background()

	resourceGroupName := "RemoveSaStosam02"
	tags := make(map[string]*string)
	tags["AppName"] = to.StringPtr("strikes")
	parameters := resourceGroup.Group{
		Name:     to.StringPtr("RemoveSaStosam02"),
		Location: to.StringPtr("japaneast"),
		Tags:     tags,
	}
	_, err = resourceGroupClient.CreateOrUpdate(ctx, resourceGroupName, parameters)
	if err != nil {
		log.Fatalf("%s: %v", "resource group creation failed", err)
	}

	client := storageAccount.NewAccountsClient((*authInfo)["subscriptionId"])

	client.Authorizer = authorizer

	var storageAccountName = "sastosam02"
	result, err := client.CheckNameAvailability(
		ctx,
		storageAccount.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(storageAccountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		log.Fatalf("%s: %v", "storage account creation failed", err)
	}
	if *result.NameAvailable != true {
		log.Fatalf("%s: %v", "storage account name not available", err)
	}

	account, err := client.Create(
		ctx,
		resourceGroupName,
		storageAccountName,
		storageAccount.AccountCreateParameters{
			Sku: &storageAccount.Sku{
				Name: storageAccount.StandardLRS},
			Location: to.StringPtr("japaneast"),
			AccountPropertiesCreateParameters: &storageAccount.AccountPropertiesCreateParameters{}})
	if err != nil {
		log.Fatal("%s: %v", "storage account creation fail ", err)
	}

	err = account.WaitForCompletion(ctx, client.Client)
	if err != nil {
		log.Fatal("Cannot get the storage account future response")
	}

	log.Printf(account.Status())

	// Retrive the the access key

	keysResult, err := client.ListKeys(ctx, resourceGroupName, storageAccountName)
	accessKey := *(keysResult.Keys)
	tableName := "someTable"
	execBasicTableOperation(tableName, storageAccountName, *accessKey[0].Value)

}

func execBasicTableOperation(tableName, storageAccountName, accessKey string) {
	client, err := newStorageTableClient(storageAccountName, accessKey)
	if err != nil {
		log.Fatal("%s: %v", "Table Client can not create", err)
	}

	// CreateOrUpdate the Table.
	// the API for table storage is going to deplicated.
	// Currently, I try to create and it fails if it already has an Table Storage.
	tableService := client.GetTableService()
	table := tableService.GetTableReference(tableName)
	err = table.Create(30, storage.EmptyPayload, nil)
	if err != nil {
		log.Fatal("%s: %v", "Table Creating error", err)
	}

}

func newStorageTableClient(name, key string) (*storage.Client, error) {
	client, err := storage.NewBasicClient(name, key)
	if err != nil {
		return nil, err
	}
	return &client, err
}
