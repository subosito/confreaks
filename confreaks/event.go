package confreaks

import (
	"bytes"
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Event struct {
	Title         string          `json:"title"`
	URL           string          `json:"url"`
	Presentations []*Presentation `json:"presentations,omitempty"`
}

func (e *Event) Fetch() error {
	b, err := fetch(e.URL)
	if err != nil {
		return err
	}

	return e.ParseDetails(bytes.NewReader(b))
}

func (e *Event) ParseDetails(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var presentations_selector = cascadia.MustCompile("div.video")
	for _, dom := range presentations_selector.MatchAll(doc) {
		p := &Presentation{}

		recorded_selector := cascadia.MustCompile(".recorded-at")
		recorded := recorded_selector.MatchFirst(dom)
		recorded_str := strings.TrimSpace(recorded.FirstChild.Data)
		recorded_at, err := time.Parse("02-Jan-06 15:04", recorded_str)
		if err == nil {
			p.Recorded = recorded_at
		}

		info_selector := cascadia.MustCompile(".main-info")
		info := info_selector.MatchFirst(dom)

		link_selector := cascadia.MustCompile(".title a")
		link := link_selector.MatchFirst(info)
		p.Title = link.LastChild.Data
		p.URL = relativePath(attrVal(link, "href")).String()

		presenters_selector := cascadia.MustCompile(".presenters a")
		presenters := []string{}
		for _, presenter := range presenters_selector.MatchAll(info) {
			presenters = append(presenters, presenter.LastChild.Data)
		}

		p.Presenters = presenters
		e.Presentations = append(e.Presentations, p)
	}

	return nil
}

func (e *Event) ParsePresentations() error {
	var wg sync.WaitGroup

	for i := range e.Presentations {
		p := e.Presentations[i]
		wg.Add(1)

		go func(p *Presentation) {
			defer wg.Done()

			for i := 0; ; i++ {
				err := p.Fetch()
				if err == nil {
					break
				} else {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(p)
	}

	wg.Wait()

	return nil
}

func (e *Event) Mkdir() error {
	err := os.MkdirAll(e.Title, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) SaveIndex() error {
	var err error

	err = e.Mkdir()
	if err != nil {
		return err
	}

	b, err := jsonMarshal(e)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(e.Title, "index.json"))
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) LoadIndex() (err error) {
	f, err := ioutil.ReadFile(filepath.Join(e.Title, indexFile))
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &e)
	if err != nil {
		return
	}

	return
}
