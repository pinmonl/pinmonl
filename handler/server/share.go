package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// Errors for share requests.
var (
	ErrShareRequired = errors.New("slug is required")
	ErrShareSlugUsed = errors.New("slug is used")
)

// shareBody defines the fields of share.
type shareBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// sharePrepareHandler starts share perparation.
//
// Share is created if not found by slug or
// otherwise cleanup action is performed.
func (s *Server) sharePrepareHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.JSON(w, ErrShareRequired, http.StatusBadRequest)
		return
	}

	var in shareBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		share  *model.Share
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err error
		share, err = s.Shares.FindSlug(ctx, user.ID, slug)
		if err == nil && share != nil {
			err = storeutils.CleanupShare(ctx, s.Sharetags, s.Sharepins, s.Pinls, s.Taggables, share.ID, true)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
		}

		if share == nil {
			share = &model.Share{
				UserID: user.ID,
				Slug:   slug,
			}
		}
		share.Name = in.Name
		share.Description = in.Description
		share.Status = model.Preparing

		if share.ID == "" {
			err = s.Shares.Create(ctx, share)
		} else {
			err = s.Shares.Update(ctx, share)
		}
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
	response.JSON(w, share, http.StatusOK)
}

// sharePublishHandler publishes share to public.
func (s *Server) sharePublishHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		share  = request.ShareFrom(ctx)
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		share.Status = model.Active
		outerr = s.Shares.Update(ctx, share)
		return outerr == nil
	})

	if outerr != nil {
		response.JSON(w, outerr, http.StatusInternalServerError)
		return
	}
	response.JSON(w, share, http.StatusOK)
}

// bindShareBySlug binds share by slug into context
// and throws error if not found.
func (s *Server) bindShareBySlug() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := request.AuthedFrom(ctx)
			slug := chi.URLParam(r, "slug")

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

// shareStatusMustBe checks share from context is with status.
func (s *Server) shareStatusMustBe(status model.Status) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			share := request.ShareFrom(r.Context())
			if share.Status != status {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// sharetagBody defines the fields of sharetag.
type sharetagBody struct {
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
	Color    string `json:"color"`
	BgColor  string `json:"bgColor"`
}

// sharetagCreateHandler creates tag of share.
func (s *Server) sharetagCreateHandler(kind model.SharetagKind) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in sharetagBody
		err := request.JSON(r, &in)
		if err != nil {
			response.JSON(w, err, http.StatusBadRequest)
			return
		}
		if in.Name == "" {
			response.JSON(w, errors.New("tag name is required"), http.StatusBadRequest)
			return
		}

		var (
			ctx      = r.Context()
			user     = request.AuthedFrom(ctx)
			share    = request.ShareFrom(ctx)
			sharetag *model.Sharetag
			code     int
			outerr   error
		)
		s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
			tag, err := s.Tags.FindOrCreate(ctx, &model.Tag{
				UserID: user.ID,
				Name:   in.Name,
			})
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			tag.Color = in.Color
			tag.BgColor = in.BgColor
			err = s.Tags.Update(ctx, tag)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			st, err := storeutils.SaveSharetag(ctx, s.Sharetags, s.Tags, user.ID, share.ID, &model.Sharetag{
				TagID:    tag.ID,
				Kind:     kind,
				ParentID: in.ParentID,
			})
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			st.Tag = tag
			sharetag = st
			return true
		})

		if outerr != nil || response.IsError(code) {
			response.JSON(w, outerr, code)
			return
		}
		response.JSON(w, sharetag, http.StatusOK)
	}
}

// sharepinBody defines the fields of sharepin.
type sharepinBody struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// sharepinCreateHandler creates pinl of share.
func (s *Server) sharepinCreateHandler(w http.ResponseWriter, r *http.Request) {
	var in sharepinBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}
	if in.URL == "" {
		response.JSON(w, errors.New("url is required"), http.StatusBadRequest)
		return
	}

	var (
		ctx      = r.Context()
		share    = request.ShareFrom(ctx)
		sharepin *model.Sharepin
		code     int
		outerr   error
	)

	pinl := &model.Pinl{
		URL:         in.URL,
		Title:       in.Title,
		Description: in.Description,
	}
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		err = s.Pinls.Create(ctx, pinl)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		sharepin = &model.Sharepin{
			ShareID: share.ID,
			PinlID:  pinl.ID,
		}
		err = s.Sharepins.Create(ctx, sharepin)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		sharepin.Pinl = pinl
		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}

	s.Queue.Add(job.NewPinlUpdated(sharepin.PinlID))
	response.JSON(w, sharepin, http.StatusOK)
}

// shareDeleteHandler handles share delete request.
func (s *Server) shareDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		share  = request.ShareFrom(ctx)
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err error
		err = storeutils.CleanupShare(ctx, s.Sharetags, s.Sharepins, s.Pinls, s.Taggables, share.ID, true)
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
