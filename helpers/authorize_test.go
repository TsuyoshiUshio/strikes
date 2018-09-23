package helpers

import (
	"testing"

	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthorizer(t *testing.T) {
	ExpectedAzureAuthLocation := "foo"
	var ExpectedAuthorizerHeaderAddr *autorest.APIKeyAuthorizer
	fakeNewAuthorizer := func(baseURI string) (autorest.Authorizer, error) {
		header := make(map[string]interface{})
		authorizer := autorest.NewAPIKeyAuthorizerWithHeaders(header)
		ExpectedAuthorizerHeaderAddr = authorizer
		return authorizer, nil
	}
	monkey.Patch(auth.NewAuthorizerFromFile, fakeNewAuthorizer)
	defer monkey.Unpatch(auth.NewAuthorizerFromFile)
	authorizer, _ := NewAuthorizer(ExpectedAzureAuthLocation)

	// AZURE_AUTH_LOCATION environment variables has the
	assert.Equal(t, ExpectedAzureAuthLocation, os.Getenv(AZURE_AUTH_LOCATION))
	assert.Equal(t, ExpectedAuthorizerHeaderAddr, authorizer)
}
