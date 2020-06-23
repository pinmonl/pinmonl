package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// PkgSelfUpdate defines the job of pkg self update
// independent from monl.
type PkgSelfUpdate struct {
	PkgID  string
	report provider.Report
}

func NewPkgSelfUpdate(pkgID string) *PkgSelfUpdate {
	return &PkgSelfUpdate{
		PkgID: pkgID,
	}
}

func (p *PkgSelfUpdate) String() string {
	return "pkg_self_update"
}

func (p *PkgSelfUpdate) Describe() []string {
	return []string{
		p.String(),
		p.PkgID,
	}
}

func (p *PkgSelfUpdate) Target() model.Morphable {
	return model.Pkg{ID: p.PkgID}
}

func (p *PkgSelfUpdate) RunAt() time.Time {
	return time.Time{}
}

func (p *PkgSelfUpdate) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)
	pkg, err := stores.Pkgs.Find(ctx, p.PkgID)
	if err != nil {
		return err
	}

	pu, err := pkg.MarshalPkgURI()
	if err != nil {
		return err
	}

	err = monler.Ping(pu.Provider, pkguri.ToURL(pu))
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

func (p *PkgSelfUpdate) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)
	defer p.report.Close()
	_, _, err := storeutils.SaveProviderReport(ctx, stores.Pkgs, stores.Stats, p.report)
	return nil, err
}

var _ Job = &PkgSelfUpdate{}
