package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

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
	found, err := p.Monls.List(ctx, &store.MonlOpts{URL: pinl.URL})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		monl = found[0]
	} else {
		monl = &model.Monl{URL: pinl.URL}
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
