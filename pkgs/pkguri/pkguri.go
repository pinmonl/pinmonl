package pkguri

import (
	"errors"
	"net/url"
	"strings"
)

// Errors
var (
	ErrNoURI         = errors.New("pkguri: uri is not defined")
	ErrNoMarshaler   = errors.New("pkguri: cannot be marshal")
	ErrNoUnmarshaler = errors.New("pkguri: cannot be unmarshal")
)

// PkgURI contains
type PkgURI struct {
	Provider string
	Host     string
	URI      string
}

// Parse fills PkgURI by the data from rawurl.
func Parse(rawurl string) (*PkgURI, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	path := strings.Trim(u.Path, "/")
	if path == "" {
		return nil, ErrNoURI
	}
	return &PkgURI{
		Provider: u.Scheme,
		Host:     u.Host,
		URI:      path,
	}, nil
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
	u := &url.URL{
		Scheme: p.Provider,
		Host:   p.Host,
		Path:   p.URI,
	}
	return u.String()
}

type Marshaler interface {
	MarshalPkgURI() (*PkgURI, error)
}

func Marshal(v interface{}) (string, error) {
	var (
		pu  *PkgURI
		err error
	)
	switch v.(type) {
	case Marshaler:
		pu, err = v.(Marshaler).MarshalPkgURI()
	default:
		err = ErrNoMarshaler
	}

	if err != nil {
		return "", err
	}
	return pu.String(), nil
}

type Unmarshaler interface {
	UnmarshalPkgURI(*PkgURI) error
}

func Unmarshal(data string, v interface{}) error {
	pu, err := Parse(data)
	if err != nil {
		return err
	}
	switch v.(type) {
	case Unmarshaler:
		return v.(Unmarshaler).UnmarshalPkgURI(pu)
	default:
		return ErrNoUnmarshaler
	}
}