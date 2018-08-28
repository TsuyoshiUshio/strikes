package storage

import (
	"errors"
	"testing"

	storageAccount "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	mock "github.com/TsuyoshiUshio/strikes/services/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateStorageAccountIfNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accountsClient := mock.NewMockIAccountsClient(ctrl)
	waitForCompletion := mock.NewMockIWaitForCompletion(ctrl)

	storageAccountClient := &StorageAccountClient{
		Client:            accountsClient,
		AutoRestClient:    nil,
		WaitForCompletion: waitForCompletion,
	}

	storageAccountName := "someAccount"
	resourceGroup := "someGroup"
	location := "japaneast"
	accountsClient.EXPECT().CheckNameAvailability(gomock.Any(),
		storageAccount.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(storageAccountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		}).Return(
		storageAccount.CheckNameAvailabilityResult{
			NameAvailable: to.BoolPtr(true),
		},
		nil,
	)
	accountsCreateFeature := storageAccount.AccountsCreateFuture{}
	accountsClient.EXPECT().Create(
		gomock.Any(),
		resourceGroup,
		storageAccountName,
		storageAccount.AccountCreateParameters{
			Sku: &storageAccount.Sku{
				Name: storageAccount.StandardLRS,
			},
			Location: to.StringPtr(location),
			AccountPropertiesCreateParameters: &storageAccount.AccountPropertiesCreateParameters{},
		},
	).Return(
		accountsCreateFeature,
		nil,
	)

	accountKey := "SomeRandomStorageAccountKey"

	var keys = []storageAccount.AccountKey{
		storageAccount.AccountKey{
			KeyName: to.StringPtr("Primary"),
			Value:   to.StringPtr(accountKey),
		},
	}

	accountsClient.EXPECT().ListKeys(
		gomock.Any(),
		resourceGroup,
		storageAccountName,
	).Return(
		storageAccount.AccountListKeysResult{
			Keys: &keys,
		},
		nil,
	)

	waitForCompletion.EXPECT().Wait(
		accountsCreateFeature,
		gomock.Any(),
		gomock.Any(),
	).Return(
		nil,
	)
	accountKeys, err := storageAccountClient.CreateStorageAccountIfNotExists(storageAccountName, resourceGroup, location)
	assert.Equal(t, accountKey, *(*accountKeys)[0].Value, "AccountKey should be equal")
	assert.Nil(t, err)
}

func TestCreateStorageAccountDeleteIfExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accountsClient := mock.NewMockIAccountsClient(ctrl)
	storageAccountClient := &StorageAccountClient{
		Client:            accountsClient,
		AutoRestClient:    nil,
		WaitForCompletion: nil,
	}
	resourceGroup := "someResource"
	storageAccountNameExists := "someStorage01"
	storageAccountNameNotExists := "someStorage02"
	accountsClient.EXPECT().GetProperties(gomock.Any(), resourceGroup, storageAccountNameExists).Return(
		storageAccount.Account{},
		nil,
	).Times(1)
	accountsClient.EXPECT().GetProperties(gomock.Any(), resourceGroup, storageAccountNameNotExists).Return(
		storageAccount.Account{},
		errors.New("Storage Account can not find"),
	).Times(1)
	accountsClient.EXPECT().Delete(gomock.Any(), resourceGroup, storageAccountNameExists).Times(1)
	storageAccountClient.DeleteIfExists(storageAccountNameExists, resourceGroup)
	storageAccountClient.DeleteIfExists(storageAccountNameNotExists, resourceGroup)

}
