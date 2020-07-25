package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/common"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// bindShare checks and binds share from request.
func (s *Server) bindShare() func(http.Handler) http.Handler {
	return common.BindShareBySlug(s.Shares, "slug")
}

// shareListHandler lists shares of the authenticated user.
func (s *Server) shareListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = request.AuthedFrom(ctx)
		pg   = request.PaginatorFrom(ctx)
	)

	opts := &store.ShareOpts{
		UserID:   user.ID,
		ListOpts: pg.ToOpts(),
	}

	sList, err := s.Shares.List(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	count, err := s.Shares.Count(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	stList, err := s.Sharetags.ListWithTag(ctx, &store.SharetagOpts{
		ShareIDs: sList.Keys(),
	})
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	sList.SetMustTagNames(stList.GetKind(model.SharetagMust).TagsByShare())
	sList.SetAnyTagNames(stList.GetKind(model.SharetagAny).TagsByShare())

	response.ListJSON(w, sList, pg.ToPageInfo(count), http.StatusOK)
}

// shareHandler finds share by id.
func (s *Server) shareHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		share = request.ShareFrom(ctx)
	)

	stList, err := s.Sharetags.ListWithTag(ctx, &store.SharetagOpts{
		ShareIDs: []string{share.ID},
	})
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	share.SetMustTagNames(stList.GetKind(model.SharetagMust).Tags())
	share.SetAnyTagNames(stList.GetKind(model.SharetagAny).Tags())

	response.JSON(w, share, http.StatusOK)
}

type shareBody struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MustTags    []string `json:"mustTags"`
	AnyTags     []string `json:"anyTags"`
}

// shareCreateHandler creates share.
func (s *Server) shareCreateHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.JSON(w, errors.New("slug is required"), http.StatusBadRequest)
		return
	}

	var in shareBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	if len(in.MustTags) == 0 {
		response.JSON(w, errors.New("must tag is required"), http.StatusBadRequest)
		return
	}
	in.AnyTags = cleanupAnyTags(in.MustTags, in.AnyTags)

	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		code   int
		outerr error
	)

	share, err := s.Shares.FindSlug(ctx, user.ID, slug)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	if share == nil {
		share = &model.Share{
			UserID: user.ID,
			Slug:   slug,
		}
	}
	share.Name = in.Name
	share.Description = in.Description

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err error
		if share.ID == "" {
			err = s.Shares.Create(ctx, share)
		} else {
			err = s.Shares.Update(ctx, share)
		}
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		if _, err := s.Sharetags.DeleteByShare(ctx, share.ID); err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		stmList, err := storeutils.AssociateSharetagByNames(ctx, s.Sharetags, s.Tags, user.ID, share.ID, model.SharetagMust, in.MustTags)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		staList, err := storeutils.AssociateSharetagByNames(ctx, s.Sharetags, s.Tags, user.ID, share.ID, model.SharetagAny, in.AnyTags)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		share.SetMustTagNames(stmList.Tags())
		share.SetAnyTagNames(staList.Tags())

		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}
	response.JSON(w, share, http.StatusOK)
}

// cleanupAnyTags removes the tags exists in mustTags from anyTags.
func cleanupAnyTags(mustTags, anyTags []string) []string {
	checks := make(map[string]int)
	for i := range mustTags {
		checks[mustTags[i]]++
	}
	out := make([]string, 0)
	for i := range anyTags {
		if _, skip := checks[anyTags[i]]; skip {
			continue
		}
		out = append(out, anyTags[i])
	}
	return out
}

// shareDeleteHandler deletes share.
func (s *Server) shareDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		share  = request.ShareFrom(ctx)
		code   int
		outerr error
	)

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		err := storeutils.CleanupShare(ctx, s.Sharetags, s.Sharepins, s.Pinls, s.Taggables, share.ID, false)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		_, err = s.Shares.Delete(ctx, share.ID)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}
		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}
	response.JSON(w, nil, http.StatusNoContent)
}

// sharePublishHandler publishes the share to pinmonl server.
func (s *Server) sharePublishHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		share = request.ShareFrom(ctx)
	)

	// TODO

	var _ = share
	if !s.Exchange.HasUser() {
		err := s.Exchange.Signup("test-user-1", "123", "test-user-1")
		if err != nil {
			err = s.Exchange.Login("test-user-1", "123")
		}
		fmt.Fprintf(w, "%v\n", err)
	}
	if !s.Exchange.HasMachine() {
		fmt.Fprintf(w, "%v\n", s.Exchange.MachineSignup())
	}
}

// sharetagListHandler lists sharetags of share.
func (s *Server) sharetagListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParseTagQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx   = r.Context()
		share = request.ShareFrom(ctx)
		pg    = request.PaginatorFrom(ctx)
	)

	opts := &store.SharetagOpts{
		ShareIDs: []string{share.ID},
		ListOpts: pg.ToOpts(),
		Kind:     field.NewNullValue(model.SharetagAny),
	}
	if query.Query != "" {
		opts.TagNamePattern = "%" + query.Query + "%"
	}
	if len(query.Names) > 0 {
		opts.TagNames = query.Names
	}
	if len(query.ParentIDs) > 0 {
		opts.ParentIDs = query.ParentIDs
	}

	stList, err := s.Sharetags.ListWithTag(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, stList.ViewTags(), http.StatusOK)
}

// sharetagCreateAnyKindHandler creates sharetag with any kind.
func (s *Server) sharetagCreateAnyKindHandler(w http.ResponseWriter, r *http.Request) {
	var in shareBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		share  = request.ShareFrom(ctx)
		code   int
		outerr error
	)

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		if _, err := s.Sharetags.DeleteByShareAndKind(ctx, share.ID, model.SharetagMust); err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		staList, err := storeutils.AssociateSharetagByNames(ctx, s.Sharetags, s.Tags, user.ID, share.ID, model.SharetagAny, in.AnyTags)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		share.SetAnyTagNames(staList.Tags())

		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}
	response.JSON(w, share, http.StatusOK)
}
