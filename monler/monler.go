package monler

import (
	"errors"
	"net/url"
)

// General errors.
var (
	ErrNoProvider = errors.New("no supported provider")
	ErrNotSupport = errors.New("url format does not support")
	ErrNotExist   = errors.New("repository does not exist")
)

// Repository manages monler provider.
type Repository struct {
	providers map[string]Provider
}

// NewRepository creates Repository.
func NewRepository() (*Repository, error) {
	return &Repository{
		providers: make(map[string]Provider),
	}, nil
}

// Register adds provider into list.
func (r *Repository) Register(provider Provider) error {
	r.providers[provider.Name()] = provider
	return nil
}

// Providers returns name of registered providers.
func (r *Repository) Providers() []string {
	var names []string
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// Get returns provider by name.
func (r *Repository) Get(providerName string) (Provider, error) {
	if p, ok := r.providers[providerName]; ok {
		return p, nil
	}
	return nil, ErrNoProvider
}

// Open creates report from provider.
func (r *Repository) Open(providerName, rawurl string, cred Credential) (Report, error) {
	p, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}
	return p.Open(rawurl, cred)
}

// Ping reports error if the url is not supported by provider.
func (r *Repository) Ping(providerName, rawurl string, cred Credential) error {
	p, err := r.Get(providerName)
	if err != nil {
		return err
	}
	return p.Ping(rawurl, cred)
}

// Provider defines the provider interface.
type Provider interface {
	// Name returns the unique name of provider.
	Name() string

	// Ping reports error if the url is not supported.
	Ping(rawurl string, cred Credential) error

	// Open creates report.
	Open(rawurl string, cred Credential) (Report, error)

	// ParseURL extracts uri from url.
	ParseURL(rawurl string) (*URL, error)
}

// URL stores uri for provider.
type URL struct {
	*url.URL
	RawURL string
	URI    string
}

// NewURLFromRaw creates URL by plain url.
func NewURLFromRaw(rawurl, uri string) (*URL, error) {
	u := URL{RawURL: rawurl, URI: uri}
	if err := u.ParseRaw(); err != nil {
		return nil, err
	}
	return &u, nil
}

// String parses URL from RawURL if URL is nil.
func (u *URL) String() string {
	if u.URL == nil {
		if err := u.ParseRaw(); err != nil {
			return ""
		}
	}
	return u.URL.String()
}

// ParseRaw parse RawURL to URL.
func (u *URL) ParseRaw() error {
	ur, err := url.Parse(u.RawURL)
	if err != nil {
		return err
	}
	u.URL = ur
	return nil
}

// Credential stores authentication information when open report.
type Credential interface {
	//
}

// Report defines the repo stats report
// where tag stats should be ordered in descending order by the date.
type Report interface {
	// URL returns the url to web page.
	URL() string

	// Provider returns the name of provider.
	Provider() string

	// ProviderURI returns the unique identifier in provider.
	ProviderURI() string

	// Stats returns the stats of repo, such as star, fork, major language.
	Stats() []*Stat

	// Next checks whether has next tag and moves the cursor.
	Next() bool

	// Prev checks whether has previous tag and moves the cursor.
	Prev() bool

	// Tag returns tag at the current cursor.
	Tag() *Stat

	// LatestTag returns the latest tag.
	LatestTag() *Stat

	// Len returns total count of tags.
	Len() int

	// Download downloads stats and tags.
	Download() error

	// Close closes the report.
	Close() error

	// Derived returns related reports.
	Derived(*Repository, Credential) ([]Report, error)
}
