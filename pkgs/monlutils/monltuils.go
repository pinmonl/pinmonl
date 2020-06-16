package monlutils

import (
	"errors"
	"net/url"
)

func NormalizeURL(rawurl string) (*url.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid url format")
	}

	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}, nil
}
