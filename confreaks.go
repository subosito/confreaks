package confreaks

import (
	"net/http"
)

type Confreaks struct {
	client *Client
}

func New(httpClient *http.Client) *Confreaks {
	return &Confreaks{NewClient(httpClient)}
}

func (c *Confreaks) Client() *Client {
	return c.client
}

func (c *Confreaks) Events() ([]Event, error) {
	v := []Event{}
	err := c.doParse(&v, "events")
	return v, err
}

func (c *Confreaks) Videos(s string) ([]Video, error) {
	v := []Video{}
	err := c.doParse(&v, "events", s, "videos")
	return v, err
}

func (c *Confreaks) doParse(v interface{}, parts ...string) error {
	u := c.client.Path(parts...)

	req, err := c.client.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.client.Do(req, v)
	if err != nil {
		return err
	}

	return nil
}
