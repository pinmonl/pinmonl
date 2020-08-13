package server

import (
	"errors"
	"net/http"

	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

// statListHandler lists stats of the pkg.
func (s *Server) statListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParseStatQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	if len(query.PkgIDs) == 0 && len(query.ParentIDs) == 0 {
		response.JSON(w, errors.New("either pkg or parent should be provided"), http.StatusBadRequest)
		return
	}

	var (
		ctx  = r.Context()
		pg   = request.PaginatorFrom(ctx)
		opts = &store.StatOpts{
			PkgIDs:    query.PkgIDs,
			ParentIDs: query.ParentIDs,
			IsLatest:  query.Latest,
			Kinds:     query.Kinds,
			ListOpts:  pg.ToOpts(),
			Orders:    []store.StatOrder{store.StatOrderByRecordDesc},
		}
	)

	if len(opts.ParentIDs) == 0 {
		opts.ParentIDs = []string{""}
	}

	stats, err := s.Stats.List(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	count, err := s.Stats.Count(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.ListJSON(w, stats, pg.ToPageInfo(count), http.StatusOK)
}
