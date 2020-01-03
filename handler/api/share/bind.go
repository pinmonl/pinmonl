package share

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// BindByID retrieves share by id and passes it into context.
func BindByID(shares store.ShareStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := urlParamID(r)
			m := model.Share{ID: id}
			ctx := r.Context()
			if err := shares.Find(ctx, &m); err != nil {
				response.NotFound(w, fmt.Errorf("id(%s) not found", id))
				return
			}
			ctx = request.WithShare(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func urlParamID(r *http.Request) string {
	id := chi.URLParam(r, "shareID")
	return id
}
