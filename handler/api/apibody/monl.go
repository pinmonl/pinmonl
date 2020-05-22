package apibody

import "github.com/pinmonl/pinmonl/model"

// Monl defines the response body of Monl.
type Monl struct {
	model.Monl
}

// NewMonl creates the response body.
func NewMonl(m model.Monl) Monl {
	return Monl{Monl: m}
}
