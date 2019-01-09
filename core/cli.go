package core

import (
	"c-get/source"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"sort"
)

var _version_ = "v1.0.0"
var _commit_ = "manual"

func CliRecv() {
	app := cli.NewApp()
	app.Version = _version_
	app.Author = "feilong"
	app.Email = "feilongphone@gmail.com"
	app.Description = _commit_

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
				return source.InfoComic(c.Args().First())
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
