package apibody

import "github.com/pinmonl/pinmonl/model"

// Stat defines the response body of Stat.
type Stat struct {
	model.Stat
	Substats []Stat `json:"substats"`
}

// NewStat creates the response body.
func NewStat(s model.Stat) Stat {
	return Stat{
		Stat:     s,
		Substats: make([]Stat, 0),
	}
}

// WithSubstats sets the value of substats.
func (s Stat) WithSubstats(ss ...Stat) Stat {
	s.Substats = ss
	return s
}
