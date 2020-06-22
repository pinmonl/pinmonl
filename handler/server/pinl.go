package server

import (
	"context"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/common"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

// bindPinl binds and checks the pinl from url param.
func (s *Server) bindPinl() func(http.Handler) http.Handler {
	return common.BindPinl(s.Pinls, "pinl")
}

// pinlListHandler lists pinls of authed user.
func (s *Server) pinlListHandler(w http.ResponseWriter, r *http.Request) {
	h := common.PinlListHandler(s.Pinls, s.Tags, s.Taggables)
	h.ServeHTTP(w, r)
}

// pinlCreateHandler creates pinl according to user request.
func (s *Server) pinlCreateHandler(w http.ResponseWriter, r *http.Request) {
	h := common.PinlCreateHandler(s.Txer, s.Pinls, s.Tags, s.Taggables, s.Queue)
	h.ServeHTTP(w, r)
}

// pinlClearHandler clears all pinls under the authed user.
func (s *Server) pinlClearHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		pList, err := s.Pinls.List(ctx, &store.PinlOpts{
			UserID: user.ID,
		})
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		for _, p := range pList {
			// Remove tag relations.
			_, err = s.Taggables.DeleteByTarget(ctx, p)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			// Remove pinl.
			_, err = s.Pinls.Delete(ctx, p.ID)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
		}

		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}
	response.JSON(w, nil, http.StatusNoContent)
}

// pinlDeleteHandler deletes pinl by the id from url param.
func (s *Server) pinlDeleteHandler(w http.ResponseWriter, r *http.Request) {
	h := common.PinlDeleteHandler(s.Txer, s.Pinls, s.Tags, s.Taggables, s.Queue)
	h.ServeHTTP(w, r)
}
