package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/apibody"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

var (
	// ErrLogin indicates the login fail message.
	ErrLogin = errors.New("please check your login and password")
)

// HandleCreate creates User.
func HandleCreate(users store.UserStore) http.HandlerFunc {
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

		response.JSON(w, apibody.NewUser(m))
	}
}

// HandleGetMe returns the user profile of current user.
func HandleGetMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, _ := request.UserFrom(r.Context())
		response.JSON(w, apibody.NewUser(m))
	}
}

// HandleUpdateMe updates the user profile of current user.
func HandleUpdateMe(users store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.ValidateUpdate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		m, _ := request.UserFrom(ctx)
		err = in.FillDirty(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = users.Update(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, apibody.NewUser(m))
	}
}
