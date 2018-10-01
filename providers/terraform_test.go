package providers

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestParseValuesHCL(t *testing.T) {
	targetPath := filepath.Join(filepath.Join(".", "test-fixture"), "values.hcl")
	result, err := parseValuesHcl(targetPath)
	assert.Nil(t, err)
	parseValuesHcl(targetPath)
	assert.Equal(t, "hello-world", (*result)["environment_base_name"])
	assert.Equal(t, "hello-world-rg", (*result)["resource_group"])
	assert.Equal(t, "hello-world/1.0.0/hello.zip", (*result)["packages_sub_dir"])
}

func TestParseArgsNormalCase(t *testing.T) {
	args := []string{
		"--set",
		"foo=bar",
		"--set",
		"bar=foo",
		"foobar",
	}
	result, err := parseValuesArgs(args)
	assert.Nil(t, err)
	assert.Equal(t, "bar", (*result)["foo"])
	assert.Equal(t, "foo", (*result)["bar"])
	assert.Equal(t, 2, len(*result))
}

func TestParseArgsWorngParameterCase(t *testing.T) {
	ExpectedError := "Parameter can not parse. : foobar"
	args := []string{
		"--set",
		"foobar",
		"--set",
		"foo=bar",
	}

	// log.Fatalf will be called then this command exit.
	fakeExit := func(int) {
		// do nothing
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	_, err := parseValuesArgs(args)
	assert.Equal(t, ExpectedError, err.Error())
}

func TestConfigureValues(t *testing.T) {
	ExpectedEnvironmentBaseName := "foo"
	ExpectedResourceGroup := "bar"
	ExpectedPackageSubDir := "foo/2.0.0/bar.zip"
	args := []string{
		"--set",
		"environment_base_name=" + ExpectedEnvironmentBaseName,
		"--set",
		"packages_sub_dir=" + ExpectedPackageSubDir,
		"bar",
	}
	m := make(map[string]string)
	m["environment_base_name"] = "foobar"
	m["resource_group"] = ExpectedResourceGroup
	m["packages_sub_dir"] = "foo"

	result, err := configureValues(&m, args)
	assert.Nil(t, err)
	assert.Equal(t, ExpectedEnvironmentBaseName, (*result)["environment_base_name"])
	assert.Equal(t, ExpectedResourceGroup, (*result)["resource_group"])
	assert.Equal(t, ExpectedPackageSubDir, (*result)["packages_sub_dir"])
}

func TestGetTerraformParameter(t *testing.T) {
	// -var 'foo=bar' is the terraform parameters.
	m := make(map[string]string)
	m["foo"] = "bar"
	m["hoge"] = "fuga"
	results := getTerraformParameter(&m)
	assert.Equal(t, "-var 'foo=bar'", (*results)[0])
	assert.Equal(t, "-var 'hoge=fuga'", (*results)[1])
}

func TestAddServicePrincipalParameters(t *testing.T) {
	fakeNewConfigContext := func() (*config.ConfigContext, error) {
		return &config.ConfigContext{}, nil
	}
	ExpectedClientID := "foo"
	ExpectedClientSecret := "bar"
	ExpectedSubscriptionID := "baz"
	ExpectedTenantID := "qux"

	fakeGetConfig := func(context *config.ConfigContext) (*config.Config, error) {
		return &config.Config{
			ClientID:       ExpectedClientID,
			ClientSecret:   ExpectedClientSecret,
			SubscriptionID: ExpectedSubscriptionID,
			TenantID:       ExpectedTenantID,
		}, nil
	}
	monkey.Patch(config.NewConfigContext, fakeNewConfigContext)
	var conf *config.ConfigContext
	monkey.PatchInstanceMethod(reflect.TypeOf(conf), "GetConfig", fakeGetConfig)
	defer monkey.UnpatchAll()
	param := []string{
		"-var", "foo='bar'",
	}
	result, err := addServicePrincipalParameters(param)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "-var", result[2])
	assert.Equal(t, fmt.Sprintf("client_id='%s'", ExpectedClientID), result[3], "ClientID is wrong.")
	assert.Equal(t, "-var", result[4])
	assert.Equal(t, fmt.Sprintf("client_secret='%s'", ExpectedClientSecret), result[5], "ClientSecret is wrong.")
	assert.Equal(t, "-var", result[6])
	assert.Equal(t, fmt.Sprintf("subscription_id='%s'", ExpectedSubscriptionID), result[7], "SubscriptionID is wrong.")
	assert.Equal(t, "-var", result[8])
	assert.Equal(t, fmt.Sprintf("tenant_id='%s'", ExpectedTenantID), result[9], "TenantID is wrong.")

}
