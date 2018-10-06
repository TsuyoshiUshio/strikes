package command

import (
	"fmt"

	"github.com/urfave/cli"
)

type NewCommand struct {
}

func (s *NewCommand) New(c *cli.Context) error {
	templateName := c.Args().Get(0)
	providerType := c.Args().Get(1)
	packageName := c.Args().Get(2)

	if templateName == "" || providerType == "" || packageName == "" {
		fmt.Println("strikes new {templateName} {providerType} {packageName}")
		fmt.Println("example: strikes basic terraform hello-world")
		return nil
	}
	if providerType != "terraform" {
		fmt.Printf("ProviderType: %s is not supported.\n", providerType)
		fmt.Println("strikes new {templateName} {providerType} {packageName}")
		fmt.Println("example: strikes basic terraform hello-world")
		return nil
	}

	fmt.Println("Generating template...")
	fmt.Printf("TemplateName: %s\n", templateName)
	fmt.Printf("ProvierType: %s\n", providerType)
	fmt.Printf("PackageName: %s\n", packageName)

	return nil
}
