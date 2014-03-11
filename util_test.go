package confreaks

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func loadFixture(name string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("fixtures/%s", name))
}

func TestJsonMarshal(t *testing.T) {
	sample := Event{
		Title: "Golang Bootcamp",
		URL:   "http://example.com/events/golang-bootcamp",
	}

	json, err := jsonMarshal(sample)
	if err != nil {
		t.Errorf("jsonMarshal returned an error %s", err)
	}

	ex := `{
  "title": "Golang Bootcamp",
  "url": "http://example.com/events/golang-bootcamp",
  "presentations": null
}`

	if jsf := fmt.Sprintf("%s", json); jsf != ex {
		t.Errorf("Output of jsonMarshal is %s which is differ with expected %s", jsf, ex)
	}
}
