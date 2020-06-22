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
	PinlID string
}

func NewPinlUpdated(pinlID string) *PinlUpdated {
	return &PinlUpdated{
		PinlID: pinlID,
	}
}

func (p *PinlUpdated) String() string {
	return "pinl_updated"
}

func (p *PinlUpdated) Describe() []string {
	return []string{
		p.String(),
		p.PinlID,
	}
}

func (p *PinlUpdated) Target() model.Morphable {
	return model.Pinl{ID: p.PinlID}
}

func (p *PinlUpdated) RunAt() time.Time {
	return time.Time{}
}

func (p *PinlUpdated) PreRun(ctx context.Context) error {
	return nil
}

func (p *PinlUpdated) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)
	pinl, err := stores.Pinls.Find(ctx, p.PinlID)
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
	found, err := stores.Monls.List(ctx, &store.MonlOpts{URL: url.String()})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		monl = found[0]
	} else {
		monl = &model.Monl{URL: url.String()}
		err := stores.Monls.Create(ctx, monl)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, NewMonlCreated(monl.ID))
	}

	pinl.MonlID = monl.ID
	err = stores.Pinls.Update(ctx, pinl)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var _ Job = &PinlUpdated{}
