package provider

import (
	"fmt"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

// Provider is the interface that must be implemented by a monler provider.
type Provider interface {
	// ProviderName is the provider name.
	ProviderName() string

	// Open creates Repo with the web url.
	// Note: It is assumed Ping check is done before calling Open.
	Open(url string) (Repo, error)

	// Parse creates Repo with the pkguri format url.
	Parse(uri string) (Repo, error)

	// Ping reports whether the url is supported by the provider.
	Ping(url string) error
}

// Repo is the interface of a repo.
type Repo interface {
	// Analyze creates Report.
	Analyze() (Report, error)

	// Derived returns urls which are related.
	Derived() ([]string, error)

	// Close closes and frees up resources.
	Close() error
}

// Report contains the analyzed data from repo.
type Report interface {
	// Stringer returns string in pkguri format.
	fmt.Stringer

	// URI returns PkgURI.
	URI() (*pkguri.PkgURI, error)

	// Stats contains list of stats other than tag.
	Stats() ([]*model.Stat, error)

	// Next reports there has next tag.
	Next() bool

	// Tag returns the tag at the cursor.
	Tag() (*model.Stat, error)

	// Close closes and frees up resources.
	Close() error
}
