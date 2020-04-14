package stat

import (
	"github.com/pinmonl/pinmonl/model"
)

// Body defines the response body of Stat.
type Body struct {
	model.Stat
	Substats []Body `json:"substats"`
}

// NewBody creates the response body.
func NewBody(s model.Stat) Body {
	return Body{
		Stat:     s,
		Substats: make([]Body, 0),
	}
}

// WithSubstats sets the value of substats.
func (b Body) WithSubstats(ss ...Body) Body {
	b.Substats = ss
	return b
}
