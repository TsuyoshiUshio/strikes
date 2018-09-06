package command

import (
	"log"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/config"
	"github.com/urfave/cli"
)

type PushCommand struct {
}

func (p *PushCommand) Push(c *cli.Context) error {

	packageDirBase := c.String("p")
	// Read the manifest file
	packageDir := filepath.Join(packageDirBase, "circuit", "manifest.yaml")
	// Read manifest file
	manifest, err := config.NewManifestFromFile(packageDir)
	if err != nil {
		log.Fatal(err)
	}
	err = manifest.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// create zip file for circuit

	// upload both zips to blob storage

	// rest api: create (or update) record

	// successfully uploaded

	return nil
}
