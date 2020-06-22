package server

import (
	"net/http"

	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// statListLatestHandler shows only latest stats of the pkg.
func (s *Server) statListLatestHandler(w http.ResponseWriter, r *http.Request) {
	h := s.statListBaseHandler(field.NewNullBool(true))
	h.ServeHTTP(w, r)
}

// statListHandler lists stats of the pkg.
func (s *Server) statListHandler(w http.ResponseWriter, r *http.Request) {
	h := s.statListBaseHandler(field.NullBool{})
	h.ServeHTTP(w, r)
}

// statListBaseHandler lists stats of the pkg.
func (s *Server) statListBaseHandler(defaultLatest field.NullBool) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		query, err := request.ParseStatQuery(r)
		if err != nil {
			response.JSON(w, err, http.StatusBadRequest)
			return
		}

		var (
			ctx = r.Context()
			pkg = request.PkgFrom(ctx)
			pg  = request.PaginatorFrom(ctx)
		)
		opts := &store.StatOpts{
			PkgIDs:   []string{pkg.ID},
			Kind:     query.Kind,
			Orders:   []store.StatOrder{store.StatOrderByRecordDesc},
			ListOpts: pg.ToOpts(),
		}
		// If defaultLatest is not null, it forces to overwrite
		// the url query.
		if defaultLatest.Valid {
			opts.IsLatest = defaultLatest
		} else {
			opts.IsLatest = query.Latest
		}

		// Search for stats.
		sList, err := s.Stats.List(ctx, opts)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		// Propagate with substats.
		sList, err = storeutils.ListStatTree(ctx, s.Stats, sList)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		response.JSON(w, sList, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
