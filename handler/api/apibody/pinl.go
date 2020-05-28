package apibody

import "github.com/pinmonl/pinmonl/model"

// Pinl defines the response body of Pinl.
type Pinl struct {
	model.Pinl
	Tags []string `json:"tags"`
	Pkgs []Pkg    `json:"pkgs"`
}

// NewPinl creates the response body.
func NewPinl(p model.Pinl) Pinl {
	return Pinl{
		Pinl: p,
		Tags: make([]string, 0),
		Pkgs: make([]Pkg, 0),
	}
}

// WithTags sets value of Tags.
func (p Pinl) WithTags(ts []model.Tag) Pinl {
	p.Tags = append(p.Tags, (model.TagList)(ts).PluckName()...)
	return p
}

// WithPkgs sets value of pkgs.
func (p Pinl) WithPkgs(pinls []model.Pkg, statMap map[string][]model.Stat) Pinl {
	p.Pkgs = make([]Pkg, len(pinls))
	for i, pinl := range pinls {
		pBody := NewPkg(pinl)
		if statMap != nil {
			pBody = pBody.WithStats(statMap[pinl.ID]...)
		}
		p.Pkgs[i] = pBody
	}
	return p
}
