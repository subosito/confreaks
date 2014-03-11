package confreaks

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"io"
	"strings"
	"time"
)

type Event struct {
	Title         string         `json:"title"`
	URL           string         `json:"url"`
	Presentations []Presentation `json:"presentations"`
}

func (e *Event) ParseDetails(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var presentations_selector = cascadia.MustCompile("div.video")
	for _, dom := range presentations_selector.MatchAll(doc) {
		presentation := Presentation{}

		recorded_selector := cascadia.MustCompile(".recorded-at")
		recorded := recorded_selector.MatchFirst(dom)
		recorded_str := strings.TrimSpace(recorded.FirstChild.Data)
		recorded_at, err := time.Parse("02-Jan-06 15:04", recorded_str)
		if err == nil {
			presentation.Recorded = recorded_at
		}

		info_selector := cascadia.MustCompile(".main-info")
		info := info_selector.MatchFirst(dom)

		link_selector := cascadia.MustCompile(".title a")
		link := link_selector.MatchFirst(info)
		presentation.Title = link.LastChild.Data
		presentation.URL = relativePath(attrVal(link, "href")).String()

		presenters_selector := cascadia.MustCompile(".presenters a")
		presenters := []string{}
		for _, presenter := range presenters_selector.MatchAll(info) {
			presenters = append(presenters, presenter.LastChild.Data)
		}

		presentation.Presenters = presenters
		e.Presentations = append(e.Presentations, presentation)
	}

	return nil
}
