package confreaks

import (
	"io/ioutil"
)

func init() {
	log.Out = ioutil.Discard
}
