package api

import "github.com/shurcooL/githubv4"

// PageInfo defines the structure of pagination
type PageInfo struct {
	StartCursor     string
	EndCursor       string
	HasNextPage     bool
	HasPreviousPage bool
}

// PageOption defines the parameters for querying
type PageOption struct {
	After  string
	Before string
	First  int
	Last   int
}

// Scalar returns map of scalar for GraphQL
func (po PageOption) Scalar() map[string]interface{} {
	out := map[string]interface{}{
		"after":  (*githubv4.String)(nil),
		"before": (*githubv4.String)(nil),
		"first":  (*githubv4.Int)(nil),
		"last":   (*githubv4.Int)(nil),
	}
	if po.After != "" {
		out["after"] = githubv4.String(po.After)
	}
	if po.Before != "" {
		out["before"] = githubv4.String(po.Before)
	}
	if po.First != 0 {
		out["first"] = githubv4.Int(po.First)
	}
	if po.Last != 0 {
		out["last"] = githubv4.Int(po.Last)
	}
	return out
}
