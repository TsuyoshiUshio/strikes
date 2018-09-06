package command

import "github.com/urfave/cli"

type PushCommand struct {
}

func (p *PushCommand) Push(c *cli.Context) error {
	// Read the manifest file
	// Read values file
	// create zip file for circuit
	// upload both zips to blob storage
	// rest api: create (or update) record
	// successfully uploaded
	return nil
}
