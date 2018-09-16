package repository

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositoryToken(t *testing.T) {
	expectedStorageAccountName := "foo"
	expectedContainerName := "bar"
	expectedSASQueryParameter := "foobar"
	expectedURL := RepositoryBaseURL + "assetserveruri"

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
			Body: helpers.NopCloser{strings.NewReader(responseDoc)},
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

func TestCreateNormalCase(t *testing.T) {

	fixture := NewFixturePackage()

	var actualURL string
	var actualContentType string
	var actualBody []byte

	fakePost := func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
		actualURL = url
		actualContentType = contentType
		actualBody, _ = ioutil.ReadAll(body)
		return &http.Response{
			StatusCode: http.StatusCreated,
		}, nil
	}
	fakeExit := func(int) {
		// do nothing
	}

	patch := monkey.Patch(http.Post, fakePost)
	patchFatalf := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	defer patchFatalf.Unpatch()

	resp, err := fixture.InputPackage.Create()

	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, fixture.ExpectedURL, actualURL, "URL is wrong.")
	assert.Equal(t, fixture.ExpectedContentType, actualContentType, "ContentType is wrong.")
	assert.Equal(t, fixture.ExpectedBody, actualBody, "Body format is wrong.")
	assert.Equal(t, fixture.ExpectedStatusCode, resp.StatusCode, "StatusCode is wrong.")

}

func TestCreateErrorCase(t *testing.T) {
	fixture := NewFixturePackage()
	expectedError := errors.New("Internal Server Error.")
	fakePost := func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
		return nil, expectedError
	}
	fakeExit := func(int) {
		// do nothing
	}
	patch := monkey.Patch(http.Post, fakePost)
	patchFatalf := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	defer patchFatalf.Unpatch()

	_, err := fixture.InputPackage.Create()
	assert.Equal(t, expectedError, err, "Different error has been thrown.")
}

func TestCreateBadRequest(t *testing.T) {
	expectedStatusCode := http.StatusBadRequest
	expectedErrorMessage := "HTTP Status code should be 201 (Created). Current Status Code is 400"
	fixture := NewFixturePackage()

	fakePost := func(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: expectedStatusCode,
		}, nil
	}
	fakeExit := func(int) {
		// do nothing
	}
	patch := monkey.Patch(http.Post, fakePost)
	patchFatalf := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	defer patchFatalf.Unpatch()
	resp, err := fixture.InputPackage.Create()

	assert.Equal(t, err.Error(), expectedErrorMessage, "Error should include expected Error Message.")
	assert.Equal(t, expectedStatusCode, resp.StatusCode, "Different error has been thrown.")
}

type fixturePackage struct {
	InputPackage        *Package
	ExpectedURL         string
	ExpectedContentType string
	ExpectedBody        []byte
	ExpectedStatusCode  int
}

func NewFixturePackage() *fixturePackage {

	inputPackage := NewPackageWithCurrentTime(
		"fooName",
		"fooDescription",
		"foo",
		"http://foo.bar",
		"http://foo.com",
		"1.0.0",
		"foobar",
		"Terraform",
	)
	expectedBody, _ := inputPackage.Marshal()
	return &fixturePackage{
		InputPackage:        inputPackage,
		ExpectedURL:         RepositoryBaseURL + "package",
		ExpectedContentType: helpers.ContentTypeApplicationJson,
		ExpectedBody:        expectedBody,
		ExpectedStatusCode:  http.StatusCreated,
	}
}

func TestGetPackageByName(t *testing.T) {
	ExpectedPackageName := "foo"
	ExpectedPackage := NewPackageWithCurrentTime(
		ExpectedPackageName,
		"desc foo",
		"bar",
		"https://foo.bar.com",
		"https://www.foo.bar.com",
		"1.0.0",
		"release foo",
		"Terraform",
	)
	ExpectedURL := RepositoryBaseURL + "package?name=" + ExpectedPackageName
	var ActualURL string
	fakeGet := func(url string) (resp *http.Response, err error) {
		ActualURL = url
		jsonPackage, _ := json.Marshal(ExpectedPackage)
		result := &http.Response{
			StatusCode: http.StatusOK,
			Body:       helpers.NopCloser{strings.NewReader(string(jsonPackage))},
		}
		return result, nil
	}
	patch := monkey.Patch(http.Get, fakeGet)
	defer patch.Unpatch()
	p, err := GetPackage(ExpectedPackageName)
	assert.Nil(t, err)
	assert.Equal(t, ExpectedPackageName, p.Name)
	assert.Equal(t, ExpectedURL, ActualURL)
}
