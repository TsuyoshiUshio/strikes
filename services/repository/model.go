package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TsuyoshiUshio/strikes/helpers"
)

const (
	RepositoryBaseURL = "https://repository.simplearchitect.club/api/"
)

type RepositoryAccessToken struct {
	StorageAccountName string
	ContainerName      string
	SASQueryParameter  string
}

type Package struct {
	Id          string `json:"id"`
	Name        string
	Description string
	Author      string
	ProjectPage string
	ProjectRepo string
	CreatedTime time.Time
	Releases    *[]Release
	IsDeleted   bool
}

type Release struct {
	Version      string
	ReleaseNote  string
	ProviderType string
	CreatedTime  time.Time
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
	var result Package
	err := json.Unmarshal(jsonBytes, &result)
	return &result, err
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
