package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pinmonl/pinmonl/monler"
)

// Name is static label of provider.
var Name = "git"

// ProviderOpts defines the options of creating provider.
type ProviderOpts struct{}

// Provider handles Git repo.
type Provider struct{}

var _ monler.Provider = &Provider{}

// NewProvider creates provider.
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
	return NewReport(&ReportOpts{
		URL: rawurl,
	})
}

// ParseURL extracts uri from url.
func (p *Provider) ParseURL(rawurl string) (*monler.URL, error) {
	return ParseURL(rawurl)
}

// Ping reports error if the url is not supported.
func Ping(rawurl string, cred monler.Credential) error {
	_, err := LsRemote(rawurl, cred)
	if ErrIsEmptyGitRepository(err) {
		return nil
	}
	return err
}

// ParseURL extracts uri from url.
func ParseURL(rawurl string) (*monler.URL, error) {
	return monler.NewURLFromRaw(rawurl, rawurl)
}

// LsRemote is equivalent call to "git ls-remote".
func LsRemote(rawurl string, cred monler.Credential) ([]*plumbing.Reference, error) {
	r := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{rawurl},
	})
	return r.List(&git.ListOptions{})
}

// ErrIsEmptyGitRepository reports the error is empty git repository.
func ErrIsEmptyGitRepository(err error) bool {
	return err == transport.ErrEmptyRemoteRepository
}
