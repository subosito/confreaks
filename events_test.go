package confreaks

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseEvents(t *testing.T) {
	b, err := loadFixture("events.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	events, err := ParseEvents(bytes.NewReader(b))
	if err != nil {
		t.Errorf("ParseEvents returned an error %s", err)
	}

	// values based on events.html
	exCount := 151
	count := len(events)
	if count != exCount {
		t.Errorf("There should be %d items, but we found %d items", exCount, count)
	}

	// test last item of events
	event := events[count-1]
	ex := Event{
		Title: "MountainWest RubyConf 2007",
		URL:   "http://confreaks.com/events/mwrc2007",
	}

	if !reflect.DeepEqual(event, ex) {
		t.Errorf("%#v is not equal with %ex", event, ex)
	}
}
