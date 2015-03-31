package main

import (
	"fmt"
	"github.com/bbiskup/edify/commands"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "edify"
	app.Usage = "EDIFACT tool"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "download",
			Aliases: []string{"d"},
			Action: func(c *cli.Context) {
				url := c.Args().First()
				err := commands.Download(url)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Action: func(c *cli.Context) {
				what := c.Args().First()
				if len(what) == 0 {
					fmt.Println("No filename given")
					os.Exit(1)
				}
				commands.ParseSimpleDataElements(what)
			},
		},
	}

	app.Run(os.Args)
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
