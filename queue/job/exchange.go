package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pinmonl-go"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/pubsub/message"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

type FetchMonl struct {
	MonlID string

	monl  *model.Monl
	pkgs  map[*model.Pkg]model.MonpkgKind
	stats map[*model.Pkg][]*model.Stat
}

func NewFetchMonl(monlID string) *FetchMonl {
	return &FetchMonl{
		MonlID: monlID,
	}
}

func (f *FetchMonl) String() string {
	return "fetch_monl"
}

func (f *FetchMonl) Describe() []string {
	return []string{
		f.String(),
		f.MonlID,
	}
}

func (f *FetchMonl) Target() model.Morphable {
	return model.Monl{ID: f.MonlID}
}

func (f *FetchMonl) RunAt() time.Time {
	return time.Time{}
}

func (f *FetchMonl) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)

	monl, err := stores.Monls.Find(ctx, f.MonlID)
	if err != nil {
		return err
	}
	f.monl = monl

	exm := ExchangeManagerFrom(ctx)
	if exm == nil {
		return ErrNoExchangeManager
	}
	client := exm.UserClient()

	pOpts := &pinmonl.PkgListOpts{}
	pOpts.URL = monl.URL
	pOpts.Size = -1
	pResp, err := client.PkgList(pOpts)
	if err != nil {
		return err
	}

	f.pkgs = make(map[*model.Pkg]model.MonpkgKind)
	f.stats = make(map[*model.Pkg][]*model.Stat)
	for i := range pResp.Data {
		mpsrc := pResp.Data[i]
		pu, err := f.parseURI(mpsrc)
		if err != nil {
			return err
		}

		pkg, err := stores.Pkgs.FindURI(ctx, pu)
		if err != nil {
			return err
		}
		if pkg == nil {
			pkg = &model.Pkg{}
			if err := pkg.UnmarshalPkgURI(pu); err != nil {
				return err
			}
		}
		pkg.Title = mpsrc.Pkg.Title
		pkg.Description = mpsrc.Pkg.Description
		pkg.CustomUri = mpsrc.Pkg.CustomUri
		f.pkgs[pkg] = model.MonpkgKind(mpsrc.Kind)

		sOpts := &pinmonl.StatListOpts{}
		sOpts.Latest = field.NewNullBool(true)
		sOpts.Pkgs = []string{mpsrc.PkgID}
		sOpts.Size = -1
		sResp, err := client.StatList(sOpts)
		if err != nil {
			return err
		}

		stats := make([]*model.Stat, len(sResp.Data))
		for j := range sResp.Data {
			ssrc := sResp.Data[j]
			stat, err := f.parseStat(ssrc)
			if err != nil {
				return err
			}
			stats[j] = stat
		}
		f.stats[pkg] = stats
	}

	return nil
}

func (f *FetchMonl) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)

	hub := PubsuberFrom(ctx)
	if hub == nil {
		return nil, ErrNoPubsuber
	}

	monl := f.monl
	monl.FetchedAt = field.Now()
	if err := stores.Monls.Update(ctx, monl); err != nil {
		return nil, err
	}

	for pkg, kind := range f.pkgs {
		var err error
		pkg.FetchedAt = field.Now()
		if pkg.ID == "" {
			err = stores.Pkgs.Create(ctx, pkg)
		} else {
			err = stores.Pkgs.Update(ctx, pkg)
		}
		if err != nil {
			return nil, err
		}

		monpkg, err := stores.Monpkgs.FindOrCreate(ctx, &model.Monpkg{
			MonlID: monl.ID,
			PkgID:  pkg.ID,
		})
		if err != nil {
			return nil, err
		}
		monpkg.Kind = kind
		if err := stores.Monpkgs.Update(ctx, monpkg); err != nil {
			return nil, err
		}

		prevStats, err := stores.Stats.List(ctx, &store.StatOpts{
			PkgIDs: []string{pkg.ID},
		})
		if err != nil {
			return nil, err
		}
		for _, prevStat := range prevStats {
			if _, err := stores.Stats.Delete(ctx, prevStat.ID); err != nil {
				return nil, err
			}
		}

		for i := range f.stats[pkg] {
			stat := f.stats[pkg][i]
			stat.PkgID = pkg.ID
			if err := stores.Stats.Create(ctx, stat); err != nil {
				return nil, err
			}
		}
	}

	pinls, err := storeutils.ListPinls(ctx, stores.Pinls, stores.Monpkgs, stores.Pinpkgs, stores.Taggables, &store.PinlOpts{
		MonlIDs: []string{monl.ID},
	})
	if err != nil {
		return nil, err
	}
	for i := range pinls {
		hub.Broadcast(message.NewPinlUpdated(pinls[i]))
	}

	return nil, nil
}

func (f *FetchMonl) parseURI(src *pinmonl.Monpkg) (*pkguri.PkgURI, error) {
	psrc := src.Pkg
	return &pkguri.PkgURI{
		Provider: psrc.Provider,
		Host:     psrc.ProviderHost,
		URI:      psrc.ProviderURI,
		Proto:    psrc.ProviderProto,
	}, nil
}

func (f *FetchMonl) parseStat(src *pinmonl.Stat) (*model.Stat, error) {
	return &model.Stat{
		RecordedAt:  src.RecordedAt,
		Kind:        model.StatKind(src.Kind),
		Value:       src.Value,
		ValueType:   model.StatValueType(src.ValueType),
		Checksum:    src.Checksum,
		Weight:      src.Weight,
		IsLatest:    src.IsLatest,
		HasChildren: src.HasChildren,
	}, nil
}

var _ Job = &FetchMonl{}
