package npm

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/pinmonl/pinmonl/monler"
)

// Name is the label of provider.
var Name = "npm"

// NPM endpoints.
var (
	DefaultEndpoint         = "https://www.npmjs.com"
	DefaultRegistryEndpoint = "https://registry.npmjs.org"
	DefaultAPIEndpoint      = "https://api.npmjs.org"
)

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles NPM url.
type Provider struct{}

var _ monler.Provider = &Provider{}

// NewProvider create NPM provider.
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
	c := &apiClient{client: &http.Client{}}
	_, err = c.getPackage(u.URI)
	if err != nil {
		return err
	}
	return nil
}

var packageRe = regexp.MustCompile(`^https?://(www\.)?npmjs\.com/package/(.+)$`)

// ParseURL extracts uri from url.
func ParseURL(rawurl string) (*monler.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	switch {
	case packageRe.MatchString(rawurl):
	default:
		return nil, monler.ErrNotSupport
	}
	sm := regexp.MustCompile(`/?package/(.+)`).FindStringSubmatch(u.Path)
	if len(sm) < 2 {
		return nil, monler.ErrNotSupport
	}
	uri := sm[len(sm)-1]
	return &monler.URL{
		URL: u,
		URI: uri,
	}, nil
}
