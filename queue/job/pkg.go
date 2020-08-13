package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// PkgCrawler defines the job of pkg self update
// independent from monl.
type PkgCrawler struct {
	PkgID  string
	report provider.Report
}

func NewPkgCrawler(pkgID string) *PkgCrawler {
	return &PkgCrawler{
		PkgID: pkgID,
	}
}

func (p *PkgCrawler) String() string {
	return "pkg_self_update"
}

func (p *PkgCrawler) Describe() []string {
	return []string{
		p.String(),
		p.PkgID,
	}
}

func (p *PkgCrawler) Target() model.Morphable {
	return model.Pkg{ID: p.PkgID}
}

func (p *PkgCrawler) RunAt() time.Time {
	return time.Time{}
}

func (p *PkgCrawler) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)
	pkg, err := stores.Pkgs.Find(ctx, p.PkgID)
	if err != nil {
		return err
	}

	pu, err := pkg.MarshalPkgURI()
	if err != nil {
		return err
	}

	repo, err := monler.Parse(pu.String())
	if err != nil {
		return err
	}

	report, err := repo.Analyze()
	if err != nil {
		return err
	}
	p.report = report
	return nil
}

func (p *PkgCrawler) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)
	defer p.report.Close()
	_, _, err := storeutils.SaveProviderReport(ctx, stores.Pkgs, stores.Stats, p.report, true)
	return nil, err
}

var _ Job = &PkgCrawler{}
