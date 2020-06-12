package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/handleutils"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue/job"
)

func providerURLParam(r *http.Request) string {
	return chi.URLParam(r, "provider")
}

func uriURLParam(r *http.Request) string {
	return chi.URLParam(r, "*")
}

func pkguriURLParam(r *http.Request) (*pkguri.PkgURI, error) {
	uri := providerURLParam(r) + "://" + uriURLParam(r)
	return pkguri.ParseProvider(uri)
}

func (s *Server) bindPkgURI() func(http.Handler) http.Handler {
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

func (s *Server) pkgLatestHandler(w http.ResponseWriter, r *http.Request) {
	pu, _ := pkguriURLParam(r)

	var (
		ctx   = r.Context()
		pkg   *model.Pkg
		isNew bool
		code  int
		err   error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		pkg, err = s.Pkgs.FindURI(ctx, pu)
		if err != nil {
			code = http.StatusInternalServerError
			return false
		}

		if pkg == nil {
			isNew = true
			pkg = &model.Pkg{}
			if err := pkg.UnmarshalPkgURI(pu); err != nil {
				code = http.StatusInternalServerError
				return false
			}
			if err := s.Pkgs.Create(ctx, pkg); err != nil {
				code = http.StatusInternalServerError
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
	}

	// Set pkg relationships.
	stats, err := handleutils.ListLatestStatsOfPkgs(ctx, s.Stats, []*model.Pkg{pkg})
	if err != nil {
		response.JSON(w, nil, http.StatusInternalServerError)
		return
	}
	pkg.Stats = &stats

	// Response.
	response.JSON(w, pkg, http.StatusOK)
}
