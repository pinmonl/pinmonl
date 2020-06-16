package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/generate"
	"github.com/pinmonl/pinmonl/pkgs/passwd"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
)

// Errors.
var (
	ErrLoginRequired  = errors.New("login and password are required")
	ErrSignupRequired = errors.New("login, password and name are required")
	ErrLoginUsed      = errors.New("login is used")
)

// authenticate binds request user into context.
func (s *Server) authenticate() func(http.Handler) http.Handler {
	return request.Authenticate(s.TokenSecret, s.Users)
}

// authorize checks the request is from a valid user.
func (s *Server) authorize() func(http.Handler) http.Handler {
	return request.Authorize()
}

// authorizeUserOnly checks the request is from a valid user
// with normal user role.
func (s *Server) authorizeUserOnly() func(http.Handler) http.Handler {
	return request.Authorize(model.NormalUser)
}

type loginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// loginHandler validates the user credentials and
// returns an access token if succeeded.
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var in loginBody
	err1 := request.JSON(r, &in)
	if err1 != nil {
		response.JSON(w, err1, http.StatusBadRequest)
		return
	}

	if in.Login == "" || in.Password == "" {
		response.JSON(w, ErrLoginRequired, http.StatusBadRequest)
		return
	}

	var (
		ctx  = r.Context()
		user *model.User
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err2 error
		user, err2 = s.Users.FindLogin(ctx, in.Login)
		if err2 != nil {
			err, code = err2, http.StatusBadRequest
			return false
		}
		err2 = passwd.CompareString(user.Password, in.Password)
		if err2 != nil {
			err, code = err2, http.StatusBadRequest
			return false
		}

		user.LastSeen = field.Time(time.Now())
		err2 = s.Users.Update(ctx, user)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	s.printToken(w, user)
}

type signupBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// signupHandler creates a user and returns with an access token.
func (s *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	var in signupBody
	err1 := request.JSON(r, &in)
	if err1 != nil {
		response.JSON(w, err1, http.StatusBadRequest)
		return
	}

	if in.Login == "" || in.Password == "" || in.Name == "" {
		response.JSON(w, ErrSignupRequired, http.StatusBadRequest)
		return
	}

	var (
		ctx  = r.Context()
		user *model.User
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		found, err2 := s.Users.FindLogin(ctx, in.Login)
		if err2 == nil && found != nil {
			err, code = ErrLoginUsed, http.StatusBadRequest
			return false
		}

		user = &model.User{
			Login:    in.Login,
			Name:     in.Name,
			Hash:     generate.UserHash(),
			LastSeen: field.Time(time.Now()),
		}
		if pw, err2 := passwd.HashString(in.Password); err2 == nil {
			user.Password = pw
		} else {
			err, code = err2, http.StatusBadRequest
			return false
		}
		err2 = s.Users.Create(ctx, user)
		if err2 != nil {
			err, code = err2, http.StatusBadRequest
			return false
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	s.printToken(w, user)
}

// aliveHandler refreshes the user.LastSeen and recreates a
// new access token.
func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = request.AuthedFrom(ctx)
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		user.LastSeen = field.Time(time.Now())
		err = s.Users.Update(ctx, user)
		if err != nil {
			err, code = err, http.StatusInternalServerError
			return false
		}
		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	s.printToken(w, user)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (s *Server) printToken(w http.ResponseWriter, user *model.User) error {
	token, err := request.GenerateJwtToken(s.TokenIssuer, s.TokenExpire, s.TokenSecret, user)
	if err != nil {
		return response.JSON(w, err, http.StatusInternalServerError)
	}
	return response.JSON(w, tokenResponse{
		Token: token,
	}, http.StatusOK)
}
