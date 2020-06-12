package git

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

var (
	TempDir  = ""
	CloneDir = "monler/git"

	CloneProgress io.Writer
)

func getCloneDir(gitURL string) (string, error) {
	h := sha256.New()
	h.Write([]byte(gitURL))

	tempDir := TempDir
	if tempDir == "" {
		tempDir = os.TempDir()
	}

	dir := fmt.Sprintf("%s/%s/%x", tempDir, CloneDir, h.Sum(nil))
	return dir, nil
}
