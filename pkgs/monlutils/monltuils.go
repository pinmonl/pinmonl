package monlutils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
)

func NormalizeURL(rawurl string) (*url.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, errors.New("invalid url format")
	}

	if u.Host == pkgdata.GithubHost {
		if strings.HasSuffix(u.Path, ".git") {
			u.Path = strings.TrimSuffix(u.Path, ".git")
		}
	}

	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   strings.TrimSuffix(u.Path, "/"),
	}, nil
}

func IsHttp(rawurl string) bool {
	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
