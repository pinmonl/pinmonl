package web

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/common"
)

func (s *Server) bindImage() func(http.Handler) http.Handler {
	return common.BindImage(s.Images, "image")
}

func (s *Server) imageHandler(w http.ResponseWriter, r *http.Request) {
	h := common.ImageHandler()
	h.ServeHTTP(w, r)
}
