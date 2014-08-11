package main

import (
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
		c := confreaks.NewConfreaks()

		err = c.FetchEvents()
		if err != nil {
			log.Fatal(err)
		}

		err = c.Save()
		if err != nil {
			log.Println(err)
		}

		log.Println("confreaks events saved")
	}

	sync := func(pattern string) {
		var err error

		c := confreaks.NewConfreaks()
		e, _ := c.GetEvent(pattern)

		err = e.Fetch()
		if err != nil {
			log.Println(err)
		}

		err = e.ParsePresentations()
		if err != nil {
			log.Println(err)
		}

		for x := range e.Presentations {
			log.Printf(" +-- %s\n", e.Presentations[x].Title)
		}
	}

	// download := func(pattern string) {
	// 	var err error

	// 	events := index(pattern)

	// 	for i := range events {
	// 		e := events[i]

	// 		log.Printf("++ %s\n", e.Title)
	// 		err = e.LoadIndex()
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		for x := range e.Presentations {
	// 			p := e.Presentations[x]

	// 			log.Printf(" +-- Downloading %s\n", p.Title)
	// 			err = p.DownloadVideo(e.Title)
	// 			if err != nil {
	// 				log.Printf(" !! unable to download video for %q\n", p.Title)
	// 			}
	// 		}
	// 	}
	// }

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
		// {
		// 	Name:  "list",
		// 	Usage: "list available events",
		// 	Action: func(cc *cli.Context) {
		// 		events := index("")
		// 		for i := range events {
		// 			fmt.Println(events[i].Title)
		// 		}
		// 	},
		// },
		{
			Name:  "sync",
			Usage: "sync event/events [EVENT TITLE]",
			Action: func(cc *cli.Context) {
				sync(cc.Args().First())
			},
		},
		// {
		// 	Name:  "download",
		// 	Usage: "download event/events [EVENT TITLE]",
		// 	Action: func(cc *cli.Context) {
		// 		download(cc.Args().First())
		// 	},
		// },
	}

	app.Run(os.Args)
}
