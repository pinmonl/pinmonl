package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
)

type pkgListQuery struct {
	URL string
}

func parsePkgListQuery(r *http.Request) (*pkgListQuery, error) {
	return &pkgListQuery{
		URL: r.URL.Query().Get("url"),
	}, nil
}

// pkgListHandler finds pkgs of the url.
func (s *Server) pkgListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := parsePkgListQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	if query.URL == "" {
		response.JSON(w, errors.New("url is required"), http.StatusBadRequest)
		return
	}

	u, err := url.Parse(query.URL)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		providerName = u.Scheme
		isHttp       = monlutils.IsHttp(query.URL)
	)

	if !isHttp && !monler.Has(providerName) {
		response.JSON(w, errors.New("unsupported protocol"), http.StatusBadRequest)
		return
	}
	if !isHttp {
		if u.Host == "" && u.Path == "" {
			response.JSON(w, errors.New("uri is required"), http.StatusBadRequest)
			return
		}
	}

	var (
		ctx    = r.Context()
		monl   *model.Monl
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

		if len(found) > 0 {
			monl = found[0]
		} else {
			monl = &model.Monl{URL: u.String()}
			err := s.Monls.Create(ctx, monl)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
		}

		return true
	})

	if monl.FetchedAt.Time().IsZero() {
		// Enqueue
		cherr := s.Queue.Add(job.NewMonlCrawler(monl.ID))
		<-cherr

		// Reload data of monl.
		if monl2, err := s.Monls.Find(ctx, monl.ID); err == nil {
			monl = monl2
		} else {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	pg := request.PaginatorFrom(ctx)
	opts := &store.MonpkgOpts{
		MonlIDs:  []string{monl.ID},
		ListOpts: pg.ToOpts(),
	}

	pkgs, err := s.Monpkgs.ListWithPkg(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	count, err := s.Monpkgs.Count(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.ListJSON(w, pkgs, pg.ToPageInfo(count), http.StatusOK)
}
