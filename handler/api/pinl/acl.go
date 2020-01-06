package pinl

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/handler/api/user"
)

// RequireOwner restricts pinl to its owner only.
func RequireOwner() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			m, has := request.PinlFrom(r.Context())
			if !has {
				response.NotFound(w, nil)
				return
			}
			if user.MatchUser(w, r, m.UserID) {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
