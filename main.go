package main

import (
	"github.com/bbiskup/edify/commands"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "edify"
	app.Usage = "EDIFACT tool"
	app.EnableBashCompletion = true

	var err error

	app.Commands = []cli.Command{
		{
			// URL e.g.: http://www.unece.org/fileadmin/DAM/trade/untdid/d14b/d14b.zip
			Name:    "download_specs",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) {
				// version: e.g. d14b
				version := c.Args().First()
				err = commands.DownloadSpecs(version)
			},
		},
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) {
				fileName := c.Args().First()
				err = commands.ParseSimpleDataElements(fileName)
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
		e := edi.NewElement("name1", "value1")
		s := edi.NewSegment("segname1")
		s.AddElement(e)

		m := edi.NewMessage("messagename1")
		m.AddSegment(s)

		i := edi.NewInterchange()
		i.AddMessage(m)

		fmt.Println(i)
	*/

}
