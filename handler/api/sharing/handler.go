package sharing

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinmonl/pinmonl/handler/api/apibody"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// HandleFind returns Share.
func HandleFind(
	shares store.ShareStore,
	sharetags store.SharetagStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := UserFrom(ctx)
		m, _ := ShareFrom(ctx)

		mtsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindMust,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		mts := mtsm[m.ID]

		response.JSON(w, apibody.NewSharing(m).WithOwner(u).WithMustTags(mts))
	}
}

// HandleListTags returns the Tags with kind "AnyTags" from Share.
func HandleListTags(sharetags store.SharetagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := ShareFrom(ctx)

		atsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindAny,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		ats := atsm[m.ID]

		resp := make([]interface{}, len(ats))
		for i, t := range ats {
			resp[i] = apibody.NewTag(t)
		}
		response.JSON(w, resp)
	}
}

// HandleListPinls returns the Pinls from Share.
func HandleListPinls(
	sharetags store.SharetagStore,
	pinls store.PinlStore,
	taggables store.TaggableStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := UserFrom(ctx)
		m, _ := ShareFrom(ctx)

		mtsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindMust,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		atsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindAny,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mtids := pluckTagIDs(mtsm[m.ID])
		atids := pluckTagIDs(atsm[m.ID])
		ps, err := pinls.List(ctx, &store.PinlOpts{
			UserID:     u.ID,
			MustTagIDs: mtids,
			AnyTagIDs:  atids,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ptm := model.MustBeMorphables(ps)
		ts, err := taggables.ListTags(ctx, &store.TaggableOpts{
			Targets: ptm,
			TagIDs:  atids,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		resp := make([]interface{}, len(ps))
		for i, p := range ps {
			pt := ts[p.ID]
			resp[i] = apibody.NewPinl(p).WithTags(pt)
		}
		response.JSON(w, resp)
	}
}

// HandleFindPinl returns Pinl from Share with detail information.
func HandleFindPinl(
	sharetags store.SharetagStore,
	pinls store.PinlStore,
	taggables store.TaggableStore,
	pkgs store.PkgStore,
	stats store.StatStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := UserFrom(ctx)
		m, _ := ShareFrom(ctx)

		mtsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindMust,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		atsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			ShareID: m.ID,
			Kind:    model.SharetagKindAny,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mtids := pluckTagIDs(mtsm[m.ID])
		atids := pluckTagIDs(atsm[m.ID])
		ps, err := pinls.List(ctx, &store.PinlOpts{
			ID:         chi.URLParam(r, "pinl"),
			UserID:     u.ID,
			MustTagIDs: mtids,
			AnyTagIDs:  atids,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		if len(ps) == 0 {
			response.BadRequest(w, nil)
			return
		}

		p := ps[0]
		ts, err := taggables.ListTags(ctx, &store.TaggableOpts{
			Targets: []model.Morphable{p},
			TagIDs:  atids,
		})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		// pps, err := pkgs.List(ctx, &store.PkgOpts{MonlURL: p.URL})
		// if err != nil {
		// 	response.InternalError(w, err)
		// 	return
		// }

		// pss, err := stats.List(ctx, &store.StatOpts{
		// 	PkgIDs:     (model.PkgList)(pps).Keys(),
		// 	WithLatest: true,
		// })
		// if err != nil {
		// 	response.InternalError(w, err)
		// 	return
		// }

		response.JSON(w, apibody.NewPinl(p).WithTags(ts[p.ID]))
	}
}

func pluckTagIDs(tags []model.Tag) []string {
	ids := make([]string, len(tags))
	for i, t := range tags {
		ids[i] = t.ID
	}
	return ids
}
