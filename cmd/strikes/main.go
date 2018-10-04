package main

import (
	"os"

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
			Action:  (&command.InstallCommand{}).Install,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search Lightning Strikes Package.",
			Action:  (&command.SearchCommand{}).Search,
		},
	}
	a := os.Args
	app.Run(a)
}
