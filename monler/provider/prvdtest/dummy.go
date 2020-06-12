package prvdtest

import (
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

type DummyProvider string

func (p DummyProvider) ProviderName() string {
	return string(p)
}

func (p DummyProvider) Open(rawurl string) (provider.Repo, error) {
	logrus.Debugf("dummyprovider: %q open with %q", p, rawurl)
	return p.open(rawurl, "")
}

func (p DummyProvider) Parse(rawuri string) (provider.Repo, error) {
	logrus.Debugf("dummyprovider: %q parse with %q", p, rawuri)
	return p.open("", rawuri)
}

func (p DummyProvider) open(rawurl, rawuri string) (DummyRepo, error) {
	return DummyRepo{
		Provider: string(p),
		RawURL:   rawurl,
		RawURI:   rawuri,
	}, nil
}

func (p DummyProvider) Ping(rawurl string) error {
	logrus.Debugf("dummyprovider: %q ping with %q", p, rawurl)
	return nil
}

type DummyRepo struct {
	Provider string
	RawURL   string
	RawURI   string
}

func (r DummyRepo) Analyze() (provider.Report, error) {
	return r.analyze()
}

func (r DummyRepo) Derived() ([]provider.Report, error) {
	return nil, nil
}

func (r DummyRepo) analyze() (DummyRepo, error) {
	return r, nil
}

func (r DummyRepo) Close() error {
	return nil
}

func (r DummyRepo) String() string {
	pu, _ := r.URI()
	return pu.String()
}

func (r DummyRepo) URI() (*pkguri.PkgURI, error) {
	pu := &pkguri.PkgURI{
		Provider: r.Provider,
		URI:      r.RawURI,
	}
	if pu.URI == "" {
		pu.URI = r.RawURL
	}
	return pu, nil
}

func (r DummyRepo) Stats() ([]*model.Stat, error) {
	return nil, nil
}

func (r DummyRepo) Next() bool {
	return false
}

func (r DummyRepo) Tag() (*model.Stat, error) {
	return nil, nil
}

var _ provider.Provider = DummyProvider("")
var _ provider.Repo = DummyRepo{}
var _ provider.Report = DummyRepo{}
