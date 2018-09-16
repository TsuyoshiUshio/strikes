package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestVersion(t *testing.T) {
	Expected := "3.1.0"
	p := &Package{
		Releases: &[]Release{
			Release{
				Version: "1.0.0",
			},
			Release{
				Version: "2.0.0",
			},
			Release{
				Version: Expected,
			},
		},
	}
	latestVersion := p.LatestVersion()
	assert.Equal(t, Expected, latestVersion)
}

func TestAssetCircuitZipURL(t *testing.T) {
	ExpectedUrl := AssetBaseURL + "foo/" + "1.0.0/" + "circuit/" + "circuit.zip"
	p := &Package{
		Name: "foo",
		Releases: &[]Release{
			Release{
				Version: "1.0.0",
			},
		},
	}
	ActualUrl := p.GetCircuitZipURL()
	assert.Equal(t, ExpectedUrl, ActualUrl)
}
