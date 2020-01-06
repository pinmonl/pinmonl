package user

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
)

// Authenticate retrieves user from the session.
func Authenticate(sess session.Store, users store.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, requestWithUser(r, sess, users))
		}
		return http.HandlerFunc(fn)
	}
}

func requestWithUser(r *http.Request, sess session.Store, users store.UserStore) *http.Request {
	val, err := sess.Get(r)
	if err != nil {
		logx.Errorf("api: fails to authenticate user, err: %s", err)
		return r
	}

	ctx := r.Context()
	user := model.User{ID: val.UserID}
	err = users.Find(ctx, &user)
	if err != nil {
		logx.Errorf("api: authenticate err: %s", err)
		return r
	}

	ctx = request.WithUser(ctx, user)
	return r.WithContext(ctx)
}
