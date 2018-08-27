package storage

import storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"

const DEFAULT_STORAGE_TABLE_NAME = "powerplantstatus"
const DEFAULT_STORAGE_ACCOUNT_NAME = "sastrikes"

type StorageAccountClient struct {
	Client *storageAccount.AccountsClient
}
