package confreaks

import (
	"encoding/json"
	"fmt"
	tdb "github.com/HouzuoGuo/tiedot/db"
	"reflect"
	"strings"
	"time"
)

var db *tdb.DB

func OpenDB(dbPath string) error {
	var err error

	db, err = tdb.OpenDB(dbPath)
	if err != nil {
		return err
	}

	return nil
}

func DB() *tdb.DB {
	return db
}

func CloseDB() error {
	return db.Close()
}

func Use(s string) (*tdb.Col, error) {
	var col *tdb.Col

	col = db.Use(s)
	if col == nil {
		err := db.Create(s)
		if err != nil {
			return nil, err
		}
	}

	return col, nil
}

func SaveEvents(events []*Event) error {
	col, err := Use("events")
	if err != nil {
		return err
	}

	if col.ApproxDocCount() == 0 {
		for i := range events {
			ev := events[i]
			log.WithField("title", ev.Title).Info("event added")
			col.Insert(map[string]interface{}{
				"Title": ev.Title,
				"URL":   ev.URL,
			})
		}

		col.Index([]string{"Title"})
		return nil
	}

	var q interface{}

	for i := range events {
		ev := events[i]
		json.Unmarshal([]byte(fmt.Sprintf(`{"eq": %q, "in": ["Title"], "limit": 1}`, ev.Title)), &q)
		result := make(map[int]struct{})
		tdb.EvalQuery(q, col, &result)

		if len(result) == 0 {
			log.WithField("title", ev.Title).Info("event added")
			col.Insert(map[string]interface{}{
				"Title": ev.Title,
				"URL":   ev.URL,
			})
		}
	}

	col.Index([]string{"Title"})
	return nil
}

func OpenEvent(title string) (ev *Event, err error) {
	col, err := Use("events")
	if err != nil {
		return
	}

	var q interface{}
	json.Unmarshal([]byte(fmt.Sprintf(`{"eq": %q, "in": ["Title"], "limit": 1}`, title)), &q)
	result := make(map[int]struct{})

	err = tdb.EvalQuery(q, col, &result)
	if err != nil {
		return
	}

	for id := range result {
		var doc map[string]interface{}

		doc, err = col.Read(id)
		if err != nil {
			return
		}

		ev = &Event{}
		ev.Title = doc["Title"].(string)
		ev.URL = doc["URL"].(string)
		ev.ID = id

		var pcol *tdb.Col
		pcol, err = Use("presentations")
		if err != nil {
			return
		}

		var pq interface{}

		json.Unmarshal([]byte(fmt.Sprintf(`{"eq": %d, "in": ["EventID"]}`, id)), &pq)

		presult := make(map[int]struct{})
		err = tdb.EvalQuery(pq, pcol, &presult)
		if err != nil {
			return
		}

		for pid := range presult {
			var pdoc map[string]interface{}

			pdoc, err = pcol.Read(pid)
			if err != nil {
				return
			}

			p := &Presentation{}
			p.Title = pdoc["Title"].(string)
			p.Description = pdoc["Description"].(string)
			p.VideoURL = pdoc["VideoURL"].(string)
			p.URL = pdoc["URL"].(string)
			p.EventID = int(pdoc["EventID"].(float64))

			pp := reflect.ValueOf(pdoc["Presenters"])
			if pp.Kind() == reflect.Slice {
				ps := make([]string, pp.Len())

				for i := 0; i < pp.Len(); i++ {
					name := pp.Index(i).Interface().(string)
					ps = append(ps, strings.TrimSpace(name))
				}

				p.Presenters = ps
			}

			tf := "2006-01-02T15:04:05Z"
			tt, _ := time.Parse(tf, pdoc["Recorded"].(string))
			p.Recorded = tt

			ev.Presentations = append(ev.Presentations, p)
		}

		return
	}

	return
}

func AllEvents() (events []*Event, err error) {
	col, err := Use("events")
	if err != nil {
		return
	}

	col.ForEachDoc(func(id int, b []byte) bool {
		ev := &Event{}
		err := json.Unmarshal(b, ev)
		if err != nil {
			return false
		}

		events = append(events, ev)
		return true
	})

	return
}

func SavePresentations(presentations []*Presentation) error {
	col, err := Use("presentations")
	if err != nil {
		return err
	}

	for i := range presentations {
		p := presentations[i]
		log.WithField("title", p.Title).Info("presentation added")
		col.Insert(map[string]interface{}{
			"Title":       p.Title,
			"Description": p.Description,
			"Presenters":  p.Presenters,
			"VideoURL":    p.VideoURL,
			"URL":         p.URL,
			"Recorded":    p.Recorded,
			"EventID":     p.EventID,
		})
	}

	col.Index([]string{"EventID"})
	return nil
}
