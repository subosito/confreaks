package confreaks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func loadFixture(name string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("fixtures/%s", name))
}

func TestParseEventDetails(t *testing.T) {
	b, err := loadFixture("event-details.html")
	if err != nil {
		t.Errorf("loadFixture failed with error: %s", err)
	}

	event, err := ParseEventDetails(bytes.NewReader(b))
	if err != nil {
		t.Errorf("ParseEventDetails returned an error %s", err)
	}

	ex := EventDetail{
		Title:       "Pharmacist or a Doctor - What does your code base need?",
		VideoURL:    "http://www.youtube.com/watch?v=c-AuPryBSZs",
		Presenters:  "Pavan Sudarshan, Anandha Krishnan",
		Description: "You might know of every single code quality \u0026 metrics tool in the Ruby ecosystem and what they give you, but do you know:\nWhich metrics do you currently need?\nDo you really need them?\nHow do you make your team members own them?\nWait, there was a metaphor in the title\nWhile a pharmacist knows about pretty much every medicine out there and what it cures, its really a doctor who figures out what is required given the symptoms of a patient.\nJust like the vitals recorded for healthy adults, infants, pregnant women or an accident patient in an ICU changes, your code base needs different metrics in different contexts to help you uncover problems.\nTalk take aways\nThrough a lot of examples, the talk helps:\nIdentify the current state of your code base\nUnderstand different metrics and what do they tell you about your code base\nDrive your team towards continuously fixing problems uncovered by these metrics",
	}

	if !reflect.DeepEqual(event, ex) {
		t.Errorf("%#v is not equal with %ex", event, ex)
	}
}

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
		URI:   "http://confreaks.com/events/mwrc2007",
	}

	if !reflect.DeepEqual(event, ex) {
		t.Errorf("%#v is not equal with %ex", event, ex)
	}
}
