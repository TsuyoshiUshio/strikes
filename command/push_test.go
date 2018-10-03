package command

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/services/repository"
	"github.com/stretchr/testify/assert"

	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/bouk/monkey"
)

// Reference FileInfo interface
// type FileInfo interface {
//     Name() string
//     Size() int64
//     Mode() FileMode
//     ModTime() time.Time
//     IsDir() bool
//     Sys() interface{}
// }

type fileInfoImpl struct {
	TargetName  string
	TargetIsDir bool
}

func (f *fileInfoImpl) Name() string {
	return f.TargetName
}

func (f *fileInfoImpl) Size() int64 {
	return 0
}

func (f *fileInfoImpl) Mode() os.FileMode {
	return os.ModePerm
}

func (f *fileInfoImpl) ModTime() time.Time {
	return time.Now()
}

func (f *fileInfoImpl) IsDir() bool {
	return f.TargetIsDir
}

func (f *fileInfoImpl) Sys() interface{} {
	return nil
}

func TestCreatePackageBlob(t *testing.T) {
	ExpectedZipFileName01 := "foo.zip"
	ExpectedZipFileName02 := "qux.zip"
	// Setup
	fakeRunDir := func(dirname string) ([]os.FileInfo, error) {
		files := []os.FileInfo{
			&fileInfoImpl{
				TargetName:  ExpectedZipFileName01,
				TargetIsDir: false,
			},
			&fileInfoImpl{
				TargetName:  "bar",
				TargetIsDir: true,
			},
			&fileInfoImpl{
				TargetName:  ".baz.tmp",
				TargetIsDir: false,
			},
			&fileInfoImpl{
				TargetName:  ExpectedZipFileName02,
				TargetIsDir: false,
			},
		}
		return files, nil
	}

	type BlobParameter struct {
		StorageAccountName string
		ContainerName      string
		BlobName           string
		SasQueryParameter  string
		URL                string
	}
	ActualBlobs := make([]BlobParameter, 0)

	ActualFilePaths := make([]string, 0)
	fakeUpload := func(b *helpers.BlockBlob, uploadFilePath string) error {
		ActualFilePaths = append(ActualFilePaths, uploadFilePath)
		ActualBlobs = append(ActualBlobs, BlobParameter{
			StorageAccountName: b.StorageAccountName,
			ContainerName:      b.ContainerName,
			BlobName:           b.BlobName,
			SasQueryParameter:  b.SASQueryParameter,
			URL:                b.BlockBlobURL.URL().Host + b.BlockBlobURL.URL().Path,
		})
		return nil
	}

	monkey.Patch(ioutil.ReadDir, fakeRunDir)
	var blockBlob *helpers.BlockBlob
	monkey.PatchInstanceMethod(reflect.TypeOf(blockBlob), "Upload", fakeUpload)
	defer monkey.UnpatchAll()
	ExpectedStorageAccount := "foo"
	ExpectedContainerName := "bar"
	ExpectedStorageAccountBaseURL := ExpectedStorageAccount + ".blob.core.windows.net/" + ExpectedContainerName
	ExpectedSASQueryParameter := "?code=baz"
	token := repository.RepositoryAccessToken{
		StorageAccountName: ExpectedStorageAccount,
		ContainerName:      ExpectedContainerName,
		SASQueryParameter:  ExpectedSASQueryParameter,
	}
	ExpectedPackageName := "qux"
	ExpectedPackageVersion := "1.0.0"
	manifest := config.Manifest{
		Name:    ExpectedPackageName,
		Version: ExpectedPackageVersion,
	}
	ExpectedPackageBase := "sample/hello-world"

	ExpectedPackagePath01 := ExpectedPackageBase + "/package/" + ExpectedZipFileName01
	ExpectedBlobName01 := ExpectedPackageName + "/" + ExpectedPackageVersion + "/package/" + ExpectedZipFileName01
	ExpectedURL01 := ExpectedStorageAccountBaseURL + "/" + ExpectedBlobName01
	ExpectedPackagePath02 := ExpectedPackageBase + "/package/" + ExpectedZipFileName02
	ExpectedBlobName02 := ExpectedPackageName + "/" + ExpectedPackageVersion + "/package/" + ExpectedZipFileName02
	ExpectedURL02 := ExpectedStorageAccountBaseURL + "/" + ExpectedBlobName02

	createPackageBlockBlob(&token, &manifest, ExpectedPackageBase)
	assert.Equal(t, ExpectedPackagePath01, ActualFilePaths[0])
	assert.Equal(t, ExpectedPackagePath02, ActualFilePaths[1])
	assert.Equal(t, ExpectedStorageAccount, ActualBlobs[0].StorageAccountName)
	assert.Equal(t, ExpectedContainerName, ActualBlobs[0].ContainerName)
	assert.Equal(t, ExpectedBlobName01, ActualBlobs[0].BlobName)
	assert.Equal(t, ExpectedBlobName02, ActualBlobs[1].BlobName)
	assert.Equal(t, ExpectedSASQueryParameter, ActualBlobs[0].SasQueryParameter)
	assert.Equal(t, ExpectedURL01, ActualBlobs[0].URL)
	assert.Equal(t, ExpectedURL02, ActualBlobs[1].URL)
}
