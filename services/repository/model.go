package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TsuyoshiUshio/strikes/helpers"
)

const (
	RepositoryBaseURL = "https://strikesbackapp.azurewebsites.net/api/"                 // TODO : go back to the simplearchitect.club after DNS works.
	AssetBaseURL      = "https://strikesrepoe9eej5x3.blob.core.windows.net/repository/" // TODO : go back to the simplearchitect.club after DNS works.
)

type RepositoryAccessToken struct {
	StorageAccountName string
	ContainerName      string
	SASQueryParameter  string
}

type Package struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	ProjectPage string     `json:"projectPage"`
	ProjectRepo string     `json:"projectRepo"`
	CreatedTime time.Time  `json:"createdTime"`
	Releases    *[]Release `json:"releases"`
	IsDeleted   bool
}

type Release struct {
	Version      string    `json:"version"`
	ReleaseNote  string    `json:"releaseNote"`
	ProviderType string    `json:"providerType"`
	CreatedTime  time.Time `json:"createdTime"`
}

type SearchPackage struct {
	Id          string `json:"id"`
	Name        string
	Description string
	Author      string
	ProjectPage string
	ProjectRepo string
	CreatedTime time.Time
	Releases    string
	// Column for Azure Search soft delete
	IsDeleted bool
}

func NewPackageWithCurrentTime(
	name,
	description,
	author,
	projectPage,
	projectRepo,
	version,
	releaseNote,
	providerType string) *Package {
	// Time is provided on the Server side.
	return &Package{
		Name:        name,
		Description: description,
		Author:      author,
		ProjectPage: projectPage,
		ProjectRepo: projectRepo,
		Releases: &[]Release{
			Release{
				Version:      version,
				ReleaseNote:  releaseNote,
				ProviderType: providerType,
			},
		},
	}
}
func NewPackageFromJson(jsonBytes []byte) (*Package, error) {
	var p = string(jsonBytes)
	log.Printf("[DEBUG] Unmarshall Package from response: %s\n", p)
	var result Package
	err := json.Unmarshal(jsonBytes, &result)
	return &result, err
}

func NewSearchPackageFromJson(jsonBytes []byte) (*[]SearchPackage, error) {
	var result []SearchPackage
	err := json.Unmarshal(jsonBytes, &result)
	return &result, err
}

func (p *Package) LatestVersion() string {
	latest := ""

	for _, release := range *p.Releases {
		if release.Version > latest {
			latest = release.Version
		}
	}
	return latest
}

func (p *Package) Create() (*http.Response, error) {
	// Serialize object

	packageJson, _ := json.Marshal(p)
	// url := "https://repository.simplearchitect.club/api/package"
	// Post request
	resp, err := http.Post(RepositoryBaseURL+"package", helpers.ContentTypeApplicationJson, bytes.NewReader(packageJson))
	if err != nil {
		log.Fatalf("Can not create a resource to repository backend server %v \n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return resp, fmt.Errorf("HTTP Status code should be 201 (Created). Current Status Code is %d", resp.StatusCode)
	}
	return resp, nil
}

func (p *Package) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Package) GetCircuitZipURL() string {
	return AssetBaseURL + p.Name + "/" + p.LatestVersion() + "/" + "circuit/" + "circuit.zip"
}
