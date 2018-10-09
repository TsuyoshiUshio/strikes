package command

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestOutputNote(t *testing.T) {
	InputDocument := `
	hello-world package successfully deployed!

	Please refer these resources: 
	
	Resource Group  : {{.ResourceGroupTemplate}}
	Azure Functions : {{.AzureFunctionsTemplate}}
	
	Go to https://{{.EnvironmentBaseNameTemplate}}app.azurewebsites.net/api/HelloWorld
	
	Happy strikes!	
	`
	ExpectedDocument := `
	hello-world package successfully deployed!

	Please refer these resources: 
	
	Resource Group  : foo
	Azure Functions : barapp.azurewebsites.net
	
	Go to https://barapp.azurewebsites.net/api/HelloWorld
	
	Happy strikes!	
	`

	ExpectedResourceGroup := "foo"
	ExpectedInstanceName := "bar"
	InputTargetPath := "baz"
	ExpectedFilePath := "baz/NOTE.txt"

	var ActualFilePath string

	fakeOpen := func(name string) (*os.File, error) {
		ActualFilePath = name
		return &os.File{}, nil
	}
	fakeReadAll := func(r io.Reader) ([]byte, error) {
		return []byte(InputDocument), nil
	}

	monkey.Patch(os.Open, fakeOpen)
	monkey.Patch(ioutil.ReadAll, fakeReadAll)
	defer monkey.UnpatchAll()
	buffer := &bytes.Buffer{}
	outputNote(InputTargetPath, ExpectedResourceGroup, ExpectedInstanceName, buffer)

	actual := buffer.String()
	assert.Equal(t, ExpectedDocument, string(actual))
	assert.Equal(t, ExpectedFilePath, ActualFilePath)
}
