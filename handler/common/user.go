package common

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

func BindUserByLogin(users *store.Users, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx   = r.Context()
				login = chi.URLParam(r, paramName)
			)

			user, err := users.FindLogin(ctx, login)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if user == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithUser(ctx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
