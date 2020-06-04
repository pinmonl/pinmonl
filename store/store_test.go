package store

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

type testWrappedListOpts struct {
	ListOpts
}

func TestPaginator(t *testing.T) {
	var (
		base      = squirrel.Select("*").From("table")
		baseQuery = "SELECT * FROM table"
	)

	var tests = []struct {
		pt    Paginator
		query string
	}{
		{
			pt:    nil,
			query: baseQuery,
		},
		{
			pt:    ListOpts{Limit: 10, Offset: 0},
			query: baseQuery + " LIMIT 10 OFFSET 0",
		},
		{
			pt: &testWrappedListOpts{ListOpts: ListOpts{
				Limit:  10,
				Offset: 0,
			}},
			query: baseQuery + " LIMIT 10 OFFSET 0",
		},
	}

	for _, test := range tests {
		b := addPagination(base, test.pt)
		got, _, _ := b.ToSql()
		want := test.query
		assert.Equal(t, want, got)
	}
}
