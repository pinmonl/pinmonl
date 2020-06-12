package pkguri

import (
	"errors"
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
)

// PkgURI contains
type PkgURI struct {
	Provider string
	Host     string
	URI      string
	Proto    string
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
	proto := DefaultProto
	if protoq := u.Query().Get("proto"); protoq != "" {
		proto = protoq
	}
	return &PkgURI{
		Provider: u.Scheme,
		Host:     u.Host,
		URI:      path,
		Proto:    proto,
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
		Path:   "/" + p.URI,
	}
	query := url.Values{}
	if p.Proto != "" && p.Proto != DefaultProto {
		query.Add("proto", p.Proto)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// URL produces accessible url with the specified protocol.
func (p *PkgURI) URL() *url.URL {
	u := &url.URL{
		Scheme: p.Proto,
		Host:   p.Host,
		Path:   p.URI,
	}
	if u.Scheme == "" {
		u.Scheme = DefaultProto
	}
	return u
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
