package pkguri

import (
	"strings"
)

func sanitizePath(path string) string {
	return strings.TrimSpace(strings.Trim(path, "/"))
}
