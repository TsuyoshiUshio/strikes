package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TsuyoshiUshio/strikes/command"
	"github.com/TsuyoshiUshio/strikes/helpers"
	"github.com/urfave/cli"
)

func main() {
	helpers.SetUpLogger()
	app := cli.NewApp()
	app.Name = "Lightning Strikes"
	app.Usage = "The Azure Functions Package management tool"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize strikes. Initialize your Power Plant.",
			Action:  (&command.InitCommand{}).Initialize,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "location, l",
					Value: "japaneast",
					Usage: "Specify location for the Strikes Resources",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Force remove current config and PowerPlant strage account if it is specified.",
				},
			},
		},
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push a package to repository server",
			Action:  (&command.PushCommand{}).Push,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "package, p",
					Value: ".",
					Usage: "Specify directory of the Strikes Package",
				},
			},
		},
		{
			Name:    "install",
			Aliases: []string{"in"},
			Usage:   "Install Lightning Strikes Package.",
			Action: func(c *cli.Context) error {
				fmt.Println("Install package ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Extracting ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Provisioning Storage Account ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Provisioning CosmosDB ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Provisioning Function App ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Deploy function to the Function App ...")
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("Ready!")
				fmt.Println("Please send message to this endpoint")
				fmt.Println("https://functions.azure.com/api/product")
				fmt.Println("Happy Strikes!")
				return nil
			},
		},
	}
	a := os.Args
	app.Run(a)
}
