package request

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

func ExtractJwtToken(r *http.Request) (string, error) {
	if token := r.Header.Get("Authorization"); token != "" && strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		if len(strings.Split(token, ".")) != 3 {
			return "", request.ErrNoTokenInRequest
		}
		return token, nil
	}
	if token := r.URL.Query().Get("token"); token != "" {
		return token, nil
	}

	return "", request.ErrNoTokenInRequest
}

func ParseJwtClaims(jwtToken string, secret []byte, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if !token.Valid {
		return err
	}
	if err := claims.Valid(); err != nil {
		return err
	}
	return nil
}

type AuthClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
	Hash   string `json:"hash"`
}

func Authenticate(secret []byte, users *store.Users) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			defer func() {
				next.ServeHTTP(w, r.WithContext(ctx))
			}()

			jwtToken, err := ExtractJwtToken(r)
			if err != nil {
				return
			}

			var claims AuthClaims
			err = ParseJwtClaims(jwtToken, secret, &claims)
			if err != nil {
				return
			}

			user, err := users.Find(ctx, claims.UserID)
			if user == nil {
				return
			}
			if user.Hash != claims.Hash {
				return
			}

			ctx = WithAuthed(ctx, user)
		}
		return http.HandlerFunc(fn)
	}
}

func Authorize(roles ...model.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user := AuthedFrom(r.Context())
			if user == nil {
				response.JSON(w, nil, http.StatusUnauthorized)
				return
			}

			if len(roles) > 0 {
				inRole := false
				for _, role := range roles {
					if role == user.Role {
						inRole = true
						break
					}
				}
				if !inRole {
					response.JSON(w, nil, http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GenerateJwtToken(issuer string, expireAfter time.Duration, secret []byte, user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		UserID: user.ID,
		Hash:   user.Hash,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(expireAfter).Unix(),
			Issuer:    issuer,
		},
	})

	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}
