package command

import (
	"fmt"

	"github.com/TsuyoshiUshio/strikes/services/repository"

	"github.com/urfave/cli"
)

type SearchCommand struct {
}

func (p *SearchCommand) Search(c *cli.Context) error {
	keyword := c.Args().Get(0)
	fmt.Printf("Now searching...%s\n", keyword)
	packages, err := repository.GetPackages(keyword)
	if err != nil {
		return nil
	}
	for _, p := range *packages {
		fmt.Printf("%s\n", p.Name)
	}
	return nil
}
