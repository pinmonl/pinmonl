package pkguri

import (
	"net/url"
	"regexp"
	"strings"
)

type Patterns map[string]*regexp.Regexp

// FindPattern finds the matched regexp pattern and its key.
func (p Patterns) FindPattern(url *url.URL) (string, *regexp.Regexp) {
	for h, re := range p {
		if url.Host != h {
			continue
		}
		path := strings.TrimPrefix(url.Path, "/")
		if re.MatchString(path) {
			return h, re
		}
	}
	return "", nil
}

func (p Patterns) MatchURL(url *url.URL) bool {
	_, re := p.FindPattern(url)
	return re != nil
}

func sanitizePath(path string) string {
	return strings.TrimSpace(strings.Trim(path, "/"))
}
