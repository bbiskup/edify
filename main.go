package main // import "github.com/bbiskup/edify"

import (
	"errors"
	"github.com/bbiskup/edify/commands"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
)

var (
	versionFlag = cli.StringFlag{
		Name:  "version",
		Value: "14B",
		Usage: "UNCE EDIFACT version",
	}
	specDirFlag = cli.StringFlag{
		Name:  "specdir, d",
		Value: "",
		Usage: "Specification directory (UNCE layout)",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "edify"
	app.Usage = "EDIFACT tool"
	app.EnableBashCompletion = true

	var err error

	app.Commands = []cli.Command{
		{
			Name:    "download_specs",
			Usage:   "Download specs from remote server",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) {
				// version: e.g. 14b
				version := c.Args().First()
				err = commands.DownloadSpecs(version)
			},
		},
		{
			Name:    "extract_specs",
			Usage:   "Extracts previously downloaded specs",
			Aliases: []string{"x"},
			Action: func(c *cli.Context) {
				// version: e.g. 14b
				version := c.Args().First()
				err = commands.ExtractSpecs(version)
			},
		},
		{
			Name:    "purge_specs",
			Usage:   "Purge previously extracted specs",
			Aliases: []string{"u"},
			Action: func(c *cli.Context) {
				purgeAll := c.Bool("all")
				// version: e.g. 14b
				version := c.Args().First()
				err = commands.PurgeSpecs(version, purgeAll)
			},

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "delete everything (including downloaded archives)",
				},
			},
		},
		{
			Name:    "parse",
			Usage:   "Parse a particular spec file",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) {
				fileNames := c.Args()
				err = commands.Parse(fileNames)
			},
		},

		{
			Name:    "query",
			Usage:   "Query a message part",
			Aliases: []string{"q"},
			Action: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					err = errors.New("Unexpected arguments")
					return
				}
				err = commands.Query(
					c.String("version"), c.String("specdir"),
					c.String("msg"), c.String("query"))
				log.Printf("##### " + c.String("dir"))
			},
			Flags: []cli.Flag{
				versionFlag,
				specDirFlag,
				cli.StringFlag{
					Name:  "msg, m",
					Value: "",
					Usage: "EDIFACT message file",
				},
				cli.StringFlag{
					Name:  "query, q",
					Value: "",
					Usage: "Query string",
				},
			},
		},
	}

	start := time.Now()

	app.Run(os.Args)
	if err != nil {
		log.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	end := time.Now()
	duration := end.Sub(start)

	log.Printf("Duration: %d ms", duration.Nanoseconds()/1e6)

	/*
		e := edi.NewElem("name1", "value1")
		s := edi.NewSeg("segname1")
		s.AddElem(e)

		m := edi.NewMessage("messagename1")
		m.AddSeg(s)

		i := edi.NewInterchange()
		i.AddMessage(m)

		fmt.Println(i)
	*/

}
