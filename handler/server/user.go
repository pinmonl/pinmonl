package server

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/common"
)

// bindUser binds user by login.
func (s *Server) bindUser() func(http.Handler) http.Handler {
	return common.BindUserByLogin(s.Users, "user")
}
