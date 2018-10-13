package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/hashicorp/hcl"
)

type Provider interface {
	CreateResource(args []string, instanceName string, overrideParameters []string) *DeploymentResult
}

type TerraformProvider struct {
	Manifest  *config.Manifest
	TargetDir string
}

type DeploymentResult struct {
	Configrations *map[string]string
}

func (d DeploymentResult) GetResourceGroup() string {
	return (*d.Configrations)["resource_group"]
}

func (d DeploymentResult) GetConfigrationsJosn() string {
	content, err := json.Marshal(d.Configrations)
	if err != nil {
		log.Fatalf("Can not parse the configration to json: %v\n", err.Error())
	}
	return string(content)
}

func NewTerraformProvider(manifest *config.Manifest, targetDir string) *TerraformProvider {
	return &TerraformProvider{
		Manifest:  manifest,
		TargetDir: targetDir,
	}
}

func (t *TerraformProvider) IsProviderCommandExists() bool {
	// Check is there is terraform command is on the path.
	cmd := exec.Command("terraform", "--help")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func (t *TerraformProvider) CreateResource(args []string, instanceName string, overrideParameters []string) *DeploymentResult {
	if !t.IsProviderCommandExists() {
		log.Fatalf("Can not find the terraform command on your path. Please check if it is on your Path environment variables")
	}

	// Execute the terraform init

	fmt.Println("Initialiing terraform ...")
	t.executeTerraformCommand("init", []string{}, []string{}, false)

	// translate parameter fit for terraformf parameters
	argsParameters, configrations := t.composeTerraformParameter(overrideParameters, instanceName)
	// then append terraform options.
	optionalParameters := []string{
		"-input=false",
		"-auto-approve",
	}

	// Execute terraform plan

	//	fmt.Println("Executing terraform plan ...")

	//	t.executeTerraformCommand("plan", *argsParameters, optionalParameters, false)

	// Execute terraform apply

	fmt.Println("Executing terraform apply ...")

	t.executeTerraformCommand("apply", *argsParameters, optionalParameters, true)

	return &DeploymentResult{
		Configrations: configrations,
	}
}

func (t *TerraformProvider) executeTerraformCommand(subCommand string, argsParameters []string, optionalParameters []string, isDump bool) {
	terraformParameter := append([]string{subCommand}, argsParameters...)
	parameters := append(terraformParameter, optionalParameters...)
	for _, param := range parameters {
		log.Printf("[DEBUG] Parameters: %s\n", param)
	}
	cmd := exec.Command("terraform", parameters...)
	cmd.Dir = t.TargetDir
	if helpers.IsDebugEnabled() || isDump { // In case of Debug or specifed dump option, it will emit the results.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// output, err := cmd.CombinedOutput()
	err := cmd.Run()
	//	log.Printf("[DEBUG] %v", string(output))
	if err != nil {
		log.Fatalf("Can not execute the terraform %s: %v\n", subCommand, err)
	}
}

func (t *TerraformProvider) DeleteResource() error {
	// delete the resource group
	return nil
}

type Dictionaries struct {
	Name    string `hcl:",key"`
	Default string
}

type Values struct {
	Variable []Dictionaries
}

func parseValuesHcl(filePath string) (*map[string]string, error) {
	var values Values
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not read values.hcl file: %v\n", err)
		return nil, err
	}
	err = hcl.Decode(&values, string(b))
	m := make(map[string]string)
	for _, v := range values.Variable {
		m[v.Name] = v.Default
	}
	return &m, nil
}

func configureValues(values *map[string]string, args []string) (*map[string]string, error) {
	parameters, err := parseValuesArgs(args)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, key := range *getMapKeys(values) {
		if (*parameters)[key] != "" {
			m[key] = (*parameters)[key]
		} else {
			m[key] = (*values)[key]
		}
	}
	return &m, nil
}

func parseValuesArgs(args []string) (*map[string]string, error) {
	m := make(map[string]string)
	for _, arg := range args {
		keyValue := strings.Split(arg, "=")
		if len(keyValue) != 2 {
			log.Fatalf("Parameter can not parse. : %v\n", arg)
			return nil, errors.New("Parameter can not parse. : " + arg)
		}
		m[keyValue[0]] = keyValue[1]
	}
	return &m, nil
}

func getMapKeys(m *map[string]string) *[]string {
	keys := reflect.ValueOf(*m).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return &strkeys
}

func convertParameterMapToStringArray(values *map[string]string) *[]string {
	keys := getMapKeys(values)
	parameters := make([]string, len(*keys)*2, len(*keys)*2)
	index := 0
	for _, k := range *keys {
		parameters[index] = "-var"
		index++
		parameters[index] = fmt.Sprintf("%s=%s", k, (*values)[k])
		index++
	}
	return &parameters
}

func getTerraformParameter(values *map[string]string) *[]string {
	parameters := addServicePrincipalParameters(values)
	return convertParameterMapToStringArray(parameters)
}

func (t *TerraformProvider) composeTerraformParameter(args []string, instanceName string) (*[]string, *map[string]string) {
	defaultValues, err := parseValuesHcl(filepath.Join(t.TargetDir, "values.hcl"))
	if err != nil {
		log.Fatalf("Can not find values.hcl on the target path %v\n", t.TargetDir)
	}

	// if the instance name specified, environment_base_name is set by the instance name
	// and if the resource group name with {environment_base_name}-rg.
	overridenValues := overrideEnvironmentBaseNameAndResourceGroupIfSpecified(defaultValues, instanceName)

	// Adding manifest parameters
	manifestAddedValues := t.addManifestParameters(overridenValues)

	// If there are parameters which is the same as the Values.hcl, override it.
	parameterValues, err := configureValues(manifestAddedValues, args)
	if err != nil {
		log.Fatalf("Can not merge manifest values and parameter values. double check the arguments. ManifestValues: %v, Args: %v\n", defaultValues, args)
	}

	// translate parameter fit for terraformf parameters
	return getTerraformParameter(parameterValues), parameterValues

}

func addServicePrincipalParameters(parameters *map[string]string) *map[string]string {
	c, err := config.NewConfigContext()
	if err != nil {
		log.Fatalf("Can not create the ConfigContext: Double check the config path is correct. : %v \n", err)
		return nil
	}
	conf, err := c.GetConfig()
	if err != nil {
		log.Fatalf("Can not create config file: Please check the ~/.strikes/conf file. :%v\n", err)
		return nil
	}

	(*parameters)["client_id"] = conf.ClientID
	(*parameters)["client_secret"] = conf.ClientSecret
	(*parameters)["subscription_id"] = conf.SubscriptionID
	(*parameters)["tenant_id"] = conf.TenantID
	return parameters
}

const (
	PackageName        = "package_name"
	PackageVersion     = "package_version"
	PackageZipNameBase = "package_zip_name"
)

const (
	EnvironmentBaseNameParameter = "environment_base_name"
	ResourceGroupParameter       = "resource_group"
)

func (t *TerraformProvider) addManifestParameters(parameters *map[string]string) *map[string]string {
	(*parameters)["package_name"] = t.Manifest.Name
	(*parameters)["package_version"] = t.Manifest.Version
	for i, name := range t.Manifest.ZipFileNames {
		(*parameters)[fmt.Sprintf("%s_%d", PackageZipNameBase, i)] = name
	}
	return parameters
}

func overrideEnvironmentBaseNameAndResourceGroupIfSpecified(parameters *map[string]string, instanceName string) *map[string]string {
	if instanceName == "" {
		return parameters
	}

	// These parameter is mandatory.
	(*parameters)[EnvironmentBaseNameParameter] = instanceName
	(*parameters)[ResourceGroupParameter] = instanceName + "-rg"

	return parameters
}
