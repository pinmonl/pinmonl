package apibody

import "github.com/pinmonl/pinmonl/model"

// Pkg defines the response body of Pkg.
type Pkg struct {
	model.Pkg
	Stats *[]Stat `json:"stats,omitempty"`
}

// NewPkg creates the response body.
func NewPkg(p model.Pkg) Pkg {
	return Pkg{Pkg: p}
}

// WithStats set the value of Body.Stats.
func (p Pkg) WithStats(stats ...model.Stat) Pkg {
	sb := make([]Stat, len(stats))
	for i, s := range stats {
		sb[i] = NewStat(s)
	}
	p.Stats = &sb
	return p
}
