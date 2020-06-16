package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/handleutils"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/pinlutils"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
)

// providerURLParam extracts the provider param from url.
func providerURLParam(r *http.Request) string {
	return chi.URLParam(r, "provider")
}

// uriURLParam extracts the uri param from url.
func uriURLParam(r *http.Request) string {
	return chi.URLParam(r, "*")
}

// pkguriURLParam extracts the pkguri param from url.
func pkguriURLParam(r *http.Request) (*pkguri.PkgURI, error) {
	uri := providerURLParam(r) + "://" + uriURLParam(r)
	return pkguri.ParseProvider(uri)
}

// checkPkgURI checks the format of the pkguri.
func (s *Server) checkPkgURI() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, err := pkguriURLParam(r)
			if err != nil {
				response.JSON(w, err, http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// bindPkgByURI binds and checks the existence of
func (s *Server) bindPkgByURI() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			pu, _ := pkguriURLParam(r)
			ctx := r.Context()
			pkg, err := s.Pkgs.FindURI(ctx, pu)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if pkg == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithPkg(ctx, pkg)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// listPkgsHandler finds pkgs of the url.
func (s *Server) listPkgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		proto = chi.URLParam(r, "proto")
		url   = proto + "://" + chi.URLParam(r, "*")
	)
	if proto != "https" {
		response.JSON(w, errors.New("currently only https is supported"), http.StatusBadRequest)
		return
	}

	// Checks url is in correct format and is not empty.
	if url == "" {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}
	if !pinlutils.IsValidURL(url) {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}
	u, err1 := monlutils.NormalizeURL(url)
	if err1 != nil {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}

	var (
		ctx   = r.Context()
		monl  *model.Monl
		isNew bool
		code  int
		err   error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		found, err2 := s.Monls.List(ctx, &store.MonlOpts{
			URL: u.String(),
		})
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		// Create monl if not found.
		if len(found) > 0 {
			monl = found[0]
		} else {
			isNew = true
			monl = &model.Monl{
				URL: u.String(),
			}
			err2 = s.Monls.Create(ctx, monl)
			if err2 != nil {
				err, code = err2, http.StatusInternalServerError
				return false
			}
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}

	if isNew {
		s.Queue.Add(job.NewMonlCreated(monl.ID, s.Monls, s.Pkgs, s.Stats, s.Monpkgs))
		response.JSON(w, nil, http.StatusCreated)
		return
	}

	// Response with the pkgs under the monl.
	pList, err := s.Monpkgs.ListWithPkg(ctx, &store.MonpkgOpts{
		MonlIDs: []string{monl.ID},
	})
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, pList.Pkgs(), http.StatusOK)
}

// findPkgHandler finds pkg by the uri.
func (s *Server) findPkgHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		pkg   *model.Pkg
		isNew bool
		code  int
		err   error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err2 error
		pu, _ := pkguriURLParam(r)
		pkg, err2 = s.Pkgs.FindURI(ctx, pu)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		if pkg == nil {
			isNew = true
			pkg = &model.Pkg{}
			if err2 = pkg.UnmarshalPkgURI(pu); err2 != nil {
				err, code = err2, http.StatusInternalServerError
				return false
			}
			if err2 = s.Pkgs.Create(ctx, pkg); err2 != nil {
				err, code = err2, http.StatusInternalServerError
				return false
			}
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}

	// Trigger to fetch stats if pkg is new.
	if isNew {
		s.Queue.Add(job.NewPkgSelfUpdate(pkg.ID, s.Pkgs, s.Stats))
		response.JSON(w, nil, http.StatusCreated)
		return
	}

	// Response.
	response.JSON(w, pkg, http.StatusOK)
}

// listLatestStatsHandler shows only latest stats of the pkg.
func (s *Server) listLatestStatsHandler(w http.ResponseWriter, r *http.Request) {
	h := s.listStatsBaseHandler(field.NewNullBool(true))
	h.ServeHTTP(w, r)
}

// listStatsHandler lists stats of the pkg.
func (s *Server) listStatsHandler(w http.ResponseWriter, r *http.Request) {
	h := s.listStatsBaseHandler(field.NullBool{})
	h.ServeHTTP(w, r)
}

// listStatsBaseHandler lists stats of the pkg.
func (s *Server) listStatsBaseHandler(defaultLatest field.NullBool) http.Handler {
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
		sList, err = handleutils.ListStatsTree(ctx, s.Stats, sList)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		response.JSON(w, sList, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
