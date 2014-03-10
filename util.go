package confreaks

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
)

const baseURI = "http://confreaks.com/"

func relativePath(pathStr string) (uri *url.URL, err error) {
	uri, err = url.Parse(baseURI)
	if err == nil {
		uri.Path = pathStr
	}

	return
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

