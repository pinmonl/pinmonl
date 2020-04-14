package github

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/pinmonl/pinmonl/monler"
	"golang.org/x/oauth2"
)

// Name is static label of provider.
var Name = "github"

// Github errors.
var (
	ErrMissingToken = errors.New("github token must be provided")
)

// Github endpoints.
var (
	DefaultEndpoint    = "https://github.com"
	DefaultAPIEndpoint = "https://api.github.com"
)

// ProviderOpts defines the options of Github provider initiation.
type ProviderOpts struct {
	Token string
}

// Provider handles Github url.
type Provider struct {
	token string
}

var _ monler.Provider = &Provider{}

// NewProvider creates Github provider.
func NewProvider(opts *ProviderOpts) (*Provider, error) {
	if opts == nil {
		return nil, ErrMissingToken
	}
	return &Provider{
		token: opts.Token,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string { return Name }

// Ping reports error if the url is not supported.
func (p *Provider) Ping(rawurl string, cred monler.Credential) error {
	return Ping(rawurl, cred)
}

// Open returns Github report.
func (p *Provider) Open(rawurl string, cred monler.Credential) (monler.Report, error) {
	u, err := ParseURL(rawurl)
	if err != nil {
		return nil, err
	}
	return NewReport(&ReportOpts{
		URI:    u.URI,
		Client: newOauthClient(p.token),
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

// ExtractURI gets owner and repository name from uri.
func ExtractURI(uri string) (string, string) {
	splits := strings.SplitN(strings.Trim(uri, "/"), "/", 3)
	if len(splits) < 2 {
		return "", ""
	}
	return splits[0], splits[1]
}

// ParseURL extracts uri from url.
func ParseURL(rawurl string) (*monler.URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	switch {
	case u.Hostname() == "github.com":
	default:
		return nil, monler.ErrNotSupport
	}
	owner, repo := ExtractURI(u.Path)
	if owner == "" || repo == "" {
		return nil, monler.ErrNotSupport
	}
	uri := owner + "/" + repo
	u.Path = "/" + uri
	return &monler.URL{
		URL: u,
		URI: uri,
	}, nil
}

func newOauthClient(token string) *http.Client {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}
