package web

import (
	"errors"
	"net/http"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pinmonl-go"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
)

func (s *Server) statListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParseStatQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	if len(query.PkgIDs) == 0 {
		response.JSON(w, errors.New("pkg is required"), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	pkg, err := s.Pkgs.Find(ctx, query.PkgIDs[0])
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	pu, err := pkg.MarshalPkgURI()
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	client := s.Exchange.UserClient()
	remotePkg, err := fetchPkg(client, pu.String())
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	remoteStats, err := fetchStats(client, remotePkg.ID, request.PaginatorFrom(ctx), query)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	response.JSON(w, remoteStats, http.StatusOK)
}

func fetchPkg(client *pinmonl.Client, url string) (*pinmonl.Pkg, error) {
	resp, err := client.PkgList(&pinmonl.PkgListOpts{URL: url})
	if err != nil {
		return nil, err
	}

	var directPkg *pinmonl.Pkg
	for i := range resp.Data {
		if resp.Data[i].Kind == pinmonl.MonpkgKind(model.MonpkgDirect) {
			directPkg = resp.Data[i].Pkg
			break
		}
	}
	if directPkg == nil {
		return nil, errors.New("no direct pkg")
	}
	return directPkg, nil
}

func fetchStats(client *pinmonl.Client, pkgID string, pg *request.Paginator, query *request.StatQuery) (*pinmonl.StatListResponse, error) {
	listopts := pinmonl.ListOpts{
		Page: pg.Page,
		Size: pg.PageSize,
	}
	if pg.PageSize == 0 {
		listopts.Size = -1
	}

	opts := &pinmonl.StatListOpts{
		Pkgs:     []string{pkgID},
		Latest:   query.Latest,
		Parents:  query.ParentIDs,
		ListOpts: listopts,
	}
	for _, k := range query.Kinds {
		opts.Kinds = append(opts.Kinds, string(k))
	}

	resp, err := client.StatList(opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
