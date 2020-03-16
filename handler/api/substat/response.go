package substat

import "github.com/pinmonl/pinmonl/model"

// Body defines the response body of Substat.
type Body struct {
	model.Substat
}

// NewBody creates the response body.
func NewBody(s model.Substat) Body {
	return Body{Substat: s}
}
