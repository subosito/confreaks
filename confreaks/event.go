package confreaks

import (
	"bytes"
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"encoding/json"
	"github.com/Sirupsen/logrus"
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

		log.WithFields(logrus.Fields{
			"event": e.Title,
			"error": err.Error(),
		}).Info("fetch error")

		return err
	}

	return e.ParseDetails(bytes.NewReader(b))
}

func (e *Event) ParseDetails(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var presentationsSelector = cascadia.MustCompile("div.video")
	for _, dom := range presentationsSelector.MatchAll(doc) {
		p := &Presentation{}

		recordedSelector := cascadia.MustCompile(".recorded-at")
		recorded := recordedSelector.MatchFirst(dom)
		recordedStr := strings.TrimSpace(recorded.FirstChild.Data)
		recordedAt, err := time.Parse("02-Jan-06 15:04", recordedStr)
		if err == nil {
			p.Recorded = recordedAt
		}

		infoSelector := cascadia.MustCompile(".main-info")
		info := infoSelector.MatchFirst(dom)

		linkSelector := cascadia.MustCompile(".title a")
		link := linkSelector.MatchFirst(info)
		p.Title = strings.TrimSpace(link.LastChild.Data)
		p.URL = relativePath(attrVal(link, "href")).String()

		presentersSelector := cascadia.MustCompile(".presenters a")
		presenters := []string{}
		for _, presenter := range presentersSelector.MatchAll(info) {
			presenters = append(presenters, presenter.LastChild.Data)
		}

		p.Presenters = presenters
		e.Presentations = append(e.Presentations, p)
	}

	return nil
}

func (e *Event) ParsePresentations() error {
	tasks := make(chan *Presentation, 12)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for p := range tasks {
				for x := 0; ; x++ {
					log.WithField("presentation", p.Title).Info("fetching")

					err := p.Fetch()
					if err == nil {
						log.WithField("presentation", p.Title).Info("fetched")
						break
					} else {
						time.Sleep(200 * time.Millisecond)
						log.WithField("presentation", p.Title).Info("re-fetching")
					}
				}
			}
			wg.Done()
		}()
	}

	for i := range e.Presentations {
		tasks <- e.Presentations[i]
	}

	close(tasks)

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

	f, err := os.Create(filepath.Join(e.Title, indexFile))
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
