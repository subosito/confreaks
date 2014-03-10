package confreaks

import (
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
	"strings"
)

type Event struct {
	Title string `json:"title"`
	URI   string `json:"uri"`
}

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
					eventUrl, _ := relativePath(a.Val)

					return Event{
						Title: n.LastChild.Data,
						URI:   eventUrl.String(),
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

type EventDetail struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Presenters  string `json:"presenters"`
	VideoURL    string `json:"video-url"`
}

func ParseEventDetails(r io.Reader) (ev EventDetail, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}

	var parse func(*html.Node)
	var extract func(*html.Node) []string
	var normalize func(s string) (string, error)

	normalize = func(s string) (string, error) {
		au, err := url.Parse(s)
		if err != nil {
			return "", err
		}

		v := url.Values{}
		v.Set("v", strings.Replace(au.Path, "/embed/", "", 1))

		au.Path = "watch"
		au.RawQuery = v.Encode()
		return au.String(), nil
	}

	extract = func(n *html.Node) []string {
		texts := []string{}

		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			if ch.LastChild != nil {
				texts = append(texts, ch.LastChild.Data)
			}
		}

		return texts
	}

	parse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "iframe":
				for _, a := range n.Attr {
					if a.Key == "src" {
						au, _ := normalize(a.Val)
						ev.VideoURL = au
					}
				}
			case "div":
				for _, a := range n.Attr {
					if a.Key == "class" {
						switch a.Val {
						case "video-title":
							ev.Title = strings.TrimSpace(n.LastChild.Data)
						case "video-presenters":
							ev.Presenters = strings.Join(extract(n), ", ")
						case "video-abstract":
							ev.Description = strings.Join(extract(n), "\n")
						}
					}
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
