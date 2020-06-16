package pinlutils

import "net/url"

func IsValidURL(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}
