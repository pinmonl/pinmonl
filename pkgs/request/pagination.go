package request

import (
	"net/http"
	"strconv"

	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

type Paginator struct {
	Page     int64
	PageSize int64
}

func NewPaginatorFromRequest(r *http.Request, pageName, sizeName string, defaultSize int64) (*Paginator, error) {
	var (
		page, size int64
		err        error
		pageq      = r.URL.Query().Get(pageName)
		sizeq      = r.URL.Query().Get(sizeName)
	)

	if pageq != "" {
		if page, err = strconv.ParseInt(pageq, 10, 64); err != nil {
			return nil, err
		}
	} else {
		page = 1
	}

	if sizeq != "" {
		if size, err = strconv.ParseInt(sizeq, 10, 64); err != nil {
			return nil, err
		}
	} else {
		size = defaultSize
	}

	return &Paginator{
		Page:     page,
		PageSize: size,
	}, nil
}

func (p *Paginator) ToOpts() store.ListOpts {
	return store.ListOpts{
		Limit:  p.PageSize,
		Offset: p.PageSize * (p.Page - 1),
	}
}

func (p *Paginator) ToPageInfo(total int64) *response.PageInfo {
	return &response.PageInfo{
		TotalCount: total,
		Page:       p.Page,
		PageSize:   p.PageSize,
	}
}

func Pagination(pageName, sizeName string, defaultSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			p, err := NewPaginatorFromRequest(r, pageName, sizeName, defaultSize)
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
