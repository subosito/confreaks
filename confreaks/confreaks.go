package confreaks

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const baseURI = "http://confreaks.com/"

var indexFile = "index.json"

type Confreaks struct {
	URL    string   `json:"url"`
	Events []*Event `json:"events"`
}

func NewConfreaks() *Confreaks {
	return &Confreaks{
		URL: relativePath("events").String(),
	}
}

func NewConfreaksFromIndex() (c *Confreaks, err error) {
	f, err := ioutil.ReadFile(indexFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &c)
	if err != nil {
		return
	}

	return
}

func (c *Confreaks) FetchEvents() error {
	b, err := fetch(c.URL)
	if err != nil {
		return err
	}

	return c.ParseEvents(bytes.NewReader(b))
}

func (c *Confreaks) ParseEvents(r io.Reader) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var parse func(*html.Node)
	var extract func(*html.Node)
	var normalize func(*html.Node) Event

	normalize = func(n *html.Node) Event {
		if n.Attr != nil {
			for _, a := range n.Attr {
				if a.Key == "href" {
					eventUrl := relativePath(a.Val)

					return Event{
						Title: n.LastChild.Data,
						URL:   eventUrl.String(),
					}
				}
			}
		}

		return Event{}
	}

	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "strong" {
			for xc := n.FirstChild; xc != nil; xc = xc.NextSibling {
				event := normalize(xc)
				if event.Title != "" {
					c.Events = append(c.Events, &event)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	parse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "event-box" {
					extract(n)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}

	parse(doc)
	return nil
}

func (c *Confreaks) SaveIndex() error {
	var err error

	b, err := jsonMarshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(indexFile)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
