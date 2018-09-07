package repository

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositoryToken(t *testing.T) {
	expectedStorageAccountName := "foo"
	expectedContainerName := "bar"
	expectedSASQueryParameter := "foobar"
	expectedURL := "https://repository.simplearchitect.club/api/assetserveruri"

	responseDoc := heredoc.Docf(`{
	"StorageAccountName":"%s",
	"ContainerName":"%s",
	"SASQueryParameter":"%s"
	}
	`, expectedStorageAccountName, expectedContainerName, expectedSASQueryParameter)
	var actualURL string
	fakeGet := func(url string) (resp *http.Response, err error) {
		actualURL = url

		return &http.Response{
			Body: nopCloser{strings.NewReader(responseDoc)},
		}, nil
	}
	patch := monkey.Patch(http.Get, fakeGet)
	defer patch.Unpatch()
	// assert.Equal(t, "aa", responseDoc)
	token, err := GetRepositoryAccessToken()
	assert.Nil(t, err, "AccessToken error object should be nil")
	assert.Equal(t, expectedURL, actualURL, "Expected URL is not passed to the http.Get")
	assert.Equal(t, expectedStorageAccountName, token.StorageAccountName)
	assert.Equal(t, expectedContainerName, token.ContainerName)
	assert.Equal(t, expectedSASQueryParameter, token.SASQueryParameter)
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
