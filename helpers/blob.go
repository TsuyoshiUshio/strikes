package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
)

type IBlockBlob interface {
	Upload(uploadFilePath string) error
	Download(filePath string) error
}

type BlockBlob struct {
	BlockBlobURL       *azblob.BlockBlobURL
	StorageAccountName string
	ContainerName      string
	BlobName           string
	SASQueryParameter  string
}

func NewBlockBlobWithSASQueryParameter(storageAccountName, containerName, blobName, sasQueryParameter string) *BlockBlob {
	blockBlobURL := createBlockBlobURLWithSasQueryParameter(storageAccountName, containerName, blobName, sasQueryParameter)
	return &BlockBlob{
		BlockBlobURL:       blockBlobURL,
		StorageAccountName: storageAccountName,
		ContainerName:      containerName,
		BlobName:           blobName,
		SASQueryParameter:  sasQueryParameter,
	}
}

func NewBlockBlob(storageAccountName, containerName, blobName string) *BlockBlob {
	blockBlobURL := createBlockBlobURL(storageAccountName, containerName, blobName)
	return &BlockBlob{
		BlockBlobURL:       blockBlobURL,
		StorageAccountName: storageAccountName,
		ContainerName:      containerName,
		BlobName:           blobName,
	}
}

func (b *BlockBlob) Upload(uploadFilePath string) error {

	f, err := ioutil.ReadFile(uploadFilePath)
	if err != nil {
		log.Fatalf("Can not read file to upload: %v \n", err)
		return err
	}
	err = b.upload(bytes.NewReader(f))

	log.Printf("%s uploaded... ", uploadFilePath)
	if err != nil {
		log.Fatalf("Can not upload file: %v \n", err)
		return err
	}
	return nil
}

func (b *BlockBlob) upload(readSeeker io.ReadSeeker) error {
	ctx := context.Background()
	_, err := b.BlockBlobURL.Upload(ctx, readSeeker, azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	return err
}

func createBlockBlobURLWithSasQueryParameter(storageAccountName, containerName, blobName, sasQueryParameter string) *azblob.BlockBlobURL {
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", storageAccountName, containerName, blobName, sasQueryParameter))
	blobURL := azblob.NewBlockBlobURL(*u,
		azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	return &blobURL
}

func createBlockBlobURL(storageAccountName, containerName, blobName string) *azblob.BlockBlobURL {
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageAccountName, containerName, blobName))
	blobURL := azblob.NewBlockBlobURL(*u,
		azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	return &blobURL
}

func (b *BlockBlob) Download(filePath string) error {
	err := downloadFile(filePath, fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", b.StorageAccountName, b.ContainerName, b.BlobName))
	if err != nil {
		log.Fatalf("Can not download file: %v", err)
		return err
	}
	return nil
}

func downloadFile(filepath, url string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}
