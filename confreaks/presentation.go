package confreaks

import (
	"bytes"
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
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

	presentation_selector := cascadia.MustCompile("div#primary-content")
	content := presentation_selector.MatchFirst(doc)

	title_selector := cascadia.MustCompile(".video-title")
	title := title_selector.MatchFirst(content)
	p.Title = strings.TrimSpace(title.LastChild.Data)

	description_selector := cascadia.MustCompile(".video-abstract")
	description := description_selector.MatchFirst(content)
	p.Description = strings.Join(extract(description), "\n")

	var video_selector cascadia.Selector
	var video *html.Node

	video_selector = cascadia.MustCompile("iframe")
	video = video_selector.MatchFirst(content)

	if video == nil {
		video_selector = cascadia.MustCompile("video source")
		video = video_selector.MatchFirst(content)
	}

	if video != nil {
		p.VideoURL = normalize(attrVal(video, "src"))
	}

	return nil
}

func (p *Presentation) DownloadVideo(dir string) error {
	if p.VideoURL == "" {
		return errors.New(fmt.Sprintf("No Video URL for %q", p.Title))
	}

	err := downloadVideo(p.VideoURL, dir)
	if err != nil {
		return err
	}

	return nil
}
