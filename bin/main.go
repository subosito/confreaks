package main

import (
	"fmt"
	tm "github.com/buger/goterm"
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

		box := tm.NewBox(50|tm.PCT, 75|tm.PCT, 0)
		et := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprint(et, "NO\tDATE EVENT\tEVENT TITLE\n")

		for i := range events {
			ev := events[i]
			fmt.Fprintf(et, "%d\t%s\t%s\n", i+1, ev.Date.Format("Jan 02, 2006"), ev.Title)
		}

		fmt.Fprintf(box, et.String())
		tm.Println(box)
		tm.Flush()
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

		fmt.Println(ev.Title)

		err = os.MkdirAll(ev.Title, 0755)
		if err != nil {
			log.Fatal(err.Error())
		}

		for x := range ev.Presentations {
			p := ev.Presentations[x]

			log.Printf(" +-- Downloading %s (%s)\n", p.Title, p.VideoURL)
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
