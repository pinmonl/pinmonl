package image

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// BindByID retrieves image by id and passes it into context.
func BindByID(name string, images store.ImageStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, name)
			m := model.Image{ID: id}
			ctx := r.Context()
			if err := images.Find(ctx, &m); err != nil {
				response.NotFound(w, fmt.Errorf("id(%s) not found", id))
				return
			}
			ctx = request.WithImage(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
