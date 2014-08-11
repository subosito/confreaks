package confreaks

import (
	"bytes"
	"testing"
)

func TestEvent_ParsePresentations(t *testing.T) {
	b, err := loadFixture("event-details.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	event := &Event{
		Title: "Pacific Northwest Scala 2013",
		URL:   "http://confreaks.com/events/pnws2013",
	}

	err = ParseEvent(bytes.NewReader(b), event)
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	err = event.FetchPresentations()
	if err != nil {
		t.Errorf("ParsePresentations returned an error %s", err)
	}

	for i := range event.Presentations {
		if p := event.Presentations[i]; p.VideoURL == "" {
			t.Errorf("Unable to fetch VideoURL for %s", p.Title)
		}
	}
}
