package repository

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Method to fetch the RepositoryAccessToken
// GET: repository.simplearchitect.club/api/assetserveruri // Repository URL
// In production, this endpoint will be protected by Azure AD
func GetRepositoryAccessToken() (*RepositoryAccessToken, error) {
	url := "https://repository.simplearchitect.club/api/assetserveruri"

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
