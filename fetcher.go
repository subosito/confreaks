package confreaks

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

func FetchEvents() (events []*Event, err error) {
	return fetchEvents()
}

func fetchEvents() (events []*Event, err error) {
	b, err := fetch("/events")
	if err != nil {
		return
	}

	return ParseEvents(bytes.NewReader(b))
}

func fetchEvent(u string, ev *Event) error {
	b, err := fetch(u)
	if err != nil {
		return err
	}

	return ParseEvent(bytes.NewReader(b), ev)
}

func fetchPresentation(u string, p *Presentation) error {
	b, err := fetch(u)
	if err != nil {
		return err
	}

	return ParsePresentation(bytes.NewReader(b), p)
}

func resolveURI(s string) (*url.URL, error) {
	uri, _ := url.Parse("http://confreaks.com/")

	urs, err := url.Parse(s)
	if err != nil {
		return &url.URL{}, err
	}

	// reset scheme and host
	urs.Scheme = ""
	urs.Host = ""

	return uri.ResolveReference(urs), nil
}

func fetch(u string) ([]byte, error) {
	uri, err := resolveURI(u)
	if err != nil {
		return []byte{}, err
	}

	res, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
