package main

import (
	"github.com/codegangsta/cli"
	"github.com/subosito/confreaks/confreaks"
	"log"
	"os"
	"strings"
)

func init() {
	cli.AppHelpTemplate = `NAME:
  {{.Name}} ({{.Version}}) - {{.Usage}}

USAGE:
  {{.Name}} [GLOBAL OPTS] command [COMMAND OPTS] [ARGS...]

COMMANDS:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}

GLOBAL OPTIONS:
  {{range .Flags}}{{.}}
  {{end}}
`
}

func main() {
	var err error

	init := func() {
		c := confreaks.NewConfreaks()

		err = c.FetchEvents()
		if err != nil {
			log.Fatal(err)
		}

		err = c.SaveIndex()
		if err != nil {
			log.Println(err)
		}

		log.Println("confreaks events saved on index.json")
	}

	sync := func(pattern string) {
		c, err := confreaks.NewConfreaksFromIndex()
		if err != nil {
			log.Fatal(err)
		}

		for i := range c.Events {
			e := c.Events[i]

			if pattern != "" && !strings.Contains(e.Title, pattern) {
				continue
			}

			log.Printf("++ %s\n", e.Title)
			err = e.Fetch()
			if err != nil {
				log.Println(err)
			}

			err = e.ParsePresentations()
			if err != nil {
				log.Println(err)
			}

			for i := range e.Presentations {
				log.Printf(" +-- %s\n", e.Presentations[i].Title)
			}

			err = e.SaveIndex()
			if err != nil {
				log.Println(err)
			}
		}
	}

	app := cli.NewApp()
	app.Name = "confreaks"
	app.Usage = "confreaks on the command line"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize confreaks",
			Action: func(cc *cli.Context) {
				init()
			},
		},
		{
			Name:  "sync",
			Usage: "sync event/events [EVENT TITLE]",
			Action: func(cc *cli.Context) {
				sync(cc.Args().First())
			},
		},
	}

	app.Run(os.Args)
}
