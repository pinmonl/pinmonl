package helm

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pinmonl/pinmonl/monler"
)

// Name is label of provider.
var Name = "helm"

// Hub endpoints.
var (
	DefaultEndpoint    = "https://hub.helm.sh"
	DefaultAPIEndpoint = "https://hub.helm.sh/api/chartsvc/v1"
	HubSourceFile      = "https://raw.githubusercontent.com/helm/hub/master/config/repo-values.yaml"
)

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles helm url.
type Provider struct {
	client *http.Client
}

var _ monler.Provider = &Provider{}

// NewProvider creates provider.
func NewProvider(opts *ProviderOpts) (*Provider, error) {
	return &Provider{
		client: &http.Client{},
	}, nil
}

// Name returns the unique name of provider.
func (p *Provider) Name() string { return Name }

// Ping reports error if the url is not supported.
func (p *Provider) Ping(rawurl string, cred monler.Credential) error {
	return Ping(rawurl, cred)
}

var uriRegex = regexp.MustCompile("(?i)^([a-z0-9-]+)/([a-z0-9-]+)")

// Open creates report.
func (p *Provider) Open(rawurl string, cred monler.Credential) (monler.Report, error) {
	u, err := ParseURL(rawurl)
	if err != nil {
		if uriRegex.MatchString(rawurl) {
			return NewReport(&ReportOpts{
				URI:    rawurl,
				Client: &http.Client{},
			})
		}
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
	c, err := newAPIClient(cred)
	if err != nil {
		return err
	}
	_, err = c.getChart(u.URI)
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
	switch {
	case u.Hostname() == "hub.helm.sh":
	default:
		return nil, monler.ErrNotSupport
	}
	sm := regexp.MustCompile(`/?charts/(.*)`).FindStringSubmatch(u.Path)
	uri := sm[len(sm)-1]
	return &monler.URL{
		URL: u,
		URI: uri,
	}, nil
}

// ExtractURI extracts the name of repo and chart from uri.
func ExtractURI(uri string) (string, string) {
	splits := strings.SplitN(uri, "/", 3)
	if len(splits) < 2 {
		return "", ""
	}
	return splits[0], splits[1]
}

// Search finds related helm charts.
func Search(query string) ([]*ChartResponse, error) {
	api, err := newAPIClient(nil)
	if err != nil {
		return nil, err
	}
	return api.searchChart("", query)
}
