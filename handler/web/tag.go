package web

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/common"
)

func (s *Server) bindTag() func(http.Handler) http.Handler {
	return common.BindTag(s.Tags, "tag")
}

func (s *Server) tagListHandler(w http.ResponseWriter, r *http.Request) {
	h := common.TagListHandler(s.Tags)
	h.ServeHTTP(w, r)
}

func (s *Server) tagHandler(w http.ResponseWriter, r *http.Request) {
	h := common.TagHandler()
	h.ServeHTTP(w, r)
}

func (s *Server) tagCreateHandler(w http.ResponseWriter, r *http.Request) {
	h := common.TagCreateHandler(s.Txer, s.Tags)
	h.ServeHTTP(w, r)
}

func (s *Server) tagUpdateHandler(w http.ResponseWriter, r *http.Request) {
	h := common.TagUpdateHandler(s.Txer, s.Tags)
	h.ServeHTTP(w, r)
}

func (s *Server) tagDeleteHandler(w http.ResponseWriter, r *http.Request) {
	h := common.TagDeleteHandler(s.Txer, s.Tags, s.Taggables)
	h.ServeHTTP(w, r)
}
