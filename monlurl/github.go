package monlurl

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/pinmonl/pinmonl/monlvars"
)

var (
	ReGithubRepo   = regexp.MustCompile(`^/([^/]+)/([^/]+)$`)
	ReGithubAuthor = regexp.MustCompile(`^/([^/]+)$`)
)

func ParseGithubUrl(u *url.URL) (*Url, error) {
	mu := &Url{
		MonlerName: monlvars.Github,
		Host:       monlvars.GithubHost,
	}
	switch {
	case ReGithubRepo.MatchString(u.Path):
		matches := ReGithubRepo.FindStringSubmatch(u.Path)
		owner, repo := matches[1], matches[2]
		mu.ObjectName = monlvars.ObjectGithubRepo
		mu.Prefix = owner
		mu.Name = repo
		return mu, nil
	case ReGithubAuthor.MatchString(u.Path):
		matches := ReGithubAuthor.FindStringSubmatch(u.Path)
		user := matches[1]
		mu.ObjectName = monlvars.ObjectGithubUser
		mu.Name = user
		return mu, nil
	default:
		return nil, ErrNoMatch
	}
}

func BuildGithubUrl(u *Url) string {
	ou := &url.URL{Scheme: "https", Host: monlvars.GithubHost}
	switch {
	case u.ObjectName == monlvars.ObjectGithubRepo:
		ou.Path = fmt.Sprintf("/%s/%s", u.Prefix, u.Name)
	case u.ObjectName == monlvars.ObjectGithubUser:
		ou.Path = fmt.Sprintf("/%s", u.Name)
	}
	return ou.String()
}
