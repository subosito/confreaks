package confreaks

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestConfreaks_Events(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/events.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("fixtures/events.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	sURL, _ := url.Parse(server.URL)

	client := New(nil)
	client.Client().BaseURL = sURL

	v, err := client.Events()
	assert.Nil(t, err)
	assert.Equal(t, 249, len(v))

	z := v[1]
	d := timeParse("2015-08-07T00:00:00.000Z")
	assert.Nil(t, err)

	w := Event{
		ID:          273,
		DisplayName: "NebraskaJS 2015",
		Conference:  "NebraskaJS",
		ShortCode:   "nebraskajs2015",
		StartAt:     d,
		EndAt:       d,
	}
	assert.Equal(t, w, z)
}

func TestConfreaks_Videos(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/events/rustcamp2015/videos.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("fixtures/videos.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	sURL, _ := url.Parse(server.URL)

	client := New(nil)
	client.Client().BaseURL = sURL

	v, err := client.Videos("rustcamp2015")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(v))

	z := v[1]
	d := timeParse("2015-08-01T16:00:00.000Z")
	p := timeParse("2015-08-14T08:30:00.000Z")
	w := Video{
		ID:           5782,
		Title:        "Writing High Performance Async IO Apps",
		Image:        "http://s3-us-west-2.amazonaws.com/confreaks-tv3/production/videos/images/000/005/782/Capture-original.PNG?1439562736",
		Slug:         "rustcamp2015-writing-high-performance-async-io-apps",
		RecordedAt:   &d,
		Event:        "RustCamp 2015",
		Rating:       "Everyone",
		Abstract:     "The talk covers how to use MIO (a lightweight non-blocking IO event loop in Rust) to write fast IO applications with Rust. It will introduce MIO's network & timeout APIs and show how to use them to create a network server. The talk will then discuss some strategies for using MIO in a multithreaded environment.",
		PostDate:     &p,
		AnnounceDate: nil,
		Host:         "youtube",
		EmbedCode:    "CjQjEMw-snk",
		Views:        1122,
		ViewsLast7:   0,
		ViewsLast30:  0,
		License:      "cc-by-sa-3.0",
		Attribution:  "",
		Presenters: []Presenter{
			Presenter{
				FirstName:     "Carl",
				LastName:      "Lerche",
				AkaName:       "",
				TwitterHandle: "carllerche",
			},
		},
	}

	assert.Equal(t, w, z)
}

func timeParse(str string) time.Time {
	d, _ := time.Parse(time.RFC3339, str)
	return d
}
