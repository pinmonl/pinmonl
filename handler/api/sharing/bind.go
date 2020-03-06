package sharing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// BindUser retrieves User by the name of User.
func BindUser(name string, users store.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userName := chi.URLParam(r, name)

			ctx := r.Context()
			m := model.User{Login: userName}
			err := users.FindLogin(ctx, &m)
			if err != nil {
				response.NotFound(w, nil)
				return
			}

			ctx = WithUser(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// BindUserShare retrieves Share by the name of Share.
func BindUserShare(name string, shares store.ShareStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			shareName := chi.URLParam(r, name)

			ctx := r.Context()
			u, _ := UserFrom(ctx)
			m := model.Share{UserID: u.ID, Name: shareName}
			err := shares.FindByName(ctx, &m)
			if err != nil {
				response.NotFound(w, nil)
				return
			}

			ctx = WithShare(ctx, m)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
