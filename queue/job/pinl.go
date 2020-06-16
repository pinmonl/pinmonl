package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/store"
)

// PinlUpdated defines the job whenever a pinl is created or updated.
//
// It finds or creates monl by a normalized url. pinl.MonlID is updated
// accordingly and creates a job for the new monl.
type PinlUpdated struct {
	PinlID  string
	Pinls   *store.Pinls
	Monls   *store.Monls
	Pkgs    *store.Pkgs
	Stats   *store.Stats
	Monpkgs *store.Monpkgs
}

func NewPinlUpdated(pinlID string, pinls *store.Pinls, monls *store.Monls, pkgs *store.Pkgs, stats *store.Stats, monpkgs *store.Monpkgs) PinlUpdated {
	return PinlUpdated{
		PinlID:  pinlID,
		Pinls:   pinls,
		Monls:   monls,
		Pkgs:    pkgs,
		Stats:   stats,
		Monpkgs: monpkgs,
	}
}

func (p PinlUpdated) String() string {
	return "pinl_updated"
}

func (p PinlUpdated) Describe() []string {
	return []string{
		p.String(),
		p.PinlID,
	}
}

func (p PinlUpdated) RunAt() time.Time {
	return time.Time{}
}

func (p PinlUpdated) Run(ctx context.Context) ([]Job, error) {
	pinl, err := p.Pinls.Find(ctx, p.PinlID)
	if err != nil {
		return nil, err
	}

	var (
		monl *model.Monl
		jobs []Job
	)
	url, err := monlutils.NormalizeURL(pinl.URL)
	if err != nil {
		return nil, err
	}
	found, err := p.Monls.List(ctx, &store.MonlOpts{URL: url.String()})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		monl = found[0]
	} else {
		monl = &model.Monl{URL: url.String()}
		err := p.Monls.Create(ctx, monl)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, NewMonlCreated(monl.ID, p.Monls, p.Pkgs, p.Stats, p.Monpkgs))
	}

	pinl.MonlID = monl.ID
	err = p.Pinls.Update(ctx, pinl)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var _ Job = PinlUpdated{}
