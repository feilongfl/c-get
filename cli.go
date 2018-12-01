package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"sort"
)

func cliRecv() {
	app := cli.NewApp()
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		//{
		//	Name:    "list-chapter",
		//	Aliases: []string{"lc"},
		//	Usage:   "list comic chapter",
		//	Action: func(c *cli.Context) error {
		//		//return listChapter(c.Args().First())
		//		return nil
		//	},
		//},
		//{
		//	Name:    "list-image",
		//	Aliases: []string{"li"},
		//	Usage:   "list comic image",
		//	Action: func(c *cli.Context) error {
		//		return nil
		//	},
		//},
		//{
		//	Name:    "info-comic",
		//	Aliases: []string{"ic"},
		//	Usage:   "get comic info",
		//	Action: func(c *cli.Context) error {
		//		return infoComic(c.Args().First())
		//	},
		//},
		{
			Name:    "download-comic",
			Aliases: []string{"dc"},
			Usage:   "download all comic images",
			Action: func(c *cli.Context) error {
				return infoComic(c.Args().First())
			},
		},
		//{
		//	Name:    "download comic",
		//	Aliases: []string{"dc"},
		//	Usage:   "list comic chapter",
		//	Action: func(c *cli.Context) error {
		//		return nil
		//	},
		//},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
