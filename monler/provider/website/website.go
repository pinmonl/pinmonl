package website

import (
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/prvdutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

type Provider struct{}

func NewProvider() (*Provider, error) {
	return &Provider{}, nil
}

func (p *Provider) ProviderName() string {
	return pkgdata.WebsiteProvider
}

func (p *Provider) Open(rawurl string) (provider.Repo, error) {
	pu, err := pkguri.ParseWebsite(rawurl)
	if err != nil {
		return nil, err
	}
	return newRepo(pu)
}

func (p *Provider) Parse(uri string) (provider.Repo, error) {
	pu, err := pkguri.NewFromURI(uri)
	if err != nil {
		return nil, err
	}
	return newRepo(pu)
}

func (p *Provider) Ping(rawurl string) error {
	return provider.ErrNoPing
}

type Repo struct {
	*prvdutils.StaticReport
	pu *pkguri.PkgURI
}

func newRepo(pu *pkguri.PkgURI) (*Repo, error) {
	report := prvdutils.NewStaticReport(pu, nil, nil)
	return &Repo{
		StaticReport: report,
		pu:           pu,
	}, nil
}

func (r *Repo) Analyze() (provider.Report, error) {
	return r, nil
}

func (r *Repo) Derived() ([]string, error) {
	return nil, nil
}

func (r *Repo) Close() error {
	return r.StaticReport.Close()
}

var _ provider.Provider = &Provider{}
var _ provider.Repo = &Repo{}
var _ provider.Report = &Repo{}
