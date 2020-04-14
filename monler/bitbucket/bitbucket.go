package bitbucket

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/monler"
)

// Name is label of provider.
var Name = "bitbucket"

// BitBucket endpoints.
var (
	DefaultEndpoint    = "https://bitbucket.org"
	DefaultAPIEndpoint = "https://api.bitbucket.org"
)

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles BitBucket url.
type Provider struct{}

var _ monler.Provider = &Provider{}

// NewProvider creates BitBucket provider.
func NewProvider(opts *ProviderOpts) (*Provider, error) {
	return &Provider{}, nil
}

// Name returns the unique name of provider.
func (p *Provider) Name() string { return Name }

// Ping reports error if the url is not supported.
func (p *Provider) Ping(rawurl string, cred monler.Credential) error {
	return Ping(rawurl, cred)
}

// Open creates report.
func (p *Provider) Open(rawurl string, cred monler.Credential) (monler.Report, error) {
	u, err := ParseURL(rawurl)
	if err != nil {
		return nil, err
	}
	return NewReport(&ReportOpts{
		URI:    u.URI,
		Client: &http.Client{},
	})
}

// ParseURL extracts uri from url.
func (p *Provider) ParseURL(rawurl string) (*monler.URL, error) {
	return ParseURL(rawurl)
}

// Ping reports error if the url is not supported.
func Ping(rawurl string, cred monler.Credential) error {
	u, err := ParseURL(rawurl)
	if err != nil {
		return err
	}
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return monler.ErrNotExist
	}
	return nil
}

// ParseURL extracts uri from url.
func ParseURL(rawurl string) (*monler.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	switch {
	case u.Hostname() == "bitbucket.com":
	default:
		return nil, monler.ErrNotSupport
	}
	ws, repo := ExtractURI(u.Path)
	if ws == "" || repo == "" {
		return nil, monler.ErrNotSupport
	}
	uri := ws + "/" + repo
	u.Path = "/" + uri
	return &monler.URL{
		URL: u,
		URI: uri,
	}, nil
}

// ExtractURI extracts uri into parts.
func ExtractURI(uri string) (string, string) {
	splits := strings.SplitN(strings.Trim(uri, "/"), "/", 3)
	if len(splits) < 2 {
		return "", ""
	}
	return splits[0], splits[1]
}
