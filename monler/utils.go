package monler

import (
	"net/url"
)

// URLNormalize removes credential, query string and hash from the URL.
func URLNormalize(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}
