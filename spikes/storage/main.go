package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
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

	client := storage.NewAccountsClient((*authInfo)["subscriptionId"])

	client.Authorizer = authorizer

	ctx := context.Background()

	var name = "sastosam01"
	result, err := client.CheckNameAvailability(
		ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(name),
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
		"RemoveSaStosam01",
		name,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Location: to.StringPtr("japaneast"),
			AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{}})
	if err != nil {
		log.Fatal("%s: %v", "storage account creation fail ", err)
	}
	log.Printf(account.Status())
}
