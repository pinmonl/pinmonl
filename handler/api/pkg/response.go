package pkg

import (
	"github.com/pinmonl/pinmonl/handler/api/stat"
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Pkg.
type Body struct {
	model.Pkg
	Stats *[]stat.Body `json:"stats,omitempty"`
}

// NewBody creates the response body.
func NewBody(p model.Pkg) Body {
	return Body{Pkg: p}
}

// WithStats set the value of Body.Stats.
func (b Body) WithStats(stats ...model.Stat) Body {
	sb := make([]stat.Body, len(stats))
	for i, s := range stats {
		sb[i] = stat.NewBody(s)
	}
	b.Stats = &sb
	return b
}
