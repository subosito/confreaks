package confreaks

import (
	"bytes"
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Level = logrus.DebugLevel
}

func FetchEvents() ([]*Event, error) {
	b, err := fetch(relativePath("events").String())
	if err != nil {
		return []*Event{}, err
	}

	return ParseEvents(bytes.NewReader(b))
}
