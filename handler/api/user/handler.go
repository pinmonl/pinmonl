package user

import (
	"errors"
	"fmt"
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

// HandleRegister handles and validates user registration.
func HandleRegister(sess session.Store, users store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		cu, err := users.List(ctx, &store.UserOpts{
			Login:    in.Login,
			ListOpts: store.ListOpts{Limit: 1},
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		if len(cu) > 0 {
			response.BadRequest(w, fmt.Errorf("login has already been used"))
			return
		}

		cu, err = users.List(ctx, &store.UserOpts{
			Email:    in.Email,
			ListOpts: store.ListOpts{Limit: 1},
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		if len(cu) > 0 {
			response.BadRequest(w, fmt.Errorf("email has already been used"))
			return
		}

		m := model.User{}
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = users.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		sv := &session.Values{UserID: m.ID}
		sess.Set(w, sv)
		response.JSON(w, Resp(m))
	}
}

// HandleLogin returns session from store after validating user credentials.
func HandleLogin(sess session.Store, users store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.ValidateLogin()
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
		sess.Set(w, sv)
		response.JSON(w, Resp(m))
	}
}

// HandleLogout clears the user session from store.
func HandleLogout(sess session.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess.Del(w, r)
		response.NoContent(w)
	}
}
