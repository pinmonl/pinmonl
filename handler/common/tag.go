package common

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/pkgs/tagutils"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

func BindTag(tags *store.Tags, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx   = r.Context()
				user  = request.AuthedFrom(ctx)
				tagID = chi.URLParam(r, paramName)
			)

			tag, err := tags.Find(ctx, tagID)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if tag == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}
			if tag.UserID != user.ID {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithTag(ctx, tag)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func TagListHandler(tags *store.Tags) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		query, err := request.ParseTagQuery(r)
		if err != nil {
			response.JSON(w, err, http.StatusBadRequest)
			return
		}

		var (
			ctx  = r.Context()
			user = request.AuthedFrom(ctx)
			pg   = request.PaginatorFrom(ctx)
		)

		opts := &store.TagOpts{
			UserID:   user.ID,
			ListOpts: pg.ToOpts(),
		}
		if query.Query != "" {
			opts.NamePattern = tagutils.ToNamePattern(query.Query)
		}
		if len(query.Names) > 0 {
			opts.Names = query.Names
		}
		if _, has := r.URL.Query()["parent"]; has {
			if len(query.ParentIDs) > 0 {
				opts.ParentIDs = query.ParentIDs
			} else {
				opts.ParentIDs = []string{""}
			}
		}

		tList, err := tags.List(ctx, opts)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		count, err := tags.Count(ctx, opts)
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}

		response.ListJSON(w, tList, pg.ToPageInfo(count), http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func TagHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
			tag = request.TagFrom(ctx)
		)

		response.JSON(w, tag, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

type TagBody struct {
	Name    string `json:"name"`
	Color   string `json:"color"`
	BgColor string `json:"bgColor"`
}

func TagCreateHandler(txer database.Txer, tags *store.Tags) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var in TagBody
		err := request.JSON(r, &in)
		if err != nil {
			response.JSON(w, nil, http.StatusBadRequest)
			return
		}
		if in.Name == "" {
			response.JSON(w, errors.New("name is required"), http.StatusBadRequest)
			return
		}

		var (
			ctx    = r.Context()
			user   = request.AuthedFrom(ctx)
			code   int
			outerr error
		)

		found, err := tags.List(ctx, &store.TagOpts{
			UserID: user.ID,
			Name:   in.Name,
		})
		if err != nil {
			response.JSON(w, err, http.StatusInternalServerError)
			return
		}
		if len(found) > 0 {
			response.JSON(w, errors.New("name is used"), http.StatusBadRequest)
			return
		}

		tag := &model.Tag{
			UserID:  user.ID,
			Name:    in.Name,
			Color:   in.Color,
			BgColor: in.BgColor,
		}

		txer.TxFunc(ctx, func(ctx context.Context) bool {
			tag2, err := storeutils.SaveTag(ctx, tags, user.ID, tag)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
			tag = tag2
			return true
		})

		if outerr != nil || response.IsError(code) {
			response.JSON(w, outerr, code)
			return
		}
		response.JSON(w, tag, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func TagUpdateHandler(txer database.Txer, tags *store.Tags) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var in TagBody
		err := request.JSON(r, &in)
		if err != nil {
			response.JSON(w, nil, http.StatusBadRequest)
			return
		}
		if in.Name == "" {
			response.JSON(w, errors.New("name is required"), http.StatusBadRequest)
			return
		}

		var (
			ctx    = r.Context()
			user   = request.AuthedFrom(ctx)
			tag    = request.TagFrom(ctx)
			code   int
			outerr error
		)

		if in.Name != tag.Name {
			found, err := tags.List(ctx, &store.TagOpts{
				UserID: user.ID,
				Name:   in.Name,
			})
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if len(found) > 0 {
				response.JSON(w, errors.New("name is used"), http.StatusBadRequest)
				return
			}
		}

		tag.Name = in.Name
		tag.Color = in.Color
		tag.BgColor = in.BgColor

		txer.TxFunc(ctx, func(ctx context.Context) bool {
			tag2, err := storeutils.SaveTag(ctx, tags, user.ID, tag)
			if err != nil {
				outerr, code = err, http.StatusInternalServerError
				return false
			}
			tag = tag2
			return true
		})

		if outerr != nil || response.IsError(code) {
			response.JSON(w, outerr, code)
			return
		}
		response.JSON(w, tag, http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func TagDeleteHandler(txer database.Txer, tags *store.Tags, taggables *store.Taggables) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = r.Context()
			tag    = request.TagFrom(ctx)
			code   int
			outerr error
		)
		txer.TxFunc(ctx, func(ctx context.Context) bool {
			_, err := storeutils.DeleteTag(ctx, tags, taggables, tag)
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
