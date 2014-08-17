package confreaks

import (
	"encoding/json"
	"fmt"
	tdb "github.com/HouzuoGuo/tiedot/db"
	"reflect"
	"sort"
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

		col = db.Use(s)
	}

	switch s {
	case "events":
		col.Index([]string{"Title"})
	case "presentations":
		col.Index([]string{"Title"})
		col.Index([]string{"EUID"})
	}

	return col, nil
}

func SaveEvents(events []*Event) error {
	col, err := Use("events")
	if err != nil {
		return err
	}

	var q interface{}

	for i := range events {
		ev := events[i]
		ev.SumHash()

		json.Unmarshal([]byte(fmt.Sprintf(`{"eq": %q, "in": ["Title"], "limit": 1}`, ev.Title)), &q)
		result := make(map[int]struct{})
		tdb.EvalQuery(q, col, &result)

		data := docFromEvent(ev)

		if len(result) == 0 {
			if _, err := col.Insert(data); err == nil {
				log.WithField("title", ev.Title).Info("event added")
			}
		} else {
			for id := range result {
				doc, err := col.Read(id)
				if err == nil {
					if ev.Hash != doc["Hash"].(string) {
						if err := col.Update(id, data); err == nil {
							log.WithField("title", ev.Title).Info("event updated")
						}
					}
				}
			}
		}
	}

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

		ev = eventFromDoc(doc)
		return
	}

	return
}

func LoadEventPresentations(ev *Event) (err error) {
	var pcol *tdb.Col
	pcol, err = Use("presentations")
	if err != nil {
		return
	}

	var pq interface{}

	err = json.Unmarshal([]byte(fmt.Sprintf(`{"eq": %q, "in": ["EUID"]}`, ev.UUID)), &pq)
	if err != nil {
		return
	}

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

		p := presentationFromDoc(pdoc)
		ev.Presentations = append(ev.Presentations, p)
	}

	return
}

type byDate []*Event

func (d byDate) Len() int           { return len(d) }
func (d byDate) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d byDate) Less(i, j int) bool { return d[i].Date.After(d[j].Date) }

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

	sort.Sort(byDate(events))
	return
}

func SavePresentations(ev *Event, presentations []*Presentation) (err error) {
	col, err := Use("presentations")
	if err != nil {
		return
	}

	for i := range presentations {
		p := presentations[i]
		p.SumHash()

		var pq interface{}

		err = json.Unmarshal([]byte(fmt.Sprintf(`{"n":[{"eq": %q, "in": ["Title"]},{"eq": %q, "in": ["EUID"]}]}`, p.Title, ev.UUID)), &pq)
		if err != nil {
			return
		}

		presult := make(map[int]struct{})
		err = tdb.EvalQuery(pq, col, &presult)
		if err != nil {
			return
		}

		data := docFromPresentation(p)

		if len(presult) == 0 {
			if _, err := col.Insert(data); err == nil {
				log.WithField("title", p.Title).Info("presentation added")
			}
		} else {
			for pid := range presult {
				doc, err := col.Read(pid)
				if err == nil {
					if p.Hash != doc["Hash"].(string) {
						if err := col.Update(pid, data); err == nil {
							log.WithField("title", p.Title).Info("presentation updated")
						}
					} else {
						log.WithField("title", p.Title).Info("presentation ok")
					}
				}
			}
		}
	}

	return nil
}

func eventFromDoc(doc map[string]interface{}) *Event {
	ev := &Event{}
	ev.Title = doc["Title"].(string)
	ev.URL = doc["URL"].(string)
	ev.Hash = doc["Hash"].(string)
	ev.UUID = doc["UUID"].(string)
	ev.Count = int32(doc["Count"].(float64))

	return ev
}

func presentationFromDoc(doc map[string]interface{}) *Presentation {
	p := &Presentation{}
	p.Title = doc["Title"].(string)
	p.Description = doc["Description"].(string)
	p.VideoURL = doc["VideoURL"].(string)
	p.URL = doc["URL"].(string)
	p.EUID = doc["EUID"].(string)
	p.Hash = doc["Hash"].(string)

	pp := reflect.ValueOf(doc["Presenters"])
	if pp.Kind() == reflect.Slice {
		ps := make([]string, pp.Len())

		for i := 0; i < pp.Len(); i++ {
			name := pp.Index(i).Interface().(string)
			ps = append(ps, strings.TrimSpace(name))
		}

		p.Presenters = ps
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", doc["Recorded"].(string))
	if err == nil {
		p.Recorded = t
	}

	return p
}

func docFromPresentation(p *Presentation) map[string]interface{} {
	return map[string]interface{}{
		"Title":       p.Title,
		"Description": p.Description,
		"Presenters":  p.Presenters,
		"VideoURL":    p.VideoURL,
		"URL":         p.URL,
		"Recorded":    p.Recorded,
		"EUID":        p.EUID,
		"Hash":        p.Hash,
	}
}

func docFromEvent(ev *Event) map[string]interface{} {
	return map[string]interface{}{
		"UUID":  ev.UUID,
		"Title": ev.Title,
		"URL":   ev.URL,
		"Date":  ev.Date,
		"Hash":  ev.Hash,
		"Count": ev.Count,
	}
}
