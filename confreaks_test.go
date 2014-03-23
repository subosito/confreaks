package confreaks

import (
	"bytes"
	"reflect"
	"testing"
)

func TestConfreaks_ParseEvents(t *testing.T) {
	b, err := loadFixture("events.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	c := NewConfreaks()
	err = c.ParseEvents(bytes.NewReader(b))
	if err != nil {
		t.Errorf("ParseEvents returned an error %s", err)
	}

	// values based on events.html
	exCount := 151
	count := len(c.Events)
	if count != exCount {
		t.Errorf("There should be %d items, but we found %d items", exCount, count)
	}

	// test last item of events
	event := c.Events[count-1]
	ex := &Event{
		Title: "MountainWest RubyConf 2007",
		URL:   "http://confreaks.com/events/mwrc2007",
	}

	if !reflect.DeepEqual(event, ex) {
		t.Errorf("%#v is not equal with %#v", event, ex)
	}
}

func TestNewConfreaksFromIndex(t *testing.T) {
	indexFile = "fixtures/index.json"

	defer func() {
		indexFile = "index.json"
	}()

	c, err := NewConfreaksFromIndex()
	if err != nil {
		t.Errorf("NewConfreaksFromIndex returned an error %s", err)
	}

	exCount := 153
	if evCount := len(c.Events); evCount != exCount {
		t.Errorf("c.Events length is %d, expected to be %d", evCount, exCount)
	}
}
