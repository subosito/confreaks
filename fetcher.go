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
	b, err := fetch(relativePath("events").String())
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

func relativePath(pathStr string) *url.URL {
	uri, _ := url.Parse("http://confreaks.com/")
	uri.Path = pathStr

	return uri
}

func fetch(uri string) ([]byte, error) {
	res, err := http.Get(uri)
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
