package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/model"
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

// pkgListHandler finds pkgs of the url.
func (s *Server) pkgListHandler(w http.ResponseWriter, r *http.Request) {
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
	u, err := monlutils.NormalizeURL(url)
	if err != nil {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}

	var (
		ctx    = r.Context()
		monl   *model.Monl
		isNew  bool
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		found, err := s.Monls.List(ctx, &store.MonlOpts{
			URL: u.String(),
		})
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
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
			err = s.Monls.Create(ctx, monl)
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

	if isNew {
		cherr := s.Queue.Add(job.NewMonlCreated(monl.ID))
		if err := <-cherr; err != nil {
			response.JSON(w, nil, http.StatusInternalServerError)
			return
		}
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

// pkgHandler finds pkg by the uri.
func (s *Server) pkgHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		pkg    *model.Pkg
		isNew  bool
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err error
		pu, _ := pkguriURLParam(r)
		pkg, err = s.Pkgs.FindURI(ctx, pu)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		if pkg == nil {
			isNew = true
			pkg = &model.Pkg{}
			if err = pkg.UnmarshalPkgURI(pu); err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
			if err = s.Pkgs.Create(ctx, pkg); err != nil {
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

	// Trigger to fetch stats if pkg is new.
	if isNew {
		cherr := s.Queue.Add(job.NewPkgSelfUpdate(pkg.ID))
		if err := <- cherr; err != nil {
			response.JSON(w, nil, http.StatusInternalServerError)
			return
		}
	}

	// Response.
	response.JSON(w, pkg, http.StatusOK)
}
