package monl

import "github.com/pinmonl/pinmonl/model"

// Body defines the response body of Monl.
type Body struct {
	model.Monl
}

// NewBody creates the response body.
func NewBody(m model.Monl) Body {
	return Body{Monl: m}
}
