package command

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/TsuyoshiUshio/strikes/template/assets"
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
	builder.Append(ui.NewReleaseNoteProcess(os.Stdin))
	builder.Append(ui.NewZipFileNameProcess(os.Stdin))
	process := builder.Build()
	parameter := ui.PackageParameter{}
	result, err := ui.Execute(process, parameter)
	if err != nil {
		return err
	}
	param := result.(ui.PackageParameter)

	param.AzureFunctionsTemplate = "{{.AzureFunctionsTemplate}}"
	param.ResourceGroupTemplate = "{{.ResourceGroupTemplate}}"
	param.EnvironmentBaseNameTemplate = "{{.EnvironmentBaseNameTemplate}}"

	fmt.Println("")
	fmt.Printf("TemplateDir: %s\nPackageName: %s\nDescription: %s\nAuthor: %s\nProjectPage: %s\nProjectRepo: %s\nReleaseNote: %s\nZipFileName: %s\n",
		param.TemplateDirPath,
		param.PackageName,
		param.Description,
		param.Author,
		param.ProjectPage,
		param.ProjectRepo,
		param.ReleaseNote,
		param.ZipFileName,
	)
	fmt.Println("")
	generateTempalte(&param)
	fmt.Println("%s template has beeen generated.", param.TemplateDirPath)
	return nil
}

func generateTempalte(parameter *ui.PackageParameter) error {
	// create a directry with the PackageName
	err := helpers.CreateDirIfNotExist(filepath.Join(".", parameter.PackageName))
	if err != nil {
		fmt.Println("Can not create the directory: %sf\n", filepath.Join(".", parameter.PackageName))
		return err
	}
	err = helpers.CreateDirIfNotExist(filepath.Join(".", parameter.PackageName, "circuit"))
	if err != nil {
		fmt.Println("Can not create the directory: %sf\n", filepath.Join(".", parameter.PackageName, "circuit"))
		return err
	}
	err = helpers.CreateDirIfNotExist(filepath.Join(".", parameter.PackageName, "package"))
	if err != nil {
		fmt.Println("Can not create the directory: %sf\n", filepath.Join(".", parameter.PackageName, "package"))
		return err
	}
	// Iterate with the File with the TmeplateDirPath
	files, _ := assets.Read(parameter.TemplateDirPath)
	d, _ := (*files).Readdir(0)
	for _, f := range d {
		file, err := os.Create(filepath.Join(".", parameter.PackageName, "circuit", f.Name()))
		if err != nil {
			fmt.Println("Can not create the file %s\n", filepath.Join(".", parameter.PackageName, "circuit", f.Name()))
			return err
		}
		source, _ := assets.Read(parameter.TemplateDirPath + "/" + f.Name())
		content, err := ioutil.ReadAll(*source)
		tmpl, err := template.New(f.Name()).Parse(string(content))
		if err != nil {
			fmt.Println("Can not create the template. %v\n", err)
			return err
		}
		err = tmpl.Execute(file, parameter)
		if err != nil {
			fmt.Println("Can not templating the file. : %s\n", f.Name())
			return err
		}
		file.Close()
	}
	// Write a file with templating

	return nil
}
