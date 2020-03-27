package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pinmonl/pinmonl/store"
)

// Pagination checks the query params and passes store.ListOpts into context.
func Pagination(defaultSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				q          = r.URL.Query()
				page int64 = 1
				size int64 = defaultSize
			)
			if qPage := q.Get("page"); qPage != "" {
				if iPage, err := strconv.ParseInt(qPage, 10, 64); err == nil && iPage > 0 {
					page = iPage
				}
			}
			if qSize := q.Get("pageSize"); qSize != "" {
				if iSize, err := strconv.ParseInt(qSize, 10, 64); err == nil && iSize > 0 {
					size = iSize
				}
			}

			opt := store.ListOpts{
				Limit:  size,
				Offset: (page - 1) * size,
			}
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), PaginationCtxKey, opt),
			))
		}
		return http.HandlerFunc(fn)
	}
}

// PaginationFrom gets store.ListOpts from context.
func PaginationFrom(ctx context.Context) *store.ListOpts {
	v, ok := ctx.Value(PaginationCtxKey).(store.ListOpts)
	if !ok {
		return nil
	}
	return &v
}
