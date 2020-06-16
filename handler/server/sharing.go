package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

// bindUserSharing binds and checks the share and its owner from url.
func (s *Server) bindUserSharing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx   = r.Context()
				login = chi.URLParam(r, "user")
				slug  = chi.URLParam(r, "share")
			)

			user, err := s.Users.FindLogin(ctx, login)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if user == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			share, err := s.Shares.FindSlug(ctx, user.ID, slug)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if share == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithUser(ctx, user)
			ctx = request.WithShare(ctx, share)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// getSharingHandler shows the information of the share and its owner.
func (s *Server) getSharingHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		user  = request.UserFrom(ctx)
		share = request.ShareFrom(ctx)
	)

	share.User = user
	response.JSON(w, share, http.StatusOK)
}

// listSharingPinlsHandler lists the pinls of the share.
func (s *Server) listSharingPinlsHandler(w http.ResponseWriter, r *http.Request) {
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

// listSharingTagsHandler lists the tags of the share.
func (s *Server) listSharingTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	stList, err := s.Sharetags.ListWithTag(ctx, &store.SharetagOpts{
		ShareIDs:       []string{share.ID},
		TagNamePattern: "%" + query.Query + "%",
		ListOpts:       pg.ToOpts(),
	})
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, stList.ViewTags(), http.StatusOK)
}
