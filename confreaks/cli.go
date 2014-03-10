package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/subosito/confreaks"
	"io/ioutil"
)

func main() {
	b, err := ioutil.ReadFile("event-details.html")
	if err != nil {
		return
	}

	/* events, _ := confreaks.ParseEvents(bytes.NewReader(b)) */
	/* jb, _ := json.MarshalIndent(events, "", "  ") */
	/* fmt.Printf("%s\n", jb) */

	event, _ := confreaks.ParseEventDetails(bytes.NewReader(b))
	jb, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("%s\n", jb)
}
