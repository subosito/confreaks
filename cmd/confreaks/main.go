package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/subosito/confreaks"
)

func main() {
	c := confreaks.New(nil)

	events := func() {
		s, err := c.Events()
		if err != nil {
			fmt.Println("!: %s", err)
			return
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprint(w, "NO\tDATE EVENT\tSHORT CODE\tEVENT TITLE\n")

		for i, v := range s {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, v.StartAt.Format("Jan 02, 2006"), v.ShortCode, v.DisplayName)
		}
		fmt.Fprintln(w)
		w.Flush()
	}

	videos := func(str string) {
		s, _ := c.Videos(str)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "NO\tVIDEO TITLE\tPRESENTERS\tVIDEO URL")

		for i, v := range s {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i+1, v.Title, v.PresentersNames(), v.URL())
		}

		fmt.Fprintln(w)
		w.Flush()
	}

	app := cli.NewApp()
	app.Name = "confreaks"
	app.Usage = "confreaks on the command line"
	app.Version = "2.0.0"
	app.Commands = []cli.Command{
		{
			Name:  "events",
			Usage: "List all events",
			Action: func(cc *cli.Context) {
				events()
			},
		},
		{
			Name:  "videos",
			Usage: "List event's videos",
			Action: func(cc *cli.Context) {
				videos(cc.Args().First())
			},
		},
	}

	app.Run(os.Args)
}
