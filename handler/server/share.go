package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/handleutils"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
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

// prepareShareHandler starts share perparation.
//
// Share is created if not found by slug or
// otherwise cleanup action is performed.
func (s *Server) prepareShareHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.JSON(w, ErrShareRequired, http.StatusBadRequest)
		return
	}

	var (
		ctx   = r.Context()
		share *model.Share
		code  int
		err   error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		user := request.UserFrom(ctx)
		share, err = s.Shares.FindSlug(ctx, user.ID, slug)
		if err == nil && share != nil {
			err = cleanupShare(ctx, s.Sharetags, s.Sharepins, s.Pinls, s.Taggables, share.ID)
			if err != nil {
				code = http.StatusInternalServerError
				return false
			}
		}

		var in shareBody
		err = request.JSON(r, &in)
		if err != nil {
			code = http.StatusBadRequest
			return false
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
			code = http.StatusInternalServerError
			return false
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	response.JSON(w, share, http.StatusOK)
}

// cleanupShare clears the relationships of share.
func cleanupShare(
	ctx context.Context,
	sharetags *store.Sharetags,
	sharepins *store.Sharepins,
	pinls *store.Pinls,
	taggables *store.Taggables,
	shareID string,
) error {
	spList, err := sharepins.List(ctx, &store.SharepinOpts{
		ShareIDs: []string{shareID},
	})
	if err != nil {
		return err
	}

	// Clean up share's pins.
	for _, sp := range spList {
		_, err = taggables.DeleteByTaggable(ctx, model.Pinl{ID: sp.PinlID})
		if err != nil {
			return err
		}
		_, err = pinls.Delete(ctx, sp.PinlID)
		if err != nil {
			return err
		}
		_, err = sharepins.Delete(ctx, sp.ID)
		if err != nil {
			return err
		}
	}

	stList, err := sharetags.List(ctx, &store.SharetagOpts{
		ShareIDs: []string{shareID},
	})
	if err != nil {
		return err
	}

	// Clean up share's tags.
	for _, st := range stList {
		_, err = sharetags.Delete(ctx, st.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// publishShareHandler publishes share to public.
func (s *Server) publishShareHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		share = request.ShareFrom(ctx)
		err   error
	)

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		share.Status = model.Active
		err = s.Shares.Update(ctx, share)
		return err == nil
	})

	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, share, http.StatusOK)
}

// bindShareBySlug binds share by slug into context
// and throws error if not found.
func (s *Server) bindShareBySlug(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := request.UserFrom(ctx)
			slug := chi.URLParam(r, paramName)

			found, err := s.Shares.List(ctx, &store.ShareOpts{
				UserID: user.ID,
				Slug:   slug,
			})
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if len(found) == 0 {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			r = r.WithContext(
				request.WithShare(ctx, found[0]),
			)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// sharetagBody defines the fields of sharetag.
type sharetagBody struct {
	Name     string             `json:"name"`
	Level    int                `json:"level"`
	ParentID string             `json:"parentId"`
	Kind     model.SharetagKind `json:"kind"`
}

// createShareTagHandler creates tag of share.
func (s *Server) createShareTagHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		user     = request.UserFrom(ctx)
		share    = request.ShareFrom(ctx)
		sharetag *model.Sharetag
		code     int
		err      error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var in sharetagBody
		err2 := request.JSON(r, &in)
		if err2 != nil {
			err, code = err2, http.StatusBadRequest
			return false
		}

		if !model.IsValidSharetagKind(in.Kind) {
			code = http.StatusBadRequest
			return false
		}

		tag, err2 := handleutils.UpdateOrCreateTag(ctx, s.Tags, user.ID, in.Name, nil)
		if err != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		st, err2 := handleutils.UpdateOrCreateSharetag(ctx, s.Sharetags, share.ID, tag.ID, &model.Sharetag{
			Level:    in.Level,
			ParentID: in.ParentID,
			Kind:     in.Kind,
		})
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		st.Tag = tag
		sharetag = st
		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	response.JSON(w, sharetag, http.StatusOK)
}

// sharepinBody defines the fields of sharepin.
type sharepinBody struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// createSharePinlHandler creates pinl of share.
func (s *Server) createSharePinlHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		user     = request.UserFrom(ctx)
		share    = request.ShareFrom(ctx)
		sharepin *model.Sharepin
		code     int
		err      error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var in sharepinBody
		err2 := request.JSON(r, &in)
		if err2 != nil {
			err, code = err2, http.StatusBadRequest
			return false
		}

		pinl := &model.Pinl{
			UserID:      user.ID,
			URL:         in.URL,
			Title:       in.Title,
			Description: in.Description,
		}
		err2 = s.Pinls.Create(ctx, pinl)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		sharepin = &model.Sharepin{
			ShareID: share.ID,
			PinlID:  pinl.ID,
		}
		err2 = s.Sharepins.Create(ctx, sharepin)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		sharepin.Pinl = pinl
		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	response.JSON(w, sharepin, http.StatusOK)
}
