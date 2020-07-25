package web

import (
	"net/http"

	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

func (s *Server) pkgListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		pg  = request.PaginatorFrom(ctx)
	)

	opts := &store.MonpkgOpts{
		ListOpts: pg.ToOpts(),
	}
	if qMonls := request.QueryCsv(r, "monl"); len(qMonls) > 0 {
		opts.MonlIDs = qMonls
	}

	pList, err := s.Monpkgs.ListWithPkg(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, pList, http.StatusOK)
}
