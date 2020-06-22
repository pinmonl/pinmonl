package common

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pinlutils"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// BindPinl finds pinl by id and binds into context.
func BindPinl(pinls *store.Pinls, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx    = r.Context()
				user   = request.AuthedFrom(ctx)
				pinlID = chi.URLParam(r, paramName)
			)

			pinl, err := pinls.Find(ctx, pinlID)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if pinl == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}
			if pinl.UserID != user.ID {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithPinl(ctx, pinl)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// PinlListHandler handles pinl listing request.
func PinlListHandler(pinls *store.Pinls, tags *store.Tags, taggables *store.Taggables) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		query, err := request.ParsePinlQuery(r)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		var (
			ctx  = r.Context()
			pg   = request.PaginatorFrom(ctx)
			user = request.AuthedFrom(ctx)
		)

		opts := &store.PinlOpts{
			UserID:   user.ID,
			Query:    query.Query,
			ListOpts: pg.ToOpts(),
		}
		if len(query.Tags) > 0 {
			opts.TagNames = query.Tags
		}

		pList, err := pinls.List(ctx, opts)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		tMap, err := storeutils.GetTags(ctx, taggables, pList.Morphables())
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		pList.SetTagNames(tMap)

		response.JSON(w, pList, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// PinlHandler shows pinl by id.
func PinlHandler(taggables *store.Taggables) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx  = r.Context()
			pinl = request.PinlFrom(ctx)
		)

		tMap, err := storeutils.GetTags(ctx, taggables, model.MorphableList{pinl})
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		pinl.SetTagNames(tMap[pinl.ID])

		response.JSON(w, pinl, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// PinlBody defines the fields of pinl body.
type PinlBody struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// PinlCreateHandler creates pinl from request.
func PinlCreateHandler(txer database.Txer, pinls *store.Pinls, tags *store.Tags, taggables *store.Taggables, queue *queue.Manager) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var in PinlBody
		err := request.JSON(r, &in)
		if err != nil {
			response.JSON(w, nil, http.StatusBadRequest)
			return
		}
		if in.URL == "" || !pinlutils.IsValidURL(in.URL) {
			response.JSON(w, errors.New("invalid url format"), http.StatusBadRequest)
			return
		}

		var (
			ctx    = r.Context()
			user   = request.AuthedFrom(ctx)
			pinl   *model.Pinl
			code   int
			outerr error
		)
		pinl = &model.Pinl{
			UserID:      user.ID,
			URL:         in.URL,
			Title:       in.Title,
			Description: in.Description,
		}

		txer.TxFunc(ctx, func(ctx context.Context) bool {
			err := pinls.Create(ctx, pinl)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			tList, err := storeutils.ReAssociateTags(ctx, tags, taggables, pinl, user.ID, in.Tags)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
			pinl.SetTagNames(tList)

			return true
		})

		if outerr != nil || response.IsError(code) {
			response.JSON(w, outerr, code)
			return
		}
		queue.Add(job.NewPinlUpdated(pinl.ID))
		response.JSON(w, pinl, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// PinlUpdateHandler updates fields of the found pinl.
func PinlUpdateHandler(txer database.Txer, pinls *store.Pinls, tags *store.Tags, taggables *store.Taggables, queue *queue.Manager) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var in PinlBody
		err := request.JSON(r, &in)
		if err != nil {
			response.JSON(w, nil, http.StatusBadRequest)
			return
		}
		if in.URL == "" || !pinlutils.IsValidURL(in.URL) {
			response.JSON(w, errors.New("invalid url format"), http.StatusBadRequest)
			return
		}

		var (
			ctx    = r.Context()
			user   = request.AuthedFrom(ctx)
			pinl   = request.PinlFrom(ctx)
			code   int
			outerr error
		)
		pinl.MonlID = ""
		pinl.URL = in.URL
		pinl.Title = in.Title
		pinl.Description = in.Description

		txer.TxFunc(ctx, func(ctx context.Context) bool {
			err := pinls.Update(ctx, pinl)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			tList, err := storeutils.ReAssociateTags(ctx, tags, taggables, pinl, user.ID, in.Tags)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
			pinl.SetTagNames(tList)

			return true
		})

		if outerr != nil || response.IsError(code) {
			response.JSON(w, outerr, code)
			return
		}
		queue.Add(job.NewPinlUpdated(pinl.ID))
		response.JSON(w, pinl, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

// PinlDeleteHandler deletes pinl if found.
func PinlDeleteHandler(txer database.Txer, pinls *store.Pinls, tags *store.Tags, taggables *store.Taggables, queue *queue.Manager) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = r.Context()
			pinl   = request.PinlFrom(ctx)
			code   int
			outerr error
		)
		txer.TxFunc(ctx, func(ctx context.Context) bool {
			_, err := taggables.DeleteByTarget(ctx, pinl)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}

			_, err = pinls.Delete(ctx, pinl.ID)
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
	return http.HandlerFunc(fn)
}
