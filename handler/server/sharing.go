package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

// bindUserSharing binds and checks the share and its owner from url.
func (s *Server) bindUserSharing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx  = r.Context()
				user = request.UserFrom(ctx)
				slug = chi.URLParam(r, "share")
			)

			share, err := s.Shares.FindSlug(ctx, user.ID, slug)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if share == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithShare(ctx, share)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// sharingHandler shows the information of the share and its owner.
func (s *Server) sharingHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		user  = request.UserFrom(ctx)
		share = request.ShareFrom(ctx)
	)

	share.User = user
	response.JSON(w, share, http.StatusOK)
}

// sharingPinlListHandler lists the pinls of the share.
func (s *Server) sharingPinlListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParsePinlQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx   = r.Context()
		share = request.ShareFrom(ctx)
		pg    = request.PaginatorFrom(ctx)
	)

	spOpts := &store.SharepinOpts{
		ShareIDs:  []string{share.ID},
		PinlQuery: query.Query,
		ListOpts:  pg.ToOpts(),
	}

	if len(query.Tags) > 0 {
		stList, err := s.Sharetags.ListWithTag(ctx, &store.SharetagOpts{
			ShareIDs: []string{share.ID},
			TagNames: query.Tags,
		})
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		spOpts.TagIDs = stList.Tags().Keys()
	}

	spList, err := s.Sharepins.ListWithPinl(ctx, spOpts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, spList.Pinls(), http.StatusOK)
}

// sharingTagListHandler lists the tags with any kind of the share.
func (s *Server) sharingTagListHandler(w http.ResponseWriter, r *http.Request) {
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
		Kind:     field.NewNullValue(model.SharetagAny),
		ListOpts: pg.ToOpts(),
	}
	if query.Query != "" {
		opts.TagNamePattern = "%" + query.Query + "%"
	}
	if len(query.Names) > 0 {
		opts.TagNames = query.Names
	}

	stList, err := s.Sharetags.ListWithTag(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, stList.ViewTags(), http.StatusOK)
}
