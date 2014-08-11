package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/subosito/confreaks"
	"log"
	"os"
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

	err = confreaks.OpenDB("_db_")
	if err != nil {
		log.Fatal(err)
	}

	init := func() {
		events, err := confreaks.FetchEvents()
		if err != nil {
			log.Fatal(err)
		}

		err = confreaks.SaveEvents(events)
		if err != nil {
			log.Println(err)
		}

		log.Println("confreaks events saved")
	}

	list := func() {
		events, err := confreaks.AllEvents()
		if err != nil {
			log.Println(err)
		}

		for i := range events {
			fmt.Println(events[i].Title)
		}
	}

	sync := func(pattern string) {
		var err error

		ev, _ := confreaks.OpenEvent(pattern)

		err = ev.FetchDetails()
		if err != nil {
			log.Println(err)
		}

		err = ev.FetchPresentations()
		if err != nil {
			log.Println(err)
		}

		err = confreaks.SavePresentations(ev.Presentations)
		if err != nil {
			log.Println(err)
		}

		for x := range ev.Presentations {
			log.Printf(" +-- %s\n", ev.Presentations[x].Title)
		}
	}

	download := func(pattern string) {
		var err error

		ev, _ := confreaks.OpenEvent(pattern)

		log.Printf("++ %s\n", ev.Title)

		for x := range ev.Presentations {
			p := ev.Presentations[x]

			log.Printf(" +-- Downloading %s\n", p.Title)
			err = p.DownloadVideo(ev.Title)
			if err != nil {
				log.Printf(" !! unable to download video for %q\n", p.Title)
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
			Name:  "list",
			Usage: "list available events",
			Action: func(cc *cli.Context) {
				list()
			},
		},
		{
			Name:  "sync",
			Usage: "sync event/events [EVENT TITLE]",
			Action: func(cc *cli.Context) {
				sync(cc.Args().First())
			},
		},
		{
			Name:  "download",
			Usage: "download event/events [EVENT TITLE]",
			Action: func(cc *cli.Context) {
				download(cc.Args().First())
			},
		},
	}

	app.Run(os.Args)
}
