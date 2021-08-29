package monlurl

import (
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/monlvars"
)

type Url struct {
	MonlerName string
	ObjectName string
	Prefix     string
	Name       string
	Scheme     string
	Host       string
}

func Parse(rawurl string) (*Url, error) {
	monler, rawurl := SplitExplicit(rawurl)
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	u = Normalize(u)

	if monler != "" {
		return ParseByExplicit(monler, u)
	} else {
		return ParseByHost(u)
	}
}

const Separator = ","

func SplitExplicit(rawurl string) (string, string) {
	schemeIdx := strings.Index(rawurl, "://")
	sepIdx := strings.Index(rawurl, Separator)
	if sepIdx > -1 && sepIdx < schemeIdx {
		return rawurl[:sepIdx], rawurl[sepIdx+1:]
	} else {
		return "", rawurl
	}
}

func ParseByHost(u *url.URL) (*Url, error) {
	switch u.Host {
	case monlvars.GithubHost:
		return ParseGithubUrl(u)
	default:
		return nil, ErrNoMatch
	}
}

func ParseByExplicit(monler string, u *url.URL) (*Url, error) {
	switch monler {
	case monlvars.Git:
		return ParseGitUrl(u)
	case monlvars.Github:
		return ParseGithubUrl(u)
	default:
		return nil, ErrNoMatch
	}
}

func Normalize(u *url.URL) *url.URL {
	if u.Host == "http" {
		u.Host = "https"
	}
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	u.Path = "/" + strings.Trim(u.Path, "/")
	return u
}

func (u *Url) String() string {
	return u.MonlerName + "," + u.SafeString()
}

func (u *Url) SafeString() string {
	switch u.MonlerName {
	case monlvars.Github:
		return BuildGithubUrl(u)
	default:
		return ""
	}
}
