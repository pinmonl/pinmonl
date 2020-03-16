package pkg

import "github.com/pinmonl/pinmonl/model"

// Body defines the response body of Pkg.
type Body struct {
	model.Pkg
}

// NewBody creates the response body.
func NewBody(p model.Pkg) Body {
	return Body{Pkg: p}
}
