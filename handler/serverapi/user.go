package serverapi

import (
	"net/http"

	"github.com/pinmonl/pinmonl/store"
)

type signupBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Signup(users *store.Users) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	}
	return fn
}
