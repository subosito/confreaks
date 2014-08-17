package confreaks

import (
	"fmt"
	"time"
)

type Presentation struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Presenters  []string  `json:"presenters"`
	VideoURL    string    `json:"video-url"`
	URL         string    `json:"url"`
	Recorded    time.Time `json:"recorded"`
	Hash        string    `json:"Hash"`
	EUID        string    `json:"-"`
}

func (p *Presentation) FetchDetails() error {
	return fetchPresentation(p.URL, p)
}

func (p *Presentation) DownloadVideo(dir string) error {
	if p.VideoURL == "" {
		return fmt.Errorf("No Video URL for %q", p.Title)
	}

	err := downloadVideo(p.VideoURL, dir)
	if err != nil {
		return err
	}

	return nil
}
