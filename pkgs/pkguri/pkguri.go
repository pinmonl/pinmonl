package pkguri

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
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
	case pkgdata.GitProvider:
		pu.URI, pu.Host = pu.Host+"/"+pu.URI, ""
	case pkgdata.GithubProvider:
		if pu.Host == pkgdata.GithubHost {
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
	case pkgdata.GithubProvider:
		if pu.Host == pkgdata.GithubHost {
			u.Host = ""
		}
	}

	return u.String(), nil
}

func getProto(proto string) string {
	if proto == "" {
		return DefaultProto
	}
	return proto
}

// ParseGit parses git url to pkguri.
func ParseGit(rawurl string) (*PkgURI, error) {
	u, err := monlutils.NormalizeURL(rawurl)
	if err != nil {
		return nil, err
	}
	return &PkgURI{
		Provider: pkgdata.GitProvider,
		URI:      u.Host + u.Path,
		Proto:    getProto(u.Scheme),
	}, nil
}

// ParseGithub parses github url to pkguri.
func ParseGithub(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != pkgdata.GithubHost {
		return nil, ErrHost
	}

	path := sanitizePath(u.Path)
	splits := strings.SplitN(path, "/", 3)
	if len(splits) < 2 {
		return nil, ErrNoURI
	}
	uri := strings.Join(splits[:2], "/")

	return &PkgURI{
		Provider: pkgdata.GithubProvider,
		URI:      uri,
		Proto:    getProto(u.Scheme),
	}, nil
}

// ParseNpm parses npm package url to pkguri.
func ParseNpm(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != pkgdata.NpmHost {
		return nil, ErrHost
	}

	path := sanitizePath(u.Path)
	if !strings.HasPrefix(path, pkgdata.NpmPrefix+"/") {
		return nil, ErrPath
	}
	uri := strings.TrimPrefix(path, pkgdata.NpmPrefix+"/")
	if uri == "" {
		return nil, ErrNoURI
	}

	return &PkgURI{
		Provider: pkgdata.NpmProvider,
		URI:      uri,
		Proto:    getProto(u.Scheme),
	}, nil
}

func ParseYoutube(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != pkgdata.YoutubeHost {
		return nil, ErrHost
	}

	re := regexp.MustCompile("^/(channel|c|user)/.+")
	if !re.MatchString(u.Path) {
		return nil, ErrPath
	}
	uri := ""
	splits := strings.SplitN(sanitizePath(u.Path), "/", 3)
	if len(splits) >= 2 {
		uri = splits[1]
	} else {
		return nil, ErrNoURI
	}

	return &PkgURI{
		Provider: pkgdata.YoutubeProvider,
		URI:      uri,
		Proto:    getProto(u.Scheme),
	}, nil
}

func IsYoutubeValidChannelId(id string) bool {
	return regexp.MustCompile("^UC[a-zA-Z0-9-_]{22,22}").MatchString(id)
}

func ParseDocker(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Host != "" && u.Host != pkgdata.DockerHost {
		return nil, ErrHost
	}

	re := regexp.MustCompile("^/(r/([^/]+/[^/]+))|(_/([^/]+))")
	if !re.MatchString(u.Path) {
		return nil, ErrPath
	}

	uri := ""
	if strings.HasPrefix(u.Path, "/r/") {
		path := strings.TrimPrefix(u.Path, "/r/")
		splits := strings.SplitN(path, "/", 3)
		if len(splits) >= 2 {
			uri = strings.Join(splits[0:2], "/")
		} else {
			return nil, ErrNoURI
		}
	} else {
		path := strings.TrimPrefix(u.Path, "/_/")
		if path == "" {
			return nil, ErrNoURI
		}
		uri = fmt.Sprintf("library/%s", path)
	}

	return &PkgURI{
		Provider: pkgdata.DockerProvider,
		URI:      uri,
		Proto:    getProto(u.Scheme),
	}, nil
}

func ParseWebsite(rawurl string) (*PkgURI, error) {
	u, err := monlutils.NormalizeURL(rawurl)
	if err != nil {
		return nil, err
	}
	return &PkgURI{
		Provider: pkgdata.WebsiteProvider,
		URI:      u.Host + u.Path,
		Proto:    getProto(u.Scheme),
	}, nil
}

func IsDockerOfficialRepository(uri string) bool {
	return strings.HasPrefix(uri, "library/")
}

func ToURL(pu *PkgURI) string {
	u := &url.URL{
		Scheme: getProto(pu.Proto),
		Host:   pu.Host,
		Path:   pu.URI,
	}

	switch pu.Provider {
	case pkgdata.GithubProvider:
		if u.Host == "" {
			u.Host = pkgdata.GithubHost
		}

	case pkgdata.GitProvider:
		splits := strings.Split(pu.URI, "/")
		if len(splits) > 0 {
			u.Host = splits[0]
		}
		if len(splits) > 1 {
			u.Path = strings.Join(splits[1:], "/")
		}

	case pkgdata.NpmProvider:
		if u.Host == "" {
			u.Host = pkgdata.NpmHost
		}
		u.Path = fmt.Sprintf("/%s/%s", pkgdata.NpmPrefix, pu.URI)

	case pkgdata.YoutubeProvider:
		if u.Host == "" {
			u.Host = pkgdata.YoutubeHost
		}
		if IsYoutubeValidChannelId(pu.URI) {
			u.Path = fmt.Sprintf("/channel/%s", pu.URI)
		}

	case pkgdata.DockerProvider:
		if u.Host == "" {
			u.Host = pkgdata.DockerHost
		}
		if IsDockerOfficialRepository(pu.URI) {
			u.Path = fmt.Sprintf("/_/%s", strings.TrimPrefix(pu.URI, "library/"))
		} else {
			u.Path = fmt.Sprintf("/r/%s", pu.URI)
		}

	case pkgdata.WebsiteProvider:
		splits := strings.Split(pu.URI, "/")
		if len(splits) > 0 {
			u.Host = splits[0]
		}
		if len(splits) > 1 {
			u.Path = strings.Join(splits[1:], "/")
		}
	}

	return u.String()
}
