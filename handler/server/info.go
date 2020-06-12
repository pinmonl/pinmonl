package server

import (
	"net/http"

	"github.com/pinmonl/pinmonl/cmd/pinmonl-server/version"
	"github.com/pinmonl/pinmonl/pkgs/response"
)

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	b := response.Body{
		"version": version.Version.String(),
	}
	response.JSON(w, b, http.StatusOK)
}
