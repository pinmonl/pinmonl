package docker

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/monler"
)

// Name is label of provider.
var Name = "docker"

// Docker endpoints and config.
var (
	DefaultEndpoint     = "https://hub.docker.com"
	DefaultAPIEndpoint  = "https://hub.docker.com"
	OfficialNamespace   = "library"
	OfficialURLPrefix   = "_"
	ThirdPartyURLPrefix = "r"
)

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles Docker url.
type Provider struct{}

var _ monler.Provider = &Provider{}

// NewProvider creates Docker provider.
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
	_, err := ParseURL(rawurl)
	if err != nil {
		return err
	}
	return nil
}

// ParseURL extracts uri from url.
func ParseURL(rawurl string) (*monler.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	var uri string
	switch {
	case u.Hostname() == "hub.docker.com":
		path := strings.Trim(u.Path, "/")
		switch {
		case strings.HasPrefix(path, OfficialURLPrefix+"/"):
			uri = OfficialNamespace + strings.TrimPrefix(path, OfficialURLPrefix)
		case strings.HasPrefix(path, "r/"):
			uri = strings.TrimPrefix(path, "r/")
		default:
			return nil, monler.ErrNotSupport
		}
	default:
		return nil, monler.ErrNotSupport
	}
	return &monler.URL{
		URL: u,
		URI: uri,
	}, nil
}
