package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockBlob(t *testing.T) {
	ExpectedBlockBlobURL := "https://foo.blob.core.windows.net/bar/baz"
	ExpectedStorageAccountName := "foo"
	ExpectedContainerName := "bar"
	ExpectedBlobName := "baz"
	blockBlob := NewBlockBlob(ExpectedStorageAccountName, ExpectedContainerName, ExpectedBlobName)
	assert.Equal(t, ExpectedBlockBlobURL, blockBlob.BlockBlobURL.String())
	assert.Equal(t, ExpectedStorageAccountName, blockBlob.StorageAccountName)
	assert.Equal(t, ExpectedContainerName, blockBlob.ContainerName)
	assert.Equal(t, ExpectedBlobName, blockBlob.BlobName)
}

func TestNewBlockBlobWithSASQueryParameter(t *testing.T) {
	ExpectedBlockBlobURL := "https://foo.blob.core.windows.net/bar/baz?code=qux"
	ExpectedStorageAccountName := "foo"
	ExpectedContainerName := "bar"
	ExpectedBlobName := "baz"
	ExpectedSASQueryParameter := "code=qux"
	blockBlob := NewBlockBlobWithSASQueryParameter(ExpectedStorageAccountName, ExpectedContainerName, ExpectedBlobName, ExpectedSASQueryParameter)
	assert.Equal(t, ExpectedBlockBlobURL, blockBlob.BlockBlobURL.String())
	assert.Equal(t, ExpectedStorageAccountName, blockBlob.StorageAccountName)
	assert.Equal(t, ExpectedContainerName, blockBlob.ContainerName)
	assert.Equal(t, ExpectedBlobName, blockBlob.BlobName)
	assert.Equal(t, ExpectedSASQueryParameter, blockBlob.SASQueryParameter)
}
