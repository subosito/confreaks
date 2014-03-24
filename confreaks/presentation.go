package confreaks

import (
	"bytes"
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

func (p *Presentation) Fetch() error {
	b, err := fetch(p.URL)
	if err != nil {
		return err
	}

	return p.ParseDetails(bytes.NewReader(b))
}

func (p *Presentation) ParseDetails(r io.Reader) error {
	var err error

	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var parse func(*html.Node)
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
						// case "video-presenters":
						// 	p.Presenters = extract(n)
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
	return nil
}

func (p *Presentation) DownloadVideo(dir string) error {
	err := downloadVideo(p.VideoURL, dir)
	if err != nil {
		return err
	}

	return nil
}
