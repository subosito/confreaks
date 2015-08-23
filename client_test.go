package confreaks

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	assert.Equal(t, c.BaseURL.String(), BaseURL)
	assert.Equal(t, c.UserAgent, UserAgent)
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)
	c.UserAgent = "UA/confreaks"

	req, err := c.NewRequest("GET", "http://example.com/", nil)
	assert.Nil(t, err)
	assert.Equal(t, req.URL.String(), "http://example.com/")
	assert.Equal(t, req.Header.Get("User-Agent"), c.UserAgent)
	assert.Equal(t, req.Header.Get("Accept"), "application/json")
	assert.Equal(t, req.Header.Get("Accept-Charset"), "utf-8")
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)

	_, err := c.NewRequest("GET", ":", nil)
	assert.NotNil(t, err)

	erl, ok := err.(*url.Error)
	assert.True(t, ok)
	assert.NotNil(t, erl)
	assert.Equal(t, erl.Op, "parse")
}
