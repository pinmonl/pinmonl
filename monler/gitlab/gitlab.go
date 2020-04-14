package gitlab

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/monler"
)

// Name is label of provider.
var Name = "gitlab"

// Gitlab endpoints.
var (
	DefaultEndpoint    = "https://gitlab.com"
	DefaultAPIEndpoint = "https://gitlab.com/api"
)

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles Gitlab url.
type Provider struct{}

var _ monler.Provider = &Provider{}

// NewProvider creates Gitlab provider.
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
	case u.Hostname() == "gitlab.com":
	default:
		return nil, monler.ErrNotSupport
	}
	return &monler.URL{
		URL: u,
		URI: strings.Trim(u.Path, "/"),
	}, nil
}
