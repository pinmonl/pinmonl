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
func (ap Pinl) WithPkgs(ps []model.Pkg, statMap map[string][]model.Stat) Pinl {
	ap.Pkgs = make([]Pkg, len(ps))
	for i, p := range ps {
		pBody := NewPkg(p)
		if statMap != nil {
			pBody = pBody.WithStats(statMap[p.ID]...)
		}
		ap.Pkgs[i] = pBody
	}
	return ap
}
