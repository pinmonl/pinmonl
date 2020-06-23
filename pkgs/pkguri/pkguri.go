package pkguri

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	DefaultProto = "https"
)

// Errors
var (
	ErrNoURI         = errors.New("pkguri: uri is not defined")
	ErrNoMarshaler   = errors.New("pkguri: cannot be marshal")
	ErrNoUnmarshaler = errors.New("pkguri: cannot be unmarshal")
	ErrHost          = errors.New("pkguri: host not match")
	ErrPath          = errors.New("pkguri: path not match")
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

	// Docker.
	DockerProvider = "docker"
	DockerHost     = "hub.docker.com"
)

// PkgURI contains
type PkgURI struct {
	Provider string
	Host     string
	URI      string
	Proto    string
}

func NewFromURI(pkguri string) (*PkgURI, error) {
	return unmarshal(pkguri)
}

// Namespace returns the namespace of repo.
func (p *PkgURI) Namespace() string {
	splits := p.URISplits()
	slen := len(splits)
	if slen == 0 {
		return ""
	}
	return strings.Join(splits[:slen-1], "/")
}

// RepoName returns the name of repo.
func (p *PkgURI) RepoName() string {
	splits := p.URISplits()
	return splits[len(splits)-1]
}

// URISplits splits uri by slash.
func (p *PkgURI) URISplits() []string {
	uri := strings.Trim(p.URI, "/")
	if uri == "" {
		return nil
	}
	return strings.Split(uri, "/")
}

// String produces url format of pkguri,
// which is supported by Parse function.
func (p *PkgURI) String() string {
	pkguri, err := p.MarshalText()
	if err != nil {
		return ""
	}
	return string(pkguri)
}

// URL produces accessible url with the specified protocol.
func (p *PkgURI) URL() string {
	return ToURL(p)
}

func (p *PkgURI) UnmarshalText(text []byte) error {
	pu, err := unmarshal(string(text))
	if err != nil {
		return err
	}
	*p = *pu
	return nil
}

func (p PkgURI) MarshalText() ([]byte, error) {
	pkguri, err := marshal(&p)
	if err != nil {
		return nil, err
	}
	return []byte(pkguri), nil
}

type Unmarshaler interface {
	UnmarshalPkgURI(*PkgURI) error
}

func Unmarshal(pkguri string, v interface{}) error {
	switch v.(type) {
	case Unmarshaler:
		pu, err := unmarshal(pkguri)
		if err != nil {
			return err
		}
		return v.(Unmarshaler).UnmarshalPkgURI(pu)
	case encoding.TextUnmarshaler:
		return v.(encoding.TextUnmarshaler).UnmarshalText([]byte(pkguri))
	default:
		return ErrNoUnmarshaler
	}
}

func unmarshal(pkguri string) (*PkgURI, error) {
	u, err := url.Parse(pkguri)
	if err != nil {
		return nil, err
	}
	path := sanitizePath(u.Path)
	if path == "" {
		return nil, ErrNoURI
	}
	proto := u.Query().Get("proto")
	if proto == "" {
		proto = DefaultProto
	}

	pu := &PkgURI{
		Provider: u.Scheme,
		Host:     u.Host,
		URI:      path,
		Proto:    proto,
	}

	switch pu.Provider {
	case GitProvider:
		pu.URI, pu.Host = pu.Host+"/"+pu.URI, ""
	case GithubProvider:
		if pu.Host == GithubHost {
			pu.Host = ""
		}
	}

	return pu, nil
}

type Marshaler interface {
	MarshalPkgURI() (*PkgURI, error)
}

func Marshal(v interface{}) (string, error) {
	switch v.(type) {
	case Marshaler:
		pu, err := v.(Marshaler).MarshalPkgURI()
		if err != nil {
			return "", err
		}
		return pu.String(), nil
	case encoding.TextMarshaler:
		text, err := v.(encoding.TextMarshaler).MarshalText()
		return string(text), err
	default:
		return "", ErrNoMarshaler
	}
}

func marshal(pu *PkgURI) (string, error) {
	u := &url.URL{
		Scheme: pu.Provider,
		Host:   pu.Host,
		Path:   "/" + pu.URI,
	}
	query := url.Values{}
	if pu.Proto != "" && pu.Proto != DefaultProto {
		query.Add("proto", pu.Proto)
	}
	u.RawQuery = query.Encode()

	switch pu.Provider {
	case GithubProvider:
		if pu.Host == GithubHost {
			u.Host = ""
		}
	}

	return u.String(), nil
}

// ParseGit parses git url to pkguri.
func ParseGit(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	proto := u.Scheme
	if proto == "" {
		proto = DefaultProto
	}
	return &PkgURI{
		Provider: GitProvider,
		URI:      u.Host + u.Path,
		Proto:    proto,
	}, nil
}

// ParseGithub parses github url to pkguri.
func ParseGithub(rawurl string) (*PkgURI, error) {
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

	proto := u.Scheme
	if proto == "" {
		proto = DefaultProto
	}

	return &PkgURI{
		Provider: GithubProvider,
		URI:      uri,
		Proto:    proto,
	}, nil
}

// ParseNpm parses npm package url to pkguri.
func ParseNpm(rawurl string) (*PkgURI, error) {
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

	proto := u.Scheme
	if proto == "" {
		proto = DefaultProto
	}

	return &PkgURI{
		Provider: NpmProvider,
		URI:      uri,
		Proto:    proto,
	}, nil
}

func ToURL(pu *PkgURI) string {
	u := &url.URL{
		Scheme: pu.Proto,
		Host:   pu.Host,
		Path:   pu.URI,
	}

	if u.Scheme == "" {
		u.Scheme = DefaultProto
	}

	switch pu.Provider {
	case GithubProvider:
		if u.Host == "" {
			u.Host = GithubHost
		}
	case GitProvider:
		splits := strings.Split(pu.URI, "/")
		if len(splits) > 0 {
			u.Host = splits[0]
		}
		if len(splits) > 1 {
			u.Path = strings.Join(splits[1:], "/")
		}
	case NpmProvider:
		if u.Host == "" {
			u.Host = NpmHost
		}
		u.Path = fmt.Sprintf("/%s/%s", NpmPrefix, strings.Trim(u.Path, "/"))
	}

	return u.String()
}
