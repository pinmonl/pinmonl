package gitderive

import (
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/monler/provider/git"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

type Derived struct {
	repo    *git.Repo
	reports []provider.Report
}

func New(repo *git.Repo) Derived {
	return Derived{repo: repo}
}

func (d *Derived) Reports() []provider.Report {
	return d.reports
}

func (d *Derived) AddNpm() error {
	pu, err := d.repo.GetNpmURI()
	if err != nil {
		return err
	}

	_, report, err := openFromURI(pu.Provider, pkguri.ToURL(pu))
	if err != nil {
		logrus.Debugf("gitderive: npm err(%v)", err)
		return err
	}

	d.reports = append(d.reports, report)
	return nil
}

func openFromURI(providerName, rawurl string) (provider.Repo, provider.Report, error) {
	err := monler.Ping(providerName, rawurl)
	if err != nil {
		return nil, nil, err
	}

	repo, err := monler.Open(providerName, rawurl)
	if err != nil {
		return nil, nil, err
	}

	report, err := repo.Analyze()
	return repo, report, err
}
