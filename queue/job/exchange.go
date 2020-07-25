package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pinmonl-go"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/sirupsen/logrus"
)

type DownloadPinlInfo struct {
	PinlID string
	client *pinmonl.Client

	monl     *model.Monl
	pkgStats map[*model.Pkg][]*model.Stat
}

func NewDownloadPinlInfo(pinlID string, client *pinmonl.Client) *DownloadPinlInfo {
	return &DownloadPinlInfo{PinlID: pinlID, client: client}
}

func (d *DownloadPinlInfo) String() string {
	return "download_pinl_info"
}

func (d *DownloadPinlInfo) Describe() []string {
	return []string{
		d.String(),
		d.PinlID,
	}
}

func (d *DownloadPinlInfo) Target() model.Morphable {
	return model.Pinl{ID: d.PinlID}
}

func (d *DownloadPinlInfo) RunAt() time.Time {
	return time.Time{}
}

func (d *DownloadPinlInfo) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)
	pinl, err := stores.Pinls.Find(ctx, d.PinlID)
	if err != nil {
		return err
	}

	u, err := monlutils.NormalizeURL(pinl.URL)
	if err != nil {
		return err
	}
	monl, err := stores.Monls.FindURL(ctx, u.String())
	if err != nil {
		return err
	}
	if monl != nil {
		logrus.Debugln(u.String(), "existed")
		d.monl = monl
		return nil
	}

	d.monl = &model.Monl{URL: u.String()}
	data := make(map[*model.Pkg][]*model.Stat)
	pData, err := d.client.PkgList(u.String(), nil)
	if err != nil {
		return err
	}
	for _, p := range pData {
		pu := &pkguri.PkgURI{
			Provider: p.Provider,
			Host:     p.ProviderHost,
			URI:      p.ProviderURI,
			Proto:    p.ProviderProto,
		}
		pkg, err := stores.Pkgs.FindURI(ctx, pu)
		if err != nil {
			return err
		}
		if pkg != nil {
			data[pkg] = nil
			continue
		}

		sData, err := d.client.StatLatestList(pu.String(), nil)
		if err != nil {
			return err
		}

		pkg = &model.Pkg{}
		err = pkg.UnmarshalPkgURI(pu)
		if err != nil {
			return err
		}

		for _, s := range sData {
			stat, err := d.parseStat(s)
			if err != nil {
				return err
			}
			data[pkg] = append(data[pkg], stat)
		}
	}
	d.pkgStats = data
	return nil
}

func (d *DownloadPinlInfo) parseStat(src *pinmonl.Stat) (*model.Stat, error) {
	stat := &model.Stat{
		RecordedAt: src.RecordedAt,
		Kind:       src.Kind,
		Value:      src.Value,
		Checksum:   src.Checksum,
		IsLatest:   src.IsLatest,
	}
	for _, subsrc := range src.Substats {
		substat, err := d.parseStat(subsrc)
		if err != nil {
			return nil, err
		}
		*stat.Substats = append(*stat.Substats, substat)
	}
	return stat, nil
}

func (d *DownloadPinlInfo) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)
	var err error

	if d.monl.ID == "" {
		err = stores.Monls.Create(ctx, d.monl)
		if err != nil {
			return nil, err
		}

		for pkg, stats := range d.pkgStats {
			if pkg.ID == "" {
				err = stores.Pkgs.Create(ctx, pkg)
				if err != nil {
					return nil, err
				}
			}

			_, err = stores.Monpkgs.FindOrCreate(ctx, &model.Monpkg{
				MonlID: d.monl.ID,
				PkgID:  pkg.ID,
			})
			if err != nil {
				return nil, err
			}

			for _, stat := range stats {
				stat.PkgID = pkg.ID
				err = stores.Stats.Create(ctx, stat)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Checks if the url of the pinl is still the same.
	pinl, err := stores.Pinls.Find(ctx, d.PinlID)
	if err != nil {
		return nil, err
	}
	u, err := monlutils.NormalizeURL(pinl.URL)
	if err != nil {
		return nil, err
	}
	if u.String() != d.monl.URL {
		return nil, nil
	}

	// Associate to monl if the url is unchanged.
	pinl.MonlID = d.monl.ID
	err = stores.Pinls.Update(ctx, pinl)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

var _ Job = &DownloadPinlInfo{}
