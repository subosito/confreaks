package main

import (
	"github.com/subosito/confreaks"
	"log"
)

func main() {
	var err error

	c := confreaks.NewConfreaks()
	err = c.FetchEvents()
	if err != nil {
		log.Fatal(err)
	}

	err = c.SaveIndex()
	if err != nil {
		log.Println(err)
	}

	for i := range c.Events {
		e := c.Events[i]

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
