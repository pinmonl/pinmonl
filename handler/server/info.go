package server

import (
	"net/http"

	"github.com/pinmonl/pinmonl/pkgs/response"
)

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	b := response.Body{
		"version": s.Version.String(),
	}
	response.JSON(w, b, http.StatusOK)
}
