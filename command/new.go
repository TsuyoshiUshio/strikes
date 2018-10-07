package command

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/TsuyoshiUshio/strikes/template/assets"
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
	packageList := assets.List(providerType)
	for i, template := range packageList {
		content, err := assets.ReadTemplateDescription("/" + providerType + "/" + template)
		if err != nil {
			log.Fatalf("Can not find Template Description for %s, error: %v\n", template, err)
			return nil
		}
		fmt.Printf("%d: %s:%s %s\n", i, template, adjustTabs(template), content)
	}

	fmt.Println("")
	fmt.Printf("Choose Template [0-%d]: ", len(packageList)-1)
	reader := bufio.NewReader(os.Stdin)
	line, prefix, err := reader.ReadLine()
	i, err := strconv.Atoi(string(line))
	if err != nil {
		fmt.Printf("Select the proper value. %s is not accepted. \n", line)
		return nil
	}
	if i > (len(packageList) - 1) {
		fmt.Printf("Select the proper value. %s is not accepted. \n", line)
		return nil
	}
	fmt.Printf("You typed: %s : %s : %v : %v \n", string(line), packageList[i], prefix, err)
	return nil
}

func adjustTabs(name string) string {
	if len(name) > 11 {
		return "\t"
	} else {
		return "\t\t"
	}
}
