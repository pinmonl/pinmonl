package model

import (
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

type Pkg struct {
	ID           string     `json:"id"`
	URL          string     `json:"url"`
	Provider     string     `json:"provider"`
	ProviderHost string     `json:"providerHost"`
	ProviderURI  string     `json:"providerUri"`
	CreatedAt    field.Time `json:"createdAt"`
	UpdatedAt    field.Time `json:"updatedAt"`

	Stats *StatList `json:"stats,omitempty"`
}

func (p Pkg) MorphKey() string  { return p.ID }
func (p Pkg) MorphName() string { return "pkg" }

func (p Pkg) MarshalPkgURI() (*pkguri.PkgURI, error) {
	return &pkguri.PkgURI{
		Provider: p.Provider,
		Host:     p.ProviderHost,
		URI:      p.ProviderURI,
	}, nil
}

func (p *Pkg) UnmarshalPkgURI(pu *pkguri.PkgURI) error {
	p.Provider = pu.Provider
	p.ProviderHost = pu.Host
	p.ProviderURI = pu.URI
	return nil
}

type PkgList []*Pkg

func (pl PkgList) Keys() []string {
	var keys []string
	for _, p := range pl {
		keys = append(keys, p.ID)
	}
	return keys
}

var _ pkguri.Marshaler = &Pkg{}
var _ pkguri.Unmarshaler = &Pkg{}
