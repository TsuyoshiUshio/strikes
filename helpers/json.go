package helpers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func ReadJson(fileName string) (*map[string]string, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can not open file %s: %v", fileName, err)
	}
	defer jsonFile.Close()
	jsonByte, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	err = json.Unmarshal(jsonByte, &result)
	return &result, err
}
