package confreaks

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestEvent_ParseDetails(t *testing.T) {
	b, err := loadFixture("event-details.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	event := &Event{
		Title: "Pacific Northwest Scala 2013",
		URL:   "http://confreaks.com/events/pnws2013",
	}

	err = event.ParseDetails(bytes.NewReader(b))
	if err != nil {
		t.Errorf("ParseDetails returned an error %s", err)
	}

	exCount := 10
	if len(event.Presentations) != exCount {
		t.Errorf("event.Presentations length %d (not %d)", len(event.Presentations), exCount)
	}

	recorded, _ := time.Parse("02-Jan-06 15:04", "19-Oct-13 17:50")
	presentation := &Presentation{
		Title:      "Twitter: Taking Hadoop Realtime with Summingbird",
		Presenters: []string{"Sam Ritchie"},
		URL:        "http://confreaks.com/videos/2768-pnws2013-twitter-taking-hadoop-realtime-with-summingbird",
		Recorded:   recorded,
	}

	lastPresentation := event.Presentations[9]
	if !reflect.DeepEqual(lastPresentation, presentation) {
		t.Errorf("Last presentation %#v is not equal %#v", lastPresentation, presentation)
	}
}
