package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/generate"
	"github.com/pinmonl/pinmonl/pkgs/passwd"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

var (
	ErrLoginRequired  = errors.New("login and password are required")
	ErrSignupRequired = errors.New("login, password and name are required")
	ErrLoginUsed      = errors.New("login is used")
)

func (s *Server) authenticate() func(http.Handler) http.Handler {
	return authenticate(s.TokenSecret, s.Users)
}

func authenticate(secret []byte, users *store.Users) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwtToken, err := request.ExtractJwtToken(r)
			if err != nil {
				response.JSON(w, err, http.StatusBadRequest)
				return
			}

			var claims authClaims
			token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if !token.Valid {
				response.JSON(w, err, http.StatusUnauthorized)
				return
			}
			if err := claims.Valid(); err != nil {
				response.JSON(w, err, http.StatusUnauthorized)
			}

			ctx := r.Context()
			user, err := users.Find(ctx, claims.UserID)
			if user == nil {
				response.JSON(w, err, http.StatusUnauthorized)
				return
			}
			if user.Hash != claims.Hash {
				response.JSON(w, nil, http.StatusUnauthorized)
				return
			}

			ctx = request.WithUser(ctx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

type loginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var in loginBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	if in.Login == "" || in.Password == "" {
		response.JSON(w, ErrLoginRequired, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := s.Users.FindLogin(ctx, in.Login)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	err = passwd.CompareString(user.Password, in.Password)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	user.LastSeen = field.Time(time.Now())
	err = s.Users.Update(ctx, user)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	token, err := s.generateToken(user)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, tokenResponse{
		Token: token,
	}, http.StatusOK)
}

type signupBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (s *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	var in signupBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	if in.Login == "" || in.Password == "" || in.Name == "" {
		response.JSON(w, ErrSignupRequired, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	found, err := s.Users.FindLogin(ctx, in.Login)
	if err == nil && found != nil {
		response.JSON(w, ErrLoginUsed, http.StatusBadRequest)
		return
	}

	user := &model.User{
		Login:    in.Login,
		Name:     in.Name,
		Hash:     generate.UserHash(),
		LastSeen: field.Time(time.Now()),
	}
	if pw, err := passwd.HashString(in.Password); err == nil {
		user.Password = pw
	} else {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	err = s.Users.Create(ctx, user)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	response.JSON(w, user, http.StatusOK)
}

func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := request.UserFrom(ctx)
	user.LastSeen = field.Time(time.Now())
	err := s.Users.Update(ctx, user)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, nil, http.StatusOK)
}

type authClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
	Hash   string `json:"hash"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (s *Server) generateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims{
		UserID: user.ID,
		Hash:   user.Hash,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(s.TokenExpire).Unix(),
			Issuer:    s.TokenIssuer,
		},
	})

	signed, err := token.SignedString(s.TokenSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}
