package response

// PageInfo stores info for pagination.
type PageInfo struct {
	TotalCount int
}

// WithPageInfo adds pagination into Body.
func WithPageInfo(b Body, p PageInfo) Body {
	b["pageInfo"] = p
	return b
}
