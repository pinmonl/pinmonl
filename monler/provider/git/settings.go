package git

import (
	"io"
	"io/ioutil"
)

var (
	TempDir     = ""
	TempPattern = "monler-git-"

	CloneProgress io.Writer
)

func getTempDir() (string, error) {
	return ioutil.TempDir(TempDir, TempPattern)
}
