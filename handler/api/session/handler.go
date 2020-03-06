package session

import (
	"errors"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkg/password"
	"github.com/pinmonl/pinmonl/session"
	"github.com/pinmonl/pinmonl/store"
)

var (
	// ErrLogin indicates the login fail message.
	ErrLogin = errors.New("please check your login and password")
)

// HandleCreate verifies user login and returns session data.
func HandleCreate(sess session.Store, users store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, ErrLogin)
			return
		}

		ctx := r.Context()
		m := model.User{Login: in.Login}
		err = users.FindLogin(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = password.Compare(m.Password, in.Password)
		if err != nil {
			response.BadRequest(w, ErrLogin)
			return
		}

		sv := &session.Values{UserID: m.ID}
		sr, err := sess.Set(w, sv)
		if err != nil {
			response.BadRequest(w, err)
			return
		}
		if sr.Data != nil {
			response.JSON(w, sr)
		} else {
			response.NoContent(w)
		}
	}
}

// HandleDelete removes session data.
func HandleDelete(sess session.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := sess.Del(w, r)
		if err != nil {
			response.BadRequest(w, err)
			return
		}
		response.NoContent(w)
	}
}
