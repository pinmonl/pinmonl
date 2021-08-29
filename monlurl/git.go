package monlurl

import (
	"fmt"
	"net/url"

	"github.com/pinmonl/pinmonl/monlvars"
)

func ParseGitUrl(u *url.URL) (*Url, error) {
	return &Url{
		MonlerName: monlvars.Git,
		ObjectName: monlvars.ObjectGitRepo,
		Name:       u.Path,
		Host:       u.Host,
		Scheme:     u.Scheme,
	}, nil
}

func BuildGitUrl(u *Url) string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Name)
}
