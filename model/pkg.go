package model

import (
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

type Pkg struct {
	ID            string     `json:"id"`
	URL           string     `json:"url"`
	Provider      string     `json:"provider"`
	ProviderHost  string     `json:"providerHost"`
	ProviderURI   string     `json:"providerUri"`
	ProviderProto string     `json:"providerProto"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ImageID       string     `json:"imageId"`
	CustomUri     string     `json:"customUri"`
	FetchedAt     field.Time `json:"fetchedAt"`
	CreatedAt     field.Time `json:"createdAt"`
	UpdatedAt     field.Time `json:"updatedAt"`

	Stats *StatList `json:"stats,omitempty"`
}

func (p Pkg) MorphKey() string  { return p.ID }
func (p Pkg) MorphName() string { return "pkg" }

func (p Pkg) MarshalPkgURI() (*pkguri.PkgURI, error) {
	return &pkguri.PkgURI{
		Provider: p.Provider,
		Host:     p.ProviderHost,
		URI:      p.ProviderURI,
		Proto:    p.ProviderProto,
	}, nil
}

func (p *Pkg) UnmarshalPkgURI(pu *pkguri.PkgURI) error {
	p.Provider = pu.Provider
	p.ProviderHost = pu.Host
	p.ProviderURI = pu.URI
	p.ProviderProto = pu.Proto
	return nil
}

func (p *Pkg) SetStats(stats StatList) {
	p.Stats = &stats
}

type PkgList []*Pkg

func (pl PkgList) Keys() []string {
	keys := make([]string, len(pl))
	for i := range pl {
		keys[i] = pl[i].ID
	}
	return keys
}

func (pl PkgList) SetStats(stats StatList) {
	for i := range pl {
		ps := stats.GetPkgID(pl[i].ID)
		pl[i].SetStats(ps)
	}
}

var _ pkguri.Marshaler = &Pkg{}
var _ pkguri.Unmarshaler = &Pkg{}
