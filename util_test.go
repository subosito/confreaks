package confreaks

import (
	"fmt"
	"io/ioutil"
)

func loadFixture(name string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("fixtures/%s", name))
}
