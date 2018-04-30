package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

/// This sample test read/write from HCL file.
func main() {

	type Directories struct {
		Name        string `hcl:",key"`
		Description string
	}

	type Config struct {
		Variable []Directories
	}

	//	var conf map[string]map[string]string
	var conf Config
	b, err := ioutil.ReadFile("sample.hcl")
	if err != nil {
		panic(err)
	}

	fmt.Println("Decoding...\n")
	fmt.Println(string(b))
	err = hcl.Decode(&conf, string(b))
	fmt.Println(conf.Variable[0].Description)
	fmt.Println(conf.Variable[0].Name)
	if err != nil {
		panic(err)
	}
}
