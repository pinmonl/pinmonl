package web

import (
	"net/http"

	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

func (s *Server) statListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParseStatQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx = r.Context()
		pg  = request.PaginatorFrom(ctx)
	)

	opts := &store.StatOpts{
		ListOpts: pg.ToOpts(),
		Kinds:    query.Kinds,
		IsLatest: query.Latest,
	}
	if len(query.PkgIDs) > 0 {
		opts.PkgIDs = query.PkgIDs
	}

	sList, err := s.Stats.List(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	sList, err = storeutils.ListStatTree(ctx, s.Stats, sList)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, sList, http.StatusOK)
}
