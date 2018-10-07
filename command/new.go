package command

import (
	"fmt"
	"os"

	"github.com/TsuyoshiUshio/strikes/ui"
	"github.com/urfave/cli"
)

type NewCommand struct {
}

func (s *NewCommand) New(c *cli.Context) error {
	providerType := c.Args().Get(0)

	if providerType == "" {
		fmt.Println("strikes new {providerType}")
		fmt.Println("example: strikes new terraform")
		return nil
	}
	if providerType != "terraform" {
		fmt.Printf("ProviderType: %s is not supported.\n", providerType)
		fmt.Println("strikes new {templateName} {providerType} {packageName}")
		fmt.Println("example: strikes basic terraform hello-world")
		return nil
	}

	// feature
	// user specify the provider type then choose the number.
	fmt.Println("")
	fmt.Println("Strikes Package Generator")
	fmt.Println("")

	builder := ui.NewProcessBuilder()
	builder.Append(ui.NewChooseTemplateProcess(providerType, os.Stdin))
	builder.Append(ui.NewPackageNameProcess(os.Stdin))
	builder.Append(ui.NewDescriptionProcess(os.Stdin))
	builder.Append(ui.NewAuthorProcess(os.Stdin))
	builder.Append(ui.NewProjectPageProcess(os.Stdin))
	builder.Append(ui.NewProjectRepoProcess(os.Stdin))
	process := builder.Build()
	parameter := ui.PackageParameter{}
	result, err := ui.Execute(process, parameter)
	if err != nil {
		return err
	}
	param := result.(ui.PackageParameter)
	fmt.Println("")
	fmt.Printf("TemplateDir: %s\nPackageName: %s\nDescription: %s\nAuthor: %s\nProjectPage: %s\nProjectRepo: %s\n",
		param.TemplateDirPath,
		param.PackageName,
		param.Description,
		param.Author,
		param.ProjectPage,
		param.ProjectRepo,
	)
	fmt.Println("")
	return nil
}
