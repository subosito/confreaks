package confreaks

import (
	"code.google.com/p/go.net/html"
	"io"
)

func ParseEvents(r io.Reader) (events []Event, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
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
					events = append(events, event)
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
	return
}
