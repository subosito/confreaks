package confreaks

import (
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
	"strings"
	"time"
)

type Presentation struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Presenters  []string  `json:"presenters"`
	VideoURL    string    `json:"video-url"`
	URL         string    `json:"url"`
	Recorded    time.Time `json:"recorded"`
}

func ParsePresentation(r io.Reader) (p Presentation, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}

	var parse func(*html.Node)
	var extract func(*html.Node) []string
	var normalize func(s string) string

	normalize = func(s string) string {
		au, err := url.Parse(s)
		if err != nil {
			return ""
		}

		v := url.Values{}
		v.Set("v", strings.Replace(au.Path, "/embed/", "", 1))

		au.Path = "watch"
		au.RawQuery = v.Encode()
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

	parse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "iframe":
				for _, a := range n.Attr {
					if a.Key == "src" {
						p.VideoURL = normalize(a.Val)
					}
				}
			case "div":
				for _, a := range n.Attr {
					if a.Key == "class" {
						switch a.Val {
						case "video-title":
							p.Title = strings.TrimSpace(n.LastChild.Data)
						case "video-presenters":
							p.Presenters = extract(n)
						case "video-abstract":
							p.Description = strings.Join(extract(n), "\n")
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
