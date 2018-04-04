package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Lightning Strikes"
	app.Usage = "The Azure Functions Package management tool"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Initialize strikes. Initialize your Power Plant.",
			Action: func(c *cli.Context) error {
				fmt.Println("Creating ...")
				time.Sleep(2000 * time.Millisecond)
				fmt.Println("Power Plant is succesfully installed.")
				return nil
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
