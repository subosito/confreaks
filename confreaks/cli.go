package main

import (
	"github.com/subosito/confreaks"
	"log"
)

func main() {
	var err error

	events, err := confreaks.LoadEvents()
	if err != nil {
		log.Fatal(err)
	}

	for i := range events {
		e := events[i]

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
