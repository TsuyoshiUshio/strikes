package helpers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestDownloadNormalScenario(t *testing.T) {
	ExpectedContent := "foo"
	ExpectedFileName := "bar"
	ExpectedUrl := "https://foo.bar.com"

	fixture := newFixture()
	fixture.SetUp(ExpectedContent)

	err := fixture.Execute(DownloadFile, ExpectedFileName, ExpectedUrl)
	assert.Nil(t, err)
	fixture.GetContent()

	assert.Equal(t, ExpectedFileName, fixture.ActualFileName, "Input filename is wrong")
	assert.Equal(t, ExpectedUrl, fixture.ActualUrl, "Input Url is wrong")
	assert.Equal(t, ExpectedContent, fixture.ActualContent, "Content mismatch")

	fixture.CleanUp()
}

type parameterFixture struct {
	ActualFileName  string
	ActualUrl       string
	ActualContent   string
	ExpectedContent string
	FakeFile        *os.File
	FakeCreate      interface{}
	FakeGet         interface{}
}

func newFixture() *parameterFixture {
	return &parameterFixture{}
}

func (f *parameterFixture) SetUp(expectedContent string) {
	f.ExpectedContent = expectedContent

	file, err := ioutil.TempFile(".", ".temp")
	if err != nil {
		panic(err)
	}
	f.FakeFile = file

	f.FakeCreate = func(name string) (*os.File, error) {
		f.ActualFileName = name
		return file, nil
	}
	f.FakeGet = func(url string) (*http.Response, error) {
		f.ActualUrl = url
		return &http.Response{
			Body: &NopCloser{strings.NewReader(f.ExpectedContent)},
		}, nil
	}

}

func (f *parameterFixture) Execute(downloadFile func(string, string) error, expectedFileName string, expectedUrl string) error {
	httpPatch := monkey.Patch(http.Get, f.FakeGet)
	defer httpPatch.Unpatch()
	osPatch := monkey.Patch(os.Create, f.FakeCreate)
	defer osPatch.Unpatch()

	return downloadFile(expectedFileName, expectedUrl)
}

func (f *parameterFixture) GetContent() {
	dat, err := ioutil.ReadFile(f.FakeFile.Name())
	if err != nil {
		panic(err)
	}
	f.ActualContent = string(dat)
}

func (f *parameterFixture) CleanUp() {
	err := os.Remove(f.FakeFile.Name())
	if err != nil {
		panic(err)
	}
}
