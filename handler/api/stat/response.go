package stat

import (
	"github.com/pinmonl/pinmonl/handler/api/substat"
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Stat.
type Body struct {
	model.Stat
	Substats []substat.Body `json:"substats"`
}

// NewBody creates the response body.
func NewBody(s model.Stat) Body {
	return Body{
		Stat:     s,
		Substats: make([]substat.Body, 0),
	}
}

// WithSubstats sets the value of substats.
func (b Body) WithSubstats(ss []model.Substat) Body {
	b.Substats = make([]substat.Body, len(ss))
	for i, s := range ss {
		b.Substats[i] = substat.NewBody(s)
	}
	return b
}
