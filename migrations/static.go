package migrations

import "github.com/gobuffalo/packr/v2"

var box = packr.New("migrations", "../migrations")

// PackrBox exports this directory in Packr.
func PackrBox() *packr.Box {
	return box
}
