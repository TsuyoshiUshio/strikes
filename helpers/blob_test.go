package helpers

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/bouk/monkey"
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

func TestUpload(t *testing.T) {
	tempFile, err := ioutil.TempFile(".", "uploadTest")
	if err != nil {
		panic(err)
	}
	tempFile.WriteString("foo")
	tempFile.Close()
	defer os.Remove(tempFile.Name())
	var ActualBodyReader io.ReadSeeker
	fakeUpload := func(_ azblob.BlockBlobURL, _ context.Context, body io.ReadSeeker, _ azblob.BlobHTTPHeaders, _ azblob.Metadata, _ azblob.BlobAccessConditions) (*azblob.BlockBlobUploadResponse, error) {
		ActualBodyReader = body
		return nil, nil
	}
	var b azblob.BlockBlobURL
	monkey.PatchInstanceMethod(reflect.TypeOf(b), "Upload", fakeUpload)
	defer monkey.UnpatchAll()
	blockBlob := NewBlockBlob("foo", "bar", "baz")
	blockBlob.Upload(tempFile.Name())
	ActualBody, _ := ioutil.ReadAll(ActualBodyReader)

	assert.Equal(t, "foo", string(ActualBody))
}
