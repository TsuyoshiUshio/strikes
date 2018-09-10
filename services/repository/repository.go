package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

// Method to fetch the RepositoryAccessToken
// GET: repository.simplearchitect.club/api/assetserveruri // Repository URL
// In production, this endpoint will be protected by Azure AD
func GetRepositoryAccessToken() (*RepositoryAccessToken, error) {
	url := RepositoryBaseURL + "assetserveruri"

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

func GetPackage(packageName string) (*Package, error) {
	resp, err := http.Get(RepositoryBaseURL + "package?name=" + packageName)
	if err != nil {
		log.Fatalf("Can not get package name: %s\n", packageName)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return NewPackageFromJson(buf.Bytes())
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("Can not fetch the package: " + packageName)
	} else {
		buf := new(bytes.Buffer)
		log.Fatalf("Backend repoistory has some problem. StatusCode: %v ResponseBody: %v\n", resp.StatusCode, buf.String())
		return nil, nil
	}
}
