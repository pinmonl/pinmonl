package web

import (
	"errors"
	"net/http"

	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

func (s *Server) pkgListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParsePkgQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	// Checks query, either id and monl id must be provided.
	if len(query.IDs) == 0 {
		response.JSON(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}

	var (
		ctx = r.Context()
		pg  = request.PaginatorFrom(ctx)
	)

	// Retrieves pkgs depends on HasPinpkgs flag.
	opts := &store.PkgOpts{
		ListOpts: pg.ToOpts(),
		IDs:      query.IDs,
	}

	pList, err := s.Pkgs.List(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	count, err := s.Pkgs.Count(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	// Load latest stats.
	if err := storeutils.LoadPkgsLatestStats(ctx, s.Stats, pList); err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	response.ListJSON(w, pList, pg.ToPageInfo(count), http.StatusOK)
}
