package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pinlutils"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
)

// bindPinl binds and checks the pinl from url param.
func (s *Server) bindPinl() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx    = r.Context()
				user   = request.AuthedFrom(ctx)
				pinlID = chi.URLParam(r, "pinl")
			)

			pinl, err := s.Pinls.Find(ctx, pinlID)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if pinl == nil {
				response.JSON(w, err, http.StatusNotFound)
				return
			}
			if user.ID != pinl.UserID {
				response.JSON(w, err, http.StatusNotFound)
				return
			}

			ctx = request.WithPinl(ctx, pinl)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// listPinlsHandler lists pinls of authed user.
func (s *Server) listPinlsHandler(w http.ResponseWriter, r *http.Request) {
	query, err := request.ParsePinlQuery(r)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	var (
		ctx  = r.Context()
		pg   = request.PaginatorFrom(ctx)
		user = request.AuthedFrom(ctx)
		opts = &store.PinlOpts{
			UserID:   user.ID,
			Query:    query.Query,
			ListOpts: pg.ToOpts(),
		}
	)

	if len(query.Tags) > 0 {
		tList, err := s.Tags.List(ctx, &store.TagOpts{
			UserID: user.ID,
			Names:  query.Tags,
		})
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		opts.TagIDs = tList.Keys()
	}

	pList, err := s.Pinls.List(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.JSON(w, pList, http.StatusOK)
}

type pinlBody struct {
	URL string `json:"url"`
}

// createPinlHandler creates pinl according to user request.
func (s *Server) createPinlHandler(w http.ResponseWriter, r *http.Request) {
	var in pinlBody
	err1 := request.JSON(r, &in)
	if err1 != nil {
		response.JSON(w, err1, http.StatusBadRequest)
		return
	}

	if in.URL == "" {
		response.JSON(w, errors.New("url cannot be empty"), http.StatusBadRequest)
		return
	}
	if !pinlutils.IsValidURL(in.URL) {
		response.JSON(w, errors.New("invalid url format"), http.StatusBadRequest)
		return
	}

	var (
		ctx  = r.Context()
		user = request.AuthedFrom(ctx)
		pinl *model.Pinl
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		pinl = &model.Pinl{
			UserID: user.ID,
			URL:    in.URL,
		}
		err2 := s.Pinls.Create(ctx, pinl)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}

	s.Queue.Add(job.NewPinlUpdated(pinl.ID, s.Pinls, s.Monls, s.Pkgs, s.Stats, s.Monpkgs))
	response.JSON(w, pinl, http.StatusOK)
}

// clearPinlsHandler clears all pinls under the authed user.
func (s *Server) clearPinlsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = request.AuthedFrom(ctx)
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		pList, err2 := s.Pinls.List(ctx, &store.PinlOpts{
			UserID: user.ID,
		})
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}

		for _, p := range pList {
			_, err2 = s.Pinls.Delete(ctx, p.ID)
			if err2 != nil {
				err, code = err2, http.StatusInternalServerError
				return false
			}
		}

		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	response.JSON(w, nil, http.StatusNoContent)
}

// deletePinlHandler deletes pinl by the id from url param.
func (s *Server) deletePinlHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		pinl = request.PinlFrom(ctx)
		code int
		err  error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		_, err2 := s.Pinls.Delete(ctx, pinl.ID)
		if err2 != nil {
			err, code = err2, http.StatusInternalServerError
			return false
		}
		return true
	})

	if err != nil || response.IsError(code) {
		response.JSON(w, err, code)
		return
	}
	response.JSON(w, nil, http.StatusNoContent)
}
