package store

import (
	"fmt"

	"github.com/pinmonl/pinmonl/database"
)

// ListOpts defines the options for listing.
type ListOpts struct {
	Limit   int64
	Offset  int64
	OrderBy map[string]string
}

func bindListOpts(opts ListOpts) database.SelectBuilder {
	br := database.SelectBuilder{}

	for col, order := range opts.OrderBy {
		br.OrderBy = append(br.OrderBy, fmt.Sprintf("%s %s", col, order))
	}
	if opts.Limit > 0 {
		br.Limit = opts.Limit
	}
	if opts.Offset > 0 {
		br.Offset = opts.Offset
	}

	return br
}
