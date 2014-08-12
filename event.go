package confreaks

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"sync"
	"time"
)

type Event struct {
	ID            int             `json:"-"`
	Title         string          `json:"title"`
	URL           string          `json:"url"`
	Presentations []*Presentation `json:"presentations,omitempty"`
}

func (e *Event) FetchDetails() error {
	b, err := fetch(e.URL)
	if err != nil {
		log.WithFields(logrus.Fields{
			"event": e.Title,
			"error": err.Error(),
		}).Info("fetch error")

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
						p.EventID = e.ID
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
