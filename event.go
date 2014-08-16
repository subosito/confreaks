package confreaks

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"strconv"
	"sync"
	"time"
)

type Event struct {
	ID               int             `json:"-"`
	Title            string          `json:"Title"`
	URL              string          `json:"URL"`
	Hash             string          `json:"Hash"`
	Date             time.Time       `json:"Date"`
	NumPresentations int             `json:"NumPresentations"`
	Presentations    []*Presentation `json:"presentations,omitempty"`
}

func (e *Event) SumHash() string {
	h := sha1.New()
	io.WriteString(h, e.Title)
	io.WriteString(h, e.URL)
	io.WriteString(h, e.Date.String())
	io.WriteString(h, strconv.Itoa(e.NumPresentations))

	e.Hash = fmt.Sprintf("%x", h.Sum(nil))
	return e.Hash
}

func (e *Event) FetchDetails() error {
	b, err := fetch(e.URL)
	if err != nil {
		log.WithFields(logrus.Fields{"event": e.Title, "error": err.Error()}).Info("fetch error")
		return err
	}

	return ParseEvent(bytes.NewReader(b), e)
}

func (e *Event) FetchPresentations() error {
	tasks := make(chan *Presentation, 12)

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for p := range tasks {
				for x := 0; ; x++ {
					log.WithField("presentation", p.Title).Debug("fetching")

					err := p.FetchDetails()
					if err == nil {
						p.EventTitle = e.Title
						log.WithField("presentation", p.Title).Debug("fetched")
						break
					} else {
						time.Sleep(200 * time.Millisecond)
						log.WithField("presentation", p.Title).Debug("re-fetching")
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
