package common

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/store"
)

func BindShare(shares *store.Shares, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx     = r.Context()
				user    = request.AuthedFrom(ctx)
				shareID = chi.URLParam(r, paramName)
			)

			share, err := shares.Find(ctx, shareID)
			if err != nil {
				response.JSON(w, err, http.StatusInternalServerError)
				return
			}
			if share == nil {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}
			if share.UserID != user.ID {
				response.JSON(w, nil, http.StatusNotFound)
				return
			}

			ctx = request.WithShare(ctx, share)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func BindShareBySlug(shares *store.Shares, paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx  = r.Context()
				user = request.AuthedFrom(ctx)
				slug = chi.URLParam(r, paramName)
			)

			share, err := shares.FindSlug(ctx, user.ID, slug)
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
