package confreaks

import (
	"bytes"
	"fmt"
	"testing"
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

	event.ParseDetails(bytes.NewReader(b))
	v, _ := jsonMarshal(event)
	fmt.Printf("%s", v)
}
