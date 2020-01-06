package user

import (
	"fmt"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
)

// Authorize checks if the user exists.
func Authorize() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, has := request.UserFrom(r.Context())
			if !has {
				response.Unauthorized(w, fmt.Errorf("please login to proceed"))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// MatchUser reports whether the user matches with that in the session.
func MatchUser(w http.ResponseWriter, r *http.Request, userID string) bool {
	user, has := request.UserFrom(r.Context())
	if !has {
		response.Unauthorized(w, fmt.Errorf("please login to proceed"))
		return false
	}
	if user.ID != userID {
		response.Unauthorized(w, fmt.Errorf("do not have enough permission"))
		return false
	}
	return true
}

// Guest limits route to anonymous user.
func Guest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, has := request.UserFrom(r.Context())
			if has {
				response.BadRequest(w, fmt.Errorf("you have logged in already"))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
