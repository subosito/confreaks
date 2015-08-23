package confreaks

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	BaseURL   = "http://confreaks.tv/api/v1"
	UserAgent = "github.com/subosito/confreaks"
)

type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u, _ := url.Parse(BaseURL)

	return &Client{
		client:    httpClient,
		UserAgent: UserAgent,
		BaseURL:   u,
	}
}

func (c *Client) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.validResponse(resp) && v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

func (c *Client) Path(parts ...string) *url.URL {
	s := append([]string{c.BaseURL.Path}, parts...)
	p := path.Join(s...) + ".json"
	u := &url.URL{Path: p}

	return c.BaseURL.ResolveReference(u)
}

func (c *Client) validResponse(resp *http.Response) bool {
	if c := resp.StatusCode; 200 <= c && c <= 299 {
		return true
	}

	return false
}
