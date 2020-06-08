package request

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/request"
)

func ExtractJwtToken(r *http.Request) (string, error) {
	if token := r.Header.Get("Authorization"); token != "" && strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		return token, nil
	}
	if token := r.URL.Query().Get("token"); token != "" {
		return token, nil
	}

	return "", request.ErrNoTokenInRequest
}
