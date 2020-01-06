package tag

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// BindByID retrieves tag by id and passes it into context.
func BindByID(name string, tags store.TagStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, name)
			m := model.Tag{ID: id}
			ctx := r.Context()
			if err := tags.Find(ctx, &m); err != nil {
				response.NotFound(w, fmt.Errorf("id(%s) not found", id))
				return
			}
			ctx = request.WithTag(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
