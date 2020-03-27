package tui

import (
	"io/ioutil"

	"github.com/pkg/browser"
)

func init() {
	browser.Stderr = ioutil.Discard
	browser.Stdout = ioutil.Discard
}

func OpenURL(url string) error {
	debugln("open url: ", url)
	return browser.OpenURL(url)
}
