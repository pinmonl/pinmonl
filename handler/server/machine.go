package server

import (
	"context"
	"net/http"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/generate"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/rs/xid"
)

// machineSignupHandler creates user with role = model.MachineUser
// and returns an access token.
func (s *Server) machineSignupHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		mach   *model.User
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		mach = &model.User{
			Login:    xid.New().String(),
			Role:     model.MachineUser,
			Hash:     generate.UserHash(),
			LastSeen: field.Now(),
		}
		err := s.Users.Create(ctx, mach)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}
		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}
	s.printToken(w, mach)
}
