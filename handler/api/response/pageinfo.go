package response

// PageInfo stores info for pagination.
type PageInfo struct {
	Count int64 `json:"count"`
}

// NewPageInfo creates page info instance.
func NewPageInfo(count int64) PageInfo {
	return PageInfo{
		Count: count,
	}
}
