package pkguri

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Provider url settings.
var (
	// Git.
	GitProvider = "git"

	// Npm.
	NpmProvider     = "npm"
	NpmHost         = "www.npmjs.com"
	NpmPrefix       = "package"
	NpmRegistryHost = "registry.npmjs.org"

	// Github.
	GithubProvider = "github"
	GithubHost     = "github.com"
)

// Errors.
var (
	ErrHost = errors.New("pkguri: host not match")
	ErrPath = errors.New("pkguri: path not match")
)

func ToURL(pu *PkgURI) *url.URL {
	switch pu.Provider {
	case NpmProvider:
		return ToNpm(pu)
	case GithubProvider:
		return ToGithub(pu)
	default:
		return pu.URL()
	}
}

func ParseProvider(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case NpmProvider:
		pu, err := ParseFromNpm(rawurl)
		if err == nil {
			return pu, nil
		}
		return ParseFromNpmRegistry(rawurl)
	case GithubProvider:
		return ParseFromGithub(rawurl)
	default:
		return Parse(rawurl)
	}
}

// ToNpm produces the url to npm package page.
func ToNpm(pu *PkgURI) *url.URL {
	u := pu.URL()
	if u.Host == "" {
		u.Host = NpmHost
	}
	u.Path = fmt.Sprintf("/%s/%s", NpmPrefix, strings.Trim(u.Path, "/"))
	return u
}

// ParseFromNpm parses npm package url to pkguri.
func ParseFromNpm(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != NpmHost {
		return nil, ErrHost
	}
	path := sanitizePath(u.Path)
	if !strings.HasPrefix(path, NpmPrefix+"/") {
		return nil, ErrPath
	}
	uri := strings.TrimPrefix(path, NpmPrefix+"/")
	if uri == "" {
		return nil, ErrNoURI
	}
	return &PkgURI{
		Provider: NpmProvider,
		URI:      uri,
	}, nil
}

// ParseFromNpmRegistry parse npm registry url to pkguri.
func ParseFromNpmRegistry(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != NpmRegistryHost {
		return nil, ErrHost
	}
	path := sanitizePath(u.Path)
	if path == "" {
		return nil, ErrNoURI
	}
	return &PkgURI{
		Provider: NpmProvider,
		URI:      path,
	}, nil
}

// ToGithub produces the url to github repository page.
func ToGithub(pu *PkgURI) *url.URL {
	u := pu.URL()
	if u.Host == "" {
		u.Host = GithubHost
	}
	return u
}

// ParseFromGithub parses github url to pkguri.
func ParseFromGithub(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != GithubHost {
		return nil, ErrHost
	}
	path := sanitizePath(u.Path)
	splits := strings.SplitN(path, "/", 3)
	if len(splits) < 2 {
		return nil, ErrNoURI
	}
	uri := strings.Join(splits[:2], "/")
	return &PkgURI{
		Provider: GithubProvider,
		URI:      uri,
	}, nil
}
