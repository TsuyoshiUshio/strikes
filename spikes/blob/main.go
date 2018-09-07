package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
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
	authInfo, err := readJson("./.config")
	// fetch the SAS token from stored Access Policy
	accountName := (*authInfo)["accountName"]
	accountKey := (*authInfo)["accountKey"]
	containerName := "repository"
	// upload zip file to the blob with subdirectory.
	qb := fetchSASQueryParameters(accountName, accountKey, containerName)
	fmt.Printf("SASKey for repository %s \n", qb)
	blobName := "hello-world/1.0.0/circuit/hello.zip"
	// upload zipfile to the blob with subdirectory
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, containerName, blobName, qb))
	blobURL := azblob.NewBlockBlobURL(*u,
		azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{}))
	ctx := context.Background()
	// reading file
	b, err := ioutil.ReadFile("./hello.zip")
	if err != nil {
		panic(err)
	}
	_, err = blobURL.Upload(ctx, bytes.NewReader(b), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	fmt.Println("Blob created.")

	// Download from blob download is open for everybody.
	downloadFile("downloaded_hello.zip", fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName))
	fmt.Println("Downloaded")

}

func fetchSASQueryParameters(accountName, accountKey, containerName string) string {

	fmt.Printf("AccountName: %s \n", accountName)
	fmt.Printf("AccountKey:  %s \n", accountKey)

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
