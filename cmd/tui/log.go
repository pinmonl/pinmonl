package tui

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var logw io.Writer

func initLogger(dev bool) {
	if !dev {
		logw = ioutil.Discard
	} else {
		f, err := os.Create("tui.log")
		if err != nil {
			panic(err)
		}
		logw = f
	}
}

func debugln(arg ...interface{}) {
	str := fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), fmt.Sprint(arg...))
	fmt.Fprint(logw, str)
}
