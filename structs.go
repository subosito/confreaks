package confreaks

import (
	"strings"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"display_name"`
	Conference  string    `json:"conference"`
	ShortCode   string    `json:"short_code"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
}

type Presenter struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AkaName       string `json:"aka_name"`
	TwitterHandle string `json:"twitter_handle"`
}

type Video struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Image        string      `json:"image"`
	Slug         string      `json:"slug"`
	RecordedAt   *time.Time  `json:"recorded_at"`
	Event        string      `json:"event"`
	Rating       string      `json:"rating"`
	Abstract     string      `json:"abstract"`
	PostDate     *time.Time  `json:"post_date"`
	AnnounceDate *time.Time  `json:"announce_date"`
	Host         string      `json:"host"`
	EmbedCode    string      `json:"embed_code"`
	Views        int         `json:"views"`
	ViewsLast7   int         `json:"views_last_7"`
	ViewsLast30  int         `json:"views_last_30"`
	License      string      `json:"license"`
	Attribution  string      `json:"attribution"`
	Presenters   []Presenter `json:"presenters"`
}

func (v Video) URL() string {
	if v.Host == "youtube" && v.EmbedCode != "" {
		return "https://www.youtube.com/watch?v=" + v.EmbedCode
	}

	return "-"
}

func (v Video) PresentersNames() string {
	s := []string{}
	for _, p := range v.Presenters {
		f := []string{strings.TrimSpace(p.FirstName), strings.TrimSpace(p.LastName)}
		s = append(s, strings.Join(f, " "))
	}

	return strings.Join(s, ", ")
}
