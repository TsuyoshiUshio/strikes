package command

import (
	"log"
	"os"

	"github.com/TsuyoshiUshio/strikes/services/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type ListCommand struct {
}

func (s *ListCommand) List(c *cli.Context) error {
	instances, err := storage.List()
	if err != nil {
		log.Fatalf("Can not get instances from the PowerPlant. : %v", err)
		return err
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"InstanceID", "InstanceName", "ResourceGroup", "PackageName", "TimeStamp"})
	var list [][]string
	for _, instance := range *instances {
		list = append(list, []string{
			instance.InstanceID,
			instance.Name,
			instance.ResourceGroup,
			instance.PackageName,
			instance.TimeStamp.Format("2006/1/2 15:04:05"),
		})
	}
	table.SetBorder(false)
	table.AppendBulk(list)
	table.Render()
	return nil
}
