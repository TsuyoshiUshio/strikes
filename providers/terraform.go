package providers

import (
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/hashicorp/hcl"
)

type Provider interface {
}

type TerraformProvider struct {
	Manifest  *config.Manifest
	TargetDir string
}

func (t *TerraformProvider) IsProviderCommandExists() error {
	// Check is there is terraform command is on the path.
	return nil
}

func (t *TerraformProvider) CreateResource(args []string) error {
	// Read and Parse Values.hcl
	// If there are parameters which is the same as the Values.hcl, override it.
	// Check if there is terraform command.
	// Execute the terraform init

	// If we could, show customer to the endpoint of the Azure Functions.

	return nil
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
	// m := make(map[string]string)
	// for _, key := range *getMapKeys(values) {

	// }
	return nil, nil
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

func getTerraformParameter(values *map[string]string) (*[]string, error) {
	return nil, nil
}
