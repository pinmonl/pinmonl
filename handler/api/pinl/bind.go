package pinl

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// BindByID retrieves pinl by id and passes it into context.
func BindByID(pinls store.PinlStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := urlParamID(r)
			m := model.Pinl{ID: id}
			ctx := r.Context()
			err := pinls.Find(ctx, &m)
			if err != nil {
				response.NotFound(w, fmt.Errorf("id(%s) not found", id))
				return
			}
			ctx = request.WithPinl(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func urlParamID(r *http.Request) string {
	id := chi.URLParam(r, "pinlID")
	return id
}
