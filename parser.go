package confreaks

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
	"strings"
	"time"
)

func ParseEvents(r io.Reader) (events []*Event, err error) {
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
					eventURL := relativePath(a.Val)

					return Event{
						Title: n.LastChild.Data,
						URL:   eventURL.String(),
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
					events = append(events, &event)
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

func ParseEvent(r io.Reader, event *Event) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var presentationsSelector = cascadia.MustCompile("div.video")
	for _, dom := range presentationsSelector.MatchAll(doc) {
		p := &Presentation{}

		recordedSelector := cascadia.MustCompile(".recorded-at")
		recorded := recordedSelector.MatchFirst(dom)
		recordedStr := strings.TrimSpace(recorded.FirstChild.Data)
		recordedAt, err := time.Parse("02-Jan-06 15:04", recordedStr)
		if err == nil {
			p.Recorded = recordedAt
		}

		infoSelector := cascadia.MustCompile(".main-info")
		info := infoSelector.MatchFirst(dom)

		linkSelector := cascadia.MustCompile(".title a")
		link := linkSelector.MatchFirst(info)
		p.Title = strings.TrimSpace(link.LastChild.Data)
		p.URL = relativePath(attrVal(link, "href")).String()

		presentersSelector := cascadia.MustCompile(".presenters a")
		presenters := []string{}
		for _, presenter := range presentersSelector.MatchAll(info) {
			presenters = append(presenters, presenter.LastChild.Data)
		}

		p.Presenters = presenters
		event.Presentations = append(event.Presentations, p)
	}

	return nil
}

func ParsePresentation(r io.Reader, p *Presentation) error {
	var err error

	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var extract func(*html.Node) []string
	var normalize func(s string) string

	normalize = func(s string) string {
		au, err := url.Parse(s)
		if err != nil {
			return ""
		}

		switch au.Host {
		case "www.youtube.com":
			// http://www.youtube.com/embed/sVd4p6oKeUA
			// => http://www.youtube.com/watch?v=sVd4p6oKeUA

			v := url.Values{}
			v.Set("v", strings.Replace(au.Path, "/embed/", "", 1))

			au.Path = "watch"
			au.RawQuery = v.Encode()
		case "player.vimeo.com":
			// http://player.vimeo.com/video/40143060?badge=0
			// => http://vimeo.com/40143060

			au.Host = "vimeo.com"
			au.Path = strings.Replace(au.Path, "/video/", "", 1)
			au.RawQuery = ""
		}

		return au.String()
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

	presentationSelector := cascadia.MustCompile("div#primary-content")
	content := presentationSelector.MatchFirst(doc)

	titleSelector := cascadia.MustCompile(".video-title")
	title := titleSelector.MatchFirst(content)
	p.Title = strings.TrimSpace(title.LastChild.Data)

	descriptionSelector := cascadia.MustCompile(".video-abstract")
	description := descriptionSelector.MatchFirst(content)
	p.Description = strings.Join(extract(description), "\n")

	var videoSelector cascadia.Selector
	var video *html.Node

	videoSelector = cascadia.MustCompile("iframe")
	video = videoSelector.MatchFirst(content)

	if video == nil {
		videoSelector = cascadia.MustCompile("video source")
		video = videoSelector.MatchFirst(content)
	}

	if video != nil {
		p.VideoURL = normalize(attrVal(video, "src"))
	}

	return nil
}
