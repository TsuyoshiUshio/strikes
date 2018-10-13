package providers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/bouk/monkey"
	"github.com/stretchr/stew/slice"
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
		"foo=bar",
		"bar=foo",
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
		"foobar",
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
		"environment_base_name=" + ExpectedEnvironmentBaseName,
		"packages_sub_dir=" + ExpectedPackageSubDir,
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
	fixture := &ServicePrincipalFixture{
		ExpectedClientID:          "foo",
		ExpectedClientSecret:      "bar",
		ExpectedSubscriptionID:    "baz",
		ExpectedTenantID:          "qux",
		ExpectedOriginalParameter: "quux",
	}
	fixture.Setup()

	defer monkey.UnpatchAll()
	// -var 'foo=bar' is the terraform parameters.
	m := make(map[string]string)
	m["foo"] = "bar"
	m["hoge"] = "fuga"
	results := getTerraformParameter(&m)
	for i := 0; i < len(*results); i = i + 2 {
		assert.Equal(t, "-var", (*results)[i])
	}
	f := func(s string) bool {
		if strings.Contains(s, "=") {
			return true
		}
		return false
	}
	values := Filter(*results, f)

	assert.True(t, slice.Contains(values, fmt.Sprintf("client_id=%s", fixture.ExpectedClientID)))
	assert.True(t, slice.Contains(values, fmt.Sprintf("client_secret=%s", fixture.ExpectedClientSecret)))
	assert.True(t, slice.Contains(values, fmt.Sprintf("subscription_id=%s", fixture.ExpectedSubscriptionID)))
	assert.True(t, slice.Contains(values, fmt.Sprintf("tenant_id=%s", fixture.ExpectedTenantID)))
	assert.True(t, slice.Contains(values, "foo=bar"))
	assert.True(t, slice.Contains(values, "hoge=fuga"))
}

func Filter(params []string, f func(s string) bool) []string {
	result := make([]string, 0)
	for _, p := range params {
		if f(p) {
			result = append(result, p)
		}
	}
	return result
}

type ServicePrincipalFixture struct {
	ExpectedClientID          string
	ExpectedClientSecret      string
	ExpectedSubscriptionID    string
	ExpectedTenantID          string
	ExpectedOriginalParameter string
	ExpectedError             string
	ActualOutputBuffer        *bytes.Buffer
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

func (s *ServicePrincipalFixture) SetupWihtoutContext() {
	fakeNewConfigContext := func() (*config.ConfigContext, error) {
		return &config.ConfigContext{}, errors.New(s.ExpectedError)
	}

	fakeGetConfig := func(context *config.ConfigContext) (*config.Config, error) {
		return &config.Config{
			ClientID:       s.ExpectedClientID,
			ClientSecret:   s.ExpectedClientSecret,
			SubscriptionID: s.ExpectedSubscriptionID,
			TenantID:       s.ExpectedTenantID,
		}, nil
	}
	s.PatchWithOutput(fakeNewConfigContext, fakeGetConfig)
}

func (s *ServicePrincipalFixture) SetupWihtoutConfig() {
	fakeNewConfigContext := func() (*config.ConfigContext, error) {
		return &config.ConfigContext{}, nil
	}

	fakeGetConfig := func(context *config.ConfigContext) (*config.Config, error) {
		return nil, errors.New(s.ExpectedError)
	}
	s.PatchWithOutput(fakeNewConfigContext, fakeGetConfig)
}

func (s *ServicePrincipalFixture) PatchWithOutput(
	fakeNewConfigContext func() (*config.ConfigContext, error),
	fakeGetConfig func(context *config.ConfigContext) (*config.Config, error),
) {
	fakeExit := func(code int) {
		// ignore.
	}
	s.ActualOutputBuffer = new(bytes.Buffer)

	log.SetOutput(s.ActualOutputBuffer)
	monkey.Patch(os.Exit, fakeExit)

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
func (s *ServicePrincipalFixture) Output() string {
	output, err := ioutil.ReadAll(s.ActualOutputBuffer)
	if err != nil {
		log.Printf("[DEBUG] %v, \n", err)
		return "" // panic doesn't work for the moneky patching.
	}
	return string(output)
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
	fixture := &ServicePrincipalFixture{
		ExpectedClientID:          "foo",
		ExpectedClientSecret:      "bar",
		ExpectedSubscriptionID:    "baz",
		ExpectedTenantID:          "qux",
		ExpectedOriginalParameter: "quux",
		ExpectedError:             "corge",
	}
	fixture.SetupWihtoutContext()

	defer monkey.UnpatchAll()
	m := make(map[string]string)
	m["foo"] = fixture.ExpectedOriginalParameter
	_ = addServicePrincipalParameters(&m)
	assert.Regexp(t, fixture.ExpectedError, string(fixture.Output()))
}

func TestAddServicePricipalParametersWithoutConfig(t *testing.T) {
	fixture := &ServicePrincipalFixture{
		ExpectedClientID:          "foo",
		ExpectedClientSecret:      "bar",
		ExpectedSubscriptionID:    "baz",
		ExpectedTenantID:          "qux",
		ExpectedOriginalParameter: "quux",
		ExpectedError:             "corge",
	}
	fixture.SetupWihtoutConfig()

	defer monkey.UnpatchAll()
	m := make(map[string]string)
	m["foo"] = fixture.ExpectedOriginalParameter
	_ = addServicePrincipalParameters(&m)
	assert.Regexp(t, fixture.ExpectedError, string(fixture.Output()))
}

func TestAddManifestParameters(t *testing.T) {
	ExpectedPackageName := "foo"
	ExpectedPackageVersion := "1.0.0"
	ExpectedPackageZipName01 := "bar.zip"
	ExpectedPackageZipName02 := "baz.zip"
	manifest := config.Manifest{
		Name:    ExpectedPackageName,
		Version: ExpectedPackageVersion,
		ZipFileNames: []string{
			ExpectedPackageZipName01,
			ExpectedPackageZipName02,
		},
	}
	provider := TerraformProvider{
		Manifest: &manifest,
	}

	currentParams := make(map[string]string)
	currentParams["foo"] = "bar"
	currentParams["bar"] = "baz"
	actual := provider.addManifestParameters(&currentParams)
	assert.Equal(t, "bar", (*actual)["foo"])
	assert.Equal(t, "baz", (*actual)["bar"])
	assert.Equal(t, ExpectedPackageName, (*actual)[PackageName])
	assert.Equal(t, ExpectedPackageVersion, (*actual)[PackageVersion])
	assert.Equal(t, ExpectedPackageZipName01, (*actual)[fmt.Sprintf("%s_0", PackageZipNameBase)])
	assert.Equal(t, ExpectedPackageZipName02, (*actual)[fmt.Sprintf("%s_1", PackageZipNameBase)])
}

func TestOverrideEnvironmentBaseNameAndResourceGroupIfSpecifiedNormalCase(t *testing.T) {
	ExpectedExistsParameterFoo := "bar"
	ExpectedExistsParameterBar := "baz"
	ExpectedEnvironmentBaseName := "qux"
	ExpectedResourceGroup := ExpectedEnvironmentBaseName + "-rg"
	m := make(map[string]string)
	m["foo"] = ExpectedExistsParameterFoo
	m["bar"] = ExpectedExistsParameterBar
	instanceName := ExpectedEnvironmentBaseName

	Actual := overrideEnvironmentBaseNameAndResourceGroupIfSpecified(&m, instanceName)
	assert.Equal(t, ExpectedExistsParameterFoo, (*Actual)["foo"])
	assert.Equal(t, ExpectedExistsParameterBar, (*Actual)["bar"])
	assert.Equal(t, ExpectedEnvironmentBaseName, (*Actual)["environment_base_name"])
	assert.Equal(t, ExpectedResourceGroup, (*Actual)["resource_group"])
}

func TestOverrideEnvironmentBaseNameAndResourceGroupIfSpecifiedNotSpecifiedCase(t *testing.T) {
	ExpectedExistsParameterFoo := "bar"
	ExpectedExistsParameterBar := "baz"
	m := make(map[string]string)
	m["foo"] = ExpectedExistsParameterFoo
	m["bar"] = ExpectedExistsParameterBar

	Actual := overrideEnvironmentBaseNameAndResourceGroupIfSpecified(&m, "")
	assert.Equal(t, ExpectedExistsParameterFoo, (*Actual)["foo"])
	assert.Equal(t, ExpectedExistsParameterBar, (*Actual)["bar"])
	_, IsActualEnvironmentBaseName := (*Actual)["environment_base_name"]
	_, IsActualResourceGroup := (*Actual)["resource_group"]
	assert.False(t, IsActualEnvironmentBaseName, "environment_base_name shold not be specified.")
	assert.False(t, IsActualResourceGroup, "resource_group shold not be specified.")
}

func TestOverrideEnvironmentBaseNameAndResourceGroupIfSpecifiedNormalWithResourceGroup(t *testing.T) {
	ExpectedExistsParameterFoo := "bar"
	ExpectedExistsParameterBar := "baz"
	ExpectedEnvironmentBaseName := "qux"
	ExpectedResourceGroup := ExpectedEnvironmentBaseName + "-rg"
	m := make(map[string]string)
	m["environment_base_name"] = "quuz"
	m["resource_group"] = "quuz-rg"
	m["foo"] = ExpectedExistsParameterFoo
	m["bar"] = ExpectedExistsParameterBar
	instanceName := ExpectedEnvironmentBaseName

	Actual := overrideEnvironmentBaseNameAndResourceGroupIfSpecified(&m, instanceName)
	assert.Equal(t, ExpectedExistsParameterFoo, (*Actual)["foo"], "Expected")
	assert.Equal(t, ExpectedExistsParameterBar, (*Actual)["bar"])
	assert.Equal(t, ExpectedEnvironmentBaseName, (*Actual)["environment_base_name"])
	assert.Equal(t, ExpectedResourceGroup, (*Actual)["resource_group"])
}
