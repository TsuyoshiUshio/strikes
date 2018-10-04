package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Method to fetch the RepositoryAccessToken
// GET: repository.simplearchitect.club/api/assetserveruri // Repository URL
// In production, this endpoint will be protected by Azure AD
func GetRepositoryAccessToken() (*RepositoryAccessToken, error) {
	url := RepositoryBaseURL + "repositoryAccessToken"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var token RepositoryAccessToken
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func GetPackages(packageName string) (*[]Package, error) {
	resp, err := getRequest(RepositoryBaseURL+"packages?name="+packageName,
		fmt.Sprintf("Can not get package name: %s\n", packageName),
		fmt.Sprintf("Can not fetch the package: "+packageName))
	if err != nil {
		return nil, err
	}
	return NewSearchPackageFromJson(*resp)
}

func GetPackage(name string) (*Package, error) {
	resp, err := getRequest(RepositoryBaseURL+"package/name/"+name,
		fmt.Sprintf("Can not get package id: %s\n", name),
		fmt.Sprintf("Can not fetch the package by name: "+name))
	if err != nil {
		return nil, err
	}
	return NewPackageFromJson(*resp)
}

func getRequest(requestUrl, errorMessage, notFoundMessage string) (*[]byte, error) {
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Fatal(errorMessage)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		result := buf.Bytes()
		return &result, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(notFoundMessage)
	} else {
		buf := new(bytes.Buffer)
		log.Fatalf("Backend repoistory has some problem. StatusCode: %v ResponseBody: %v\n", resp.StatusCode, buf.String())
		return nil, nil
	}
}
