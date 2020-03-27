package pmapi

import (
	"fmt"
	"net/url"
)

// PageOpts defines the pagination parameter.
type PageOpts struct {
	Page     int64
	PageSize int64
}

// Query converts PageOpts to query params.
func (p PageOpts) Query() url.Values {
	query := url.Values{}
	if p.Page > 0 {
		query.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PageSize > 0 {
		query.Set("page_size", fmt.Sprintf("%d", p.PageSize))
	}
	return query
}

// PageInfo defines the body of page-info API.
type PageInfo struct {
	Count int64 `json:"count"`
}
