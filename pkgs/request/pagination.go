package request

import (
	"net/http"
	"strconv"

	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

type Paginator struct {
	Page     int
	PageSize int
}

func NewPaginatorFromRequest(r *http.Request, pageName, sizeName string) (*Paginator, error) {
	var (
		page, size int
		err        error
	)
	if page, err = strconv.Atoi(r.URL.Query().Get(pageName)); err != nil {
		return nil, err
	}
	if size, err = strconv.Atoi(r.URL.Query().Get(sizeName)); err != nil {
		return nil, err
	}
	return &Paginator{
		Page:     page,
		PageSize: size,
	}, nil
}

func (p *Paginator) ToOpts() store.ListOpts {
	return store.ListOpts{
		Limit:  p.PageSize,
		Offset: p.Page * (p.PageSize - 1),
	}
}

func Pagination(pageName, sizeName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			p, err := NewPaginatorFromRequest(r, pageName, sizeName)
			if err != nil {
				response.JSON(w, err, http.StatusBadRequest)
				return
			}

			r = r.WithContext(
				WithPaginator(r.Context(), p),
			)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
