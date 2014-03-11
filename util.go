package confreaks

import (
	"code.google.com/p/go.net/html"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURI = "http://confreaks.com/"

func relativePath(pathStr string) *url.URL {
	uri, _ := url.Parse(baseURI)
	uri.Path = pathStr

	return uri
}

func fetch(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func jsonMarshal(n interface{}) ([]byte, error) {
	return json.MarshalIndent(n, "", "  ")
}

func attrVal(n *html.Node, k string) string {
	for _, a := range n.Attr {
		if a.Key == k {
			return a.Val
		}
	}

	return ""
}
