package monl

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pinmonl/pinmonl/monlvars"
)

type URL struct {
	MonlerName string
	ObjectName string
	Prefix     string
	Name       string
	Protocol   string
}

func ParseURL(rawurl string) (*URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimLeft(u.Path, "/")
	switch u.Host {
	case monlvars.GithubHost:
		return parseGithubURL(u)
	default:
		return nil, ErrNoMatch
	}
}

var (
	ReGithubRepo   = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
	ReGithubAuthor = regexp.MustCompile(`^([^/]+)$`)
)

func parseGithubURL(u *url.URL) (*URL, error) {
	mu := &URL{MonlerName: monlvars.Github}
	if ReGithubRepo.MatchString(u.Path) {
		matches := ReGithubRepo.FindStringSubmatch(u.Path)
		owner, repo := matches[1], matches[2]
		mu.ObjectName = monlvars.GithubRepo
		mu.Prefix = owner
		mu.Name = repo
		return mu, nil
	}
	if ReGithubAuthor.MatchString(u.Path) {
		matches := ReGithubAuthor.FindStringSubmatch(u.Path)
		user := matches[1]
		mu.ObjectName = monlvars.GithubUser
		mu.Name = user
		return mu, nil
	}
	return nil, ErrNoMatch
}
