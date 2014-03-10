package confreaks

import (
	"fmt"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	sample := Event{
		Title: "Golang Bootcamp",
		URI:   "http://example.com/events/golang-bootcamp",
	}

	json, err := jsonMarshal(sample)
	if err != nil {
		t.Errorf("jsonMarshal returned an error %s", err)
	}

	ex := `{
  "title": "Golang Bootcamp",
  "uri": "http://example.com/events/golang-bootcamp"
}`

	if jsf := fmt.Sprintf("%s", json); jsf != ex {
		t.Errorf("Output of jsonMarshal is %s which is differ with expected %s", jsf, ex)
	}
}
