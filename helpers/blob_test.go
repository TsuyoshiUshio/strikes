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
