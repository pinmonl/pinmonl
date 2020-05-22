package share

import (
	"net/http"

	"github.com/pinmonl/pinmonl/handler/api/apibody"
	"github.com/pinmonl/pinmonl/handler/api/apiutils"
	"github.com/pinmonl/pinmonl/handler/api/request"
	"github.com/pinmonl/pinmonl/handler/api/response"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// HandleList returns shares.
func HandleList(shares store.ShareStore, sharetags store.SharetagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		ms, err := shares.List(ctx, &store.ShareOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mts, err := sharetags.ListTags(ctx, &store.SharetagOpts{
			Kind:     model.SharetagKindMust,
			ShareIDs: (model.ShareList)(ms).Keys(),
		})

		resp := make([]interface{}, len(ms))
		for i, m := range ms {
			resp[i] = apibody.NewShare(m).WithMustTags(mts[m.ID])
		}
		response.JSON(w, resp)
	}
}

// HandleFind returns share and its relations.
func HandleFind(shares store.ShareStore, sharetags store.SharetagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.ShareFrom(ctx)

		mtsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{ShareID: m.ID, Kind: model.SharetagKindMust})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		mts := mtsm[m.ID]

		atsm, err := sharetags.ListTags(ctx, &store.SharetagOpts{ShareID: m.ID, Kind: model.SharetagKindAny})
		if err != nil {
			response.InternalError(w, err)
			return
		}
		ats := atsm[m.ID]

		response.JSON(w, apibody.NewShare(m).WithMustTags(mts).WithAnyTags(ats))
	}
}

// HandleCreate validates and create share from user input.
func HandleCreate(shares store.ShareStore, sharetags store.SharetagStore, tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		m := model.Share{UserID: u.ID}
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = shares.Create(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mts, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.MustTags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.SharetagKindMust, rebuildMustTags(mts))
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ats, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.AnyTags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.SharetagKindAny, ats)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, apibody.NewShare(m).WithMustTags(mts).WithAnyTags(ats))
	}
}

// HandleUpdate validates and updates share from user input.
func HandleUpdate(shares store.ShareStore, sharetags store.SharetagStore, tags store.TagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := ReadInput(r.Body)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = in.Validate()
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		m, _ := request.ShareFrom(ctx)
		err = in.Fill(&m)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = shares.Update(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		mts, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.MustTags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.SharetagKindMust, rebuildMustTags(mts))
		if err != nil {
			response.InternalError(w, err)
			return
		}

		ats, err := apiutils.FindOrCreateTagsByName(ctx, tags, u, in.AnyTags)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		err = sharetags.ReAssocTags(ctx, m, model.SharetagKindAny, rebuildAnyTags(ats))
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, apibody.NewShare(m).WithMustTags(mts).WithAnyTags(ats))
	}
}

// HandleDelete removes share and its relations.
func HandleDelete(shares store.ShareStore, sharetags store.SharetagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		m, _ := request.ShareFrom(ctx)

		err := shares.Delete(ctx, &m)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = sharetags.ClearByKind(ctx, m, model.SharetagKindMust)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		err = sharetags.ClearByKind(ctx, m, model.SharetagKindAny)
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.NoContent(w)
	}
}

// HandlePageInfo returns the page info of Share.
func HandlePageInfo(shares store.ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, _ := request.UserFrom(ctx)
		count, err := shares.Count(ctx, &store.ShareOpts{UserID: u.ID})
		if err != nil {
			response.InternalError(w, err)
			return
		}

		response.JSON(w, response.NewPageInfo(count))
	}
}

func rebuildMustTags(tags []model.Tag) []model.Tag {
	for i := range tags {
		tags[i].ParentID = ""
	}
	return tags
}

func rebuildAnyTags(tags []model.Tag) []model.Tag {
	byID := make(map[string]model.Tag)
	for _, t := range tags {
		byID[t.ID] = t
	}
	for i, t := range tags {
		if _, has := byID[t.ParentID]; !has {
			tags[i].ParentID = ""
		}
	}
	return tags
}
