package confreaks

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

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

func downloadVideo(u, dir string) error {
	var stderr bytes.Buffer

	cmd := exec.Command("/usr/bin/youtube-dl", "-o", fmt.Sprintf("%s/%%(title)s-%%(id)s.%%(ext)s", dir), fmt.Sprintf("%s", u))
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
