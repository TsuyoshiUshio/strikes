package providers

import (
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

type ServicePrincipalFixture struct {
	ExpectedClientID          string
	ExpectedClientSecret      string
	ExpectedSubscriptionID    string
	ExpectedTenantID          string
	ExpectedOriginalParameter string
}

func (s *ServicePrincipalFixture) Setup() {
	fakeNewConfigContext := func() (*config.ConfigContext, error) {
		return &config.ConfigContext{}, nil
	}

	fakeGetConfig := func(context *config.ConfigContext) (*config.Config, error) {
		return &config.Config{
			ClientID:       s.ExpectedClientID,
			ClientSecret:   s.ExpectedClientSecret,
			SubscriptionID: s.ExpectedSubscriptionID,
			TenantID:       s.ExpectedTenantID,
		}, nil
	}

	s.Patch(fakeNewConfigContext, fakeGetConfig)
}

func (s *ServicePrincipalFixture) Patch(
	fakeNewConfigContext func() (*config.ConfigContext, error),
	fakeGetConfig func(context *config.ConfigContext) (*config.Config, error),
) {
	monkey.Patch(config.NewConfigContext, fakeNewConfigContext)
	var conf *config.ConfigContext
	monkey.PatchInstanceMethod(reflect.TypeOf(conf), "GetConfig", fakeGetConfig)
}

func TestAddServicePrincipalParameters(t *testing.T) {
	fixture := &ServicePrincipalFixture{
		ExpectedClientID:          "foo",
		ExpectedClientSecret:      "bar",
		ExpectedSubscriptionID:    "baz",
		ExpectedTenantID:          "qux",
		ExpectedOriginalParameter: "quux",
	}
	fixture.Setup()

	defer monkey.UnpatchAll()
	m := make(map[string]string)
	m["foo"] = fixture.ExpectedOriginalParameter
	result := addServicePrincipalParameters(&m)

	assert.Equal(t, fixture.ExpectedOriginalParameter, (*result)["foo"], "Original Parameter is wrong.")
	assert.Equal(t, fixture.ExpectedClientID, (*result)["client_id"], "ClientID is wrong.")
	assert.Equal(t, fixture.ExpectedClientSecret, (*result)["client_secret"], "ClientSecret is wrong.")
	assert.Equal(t, fixture.ExpectedSubscriptionID, (*result)["subscription_id"], "SubscriptionID is wrong.")
	assert.Equal(t, fixture.ExpectedTenantID, (*result)["tenant_id"], "TenantID is wrong.")
}

func TestAddServicePricipalParametersWithoutContext(t *testing.T) {

}
