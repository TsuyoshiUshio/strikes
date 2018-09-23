package helpers

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
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
	ExpectedBody := "foo"
	ExpectedBlockBlobURL := "https://foo.blob.core.windows.net/bar/baz"
	tempFile, err := ioutil.TempFile(".", "uploadTest")
	if err != nil {
		panic(err)
	}
	tempFile.WriteString(ExpectedBody)
	tempFile.Close()
	defer os.Remove(tempFile.Name())
	var ActualBodyReader io.ReadSeeker
	var ActualBlockBlobURL azblob.BlockBlobURL
	fakeUpload := func(blockBlobURL azblob.BlockBlobURL, _ context.Context, body io.ReadSeeker, _ azblob.BlobHTTPHeaders, _ azblob.Metadata, _ azblob.BlobAccessConditions) (*azblob.BlockBlobUploadResponse, error) {
		ActualBodyReader = body
		ActualBlockBlobURL = blockBlobURL
		return nil, nil
	}
	var b azblob.BlockBlobURL
	monkey.PatchInstanceMethod(reflect.TypeOf(b), "Upload", fakeUpload)
	defer monkey.UnpatchAll()
	blockBlob := NewBlockBlob("foo", "bar", "baz")
	blockBlob.Upload(tempFile.Name())
	ActualBody, _ := ioutil.ReadAll(ActualBodyReader)

	assert.Equal(t, ExpectedBody, string(ActualBody), "expected body is wrong")
	assert.Equal(t, ExpectedBlockBlobURL, ActualBlockBlobURL.String(), "blob url is wrong.")
}

type StringReaderCloser struct {
	Reader io.Reader
}

func (s *StringReaderCloser) Close() error {
	return nil
}
func (s *StringReaderCloser) Read(b []byte) (n int, err error) {
	return s.Reader.Read(b)
}

func TestDownload(t *testing.T) {
	ExpectedContent := "foobar"
	blockBlob := NewBlockBlob("foo", "bar", "baz")
	fakeGet := func(url string) (resp *http.Response, err error) {
		readerCloser := &StringReaderCloser{
			Reader: strings.NewReader(ExpectedContent),
		}
		return &http.Response{
			Body: readerCloser,
		}, nil
	}
	monkey.Patch(http.Get, fakeGet)
	defer monkey.Unpatch(http.Get)
	tempFileName := "download.tmp"
	blockBlob.Download(tempFileName)
	defer os.Remove(tempFileName)
	content, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, ExpectedContent, string(content))

}
