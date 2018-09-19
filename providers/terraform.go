package providers

import (
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
}

type TerraformProvider struct {
	Manifest  *config.Manifest
	TargetDir string
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

func (t *TerraformProvider) CreateResource(args []string) {
	if !t.IsProviderCommandExists() {
		log.Fatalf("Can not find the terraform command on your path. Please check if it is on your Path environment variables")
	}

	// Execute the terraform init

	fmt.Println("Initialiing terraform ...")
	t.executeTerraformCommand("init", []string{}, []string{}, false)

	// translate parameter fit for terraformf parameters
	argsParameters := t.composeTerraformParameter(args)
	// then append terraform options.
	optionalParameters := []string{}

	// Execute terraform plan

	fmt.Println("Executing terraform plan ...")

	t.executeTerraformCommand("plan", *argsParameters, optionalParameters, false)

	// Execute terraform apply

	fmt.Println("Executing terraform apply ...")

	t.executeTerraformCommand("apply", *argsParameters, optionalParameters, true)
}

func (t *TerraformProvider) executeTerraformCommand(subCommand string, argsParameters []string, optionalParameters []string, isDump bool) {
	terraformParameter := append([]string{subCommand}, argsParameters...)
	parameters := append(terraformParameter, optionalParameters...)

	cmd := exec.Command("terraform", parameters...)
	cmd.Path = t.TargetDir
	if helpers.IsDebugEnabled() || isDump { // In case of Debug or specifed dump option, it will emit the results.
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	_, err := cmd.CombinedOutput()
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
	flag := false
	m := make(map[string]string)
	for _, arg := range args {
		if flag {
			keyValue := strings.Split(arg, "=")
			if len(keyValue) != 2 {
				log.Fatalf("Parameter can not parse. : %v\n", arg)
				return nil, errors.New("Parameter can not parse. : " + arg)
			}
			m[keyValue[0]] = keyValue[1]
		}
		if arg == "--set" {
			flag = true
		} else {
			flag = false
		}

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

func getTerraformParameter(values *map[string]string) *[]string {
	keys := getMapKeys(values)
	parameters := make([]string, len(*keys), len(*keys))
	for i, k := range *keys {
		parameters[i] = fmt.Sprintf("-var '%s=%s'", k, (*values)[k])
	}
	return &parameters
}

func (t *TerraformProvider) composeTerraformParameter(args []string) *[]string {
	defaultValues, err := parseValuesHcl(filepath.Join(t.TargetDir, "values.hcl"))
	if err != nil {
		log.Fatalf("Can not find values.hcl on the target path %v\n", t.TargetDir)
	}
	// If there are parameters which is the same as the Values.hcl, override it.
	parameterValues, err := configureValues(defaultValues, args)
	if err != nil {
		log.Fatalf("Can not merge manifest values and parameter values. double check the arguments. ManifestValues: %v, Args: %v\n", defaultValues, args)
	}

	// translate parameter fit for terraformf parameters
	return getTerraformParameter(parameterValues)

}
