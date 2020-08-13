package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/pinmonl/pinmonl/handler/common"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/pinlutils"
	"github.com/pinmonl/pinmonl/pkgs/request"
	"github.com/pinmonl/pinmonl/pkgs/response"
	"github.com/pinmonl/pinmonl/pkgs/tagutils"
	"github.com/pinmonl/pinmonl/pubsub/message"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

func (s *Server) bindPinl() func(http.Handler) http.Handler {
	return common.BindPinl(s.Pinls, "pinl")
}

func (s *Server) pinlListHandler(w http.ResponseWriter, r *http.Request) {
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
		NoTag:    query.NoTag,
		Orders:   []store.PinlOrder{store.PinlOrderByLatest},
	}
	if len(query.Tags) > 0 {
		opts.TagNamePatterns = make([]string, len(query.Tags))
		for i := range query.Tags {
			opts.TagNamePatterns[i] = tagutils.ToNamePattern(query.Tags[i])
		}
	}

	pList, err := storeutils.ListPinlsWithLatestStats(ctx, s.Pinls, s.Monpkgs, s.Stats, s.Taggables, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	count, err := s.Pinls.Count(ctx, opts)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}

	response.ListJSON(w, pList, pg.ToPageInfo(count), http.StatusOK)
}

func (s *Server) pinlHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		pinl = request.PinlFrom(ctx)
	)

	pinl2, err := storeutils.PinlWithLatestStats(ctx, s.Pinls, s.Monpkgs, s.Stats, s.Taggables, pinl.ID)
	if err != nil {
		response.JSON(w, err, http.StatusInternalServerError)
		return
	}
	response.JSON(w, pinl2, http.StatusOK)
}

type pinlBody common.PinlBody

func validatePinlBody(in pinlBody) error {
	if in.URL == "" || !pinlutils.IsValidURL(in.URL) {
		return errors.New("invalid url format")
	}
	return nil
}

func savePinl(
	ctx context.Context,
	pinls *store.Pinls, images *store.Images,
	tags *store.Tags, taggables *store.Taggables,
	monls *store.Monls, monpkgs *store.Monpkgs, stats *store.Stats,
	pinl *model.Pinl, userId string, tagNames []string,
) (*model.Pinl, *model.Monl, error) {
	nurl, err := monlutils.NormalizeURL(pinl.URL)
	if err != nil {
		return nil, nil, err
	}
	found, err := monls.List(ctx, &store.MonlOpts{
		URL: nurl.String(),
	})
	if err != nil {
		return nil, nil, err
	}

	var monl *model.Monl
	if len(found) > 0 {
		monl = found[0]
	} else {
		monl = &model.Monl{URL: nurl.String()}
		if err := monls.Create(ctx, monl); err != nil {
			return nil, nil, err
		}
	}
	pinl.MonlID = monl.ID

	if out, _, err := storeutils.SavePinl(ctx, pinls, images, pinl, false); err == nil {
		pinl = out
	} else {
		return nil, nil, err
	}

	if _, err = storeutils.ReAssociateTags(ctx, tags, taggables, pinl, userId, tagNames); err != nil {
		return nil, nil, err
	}

	if out, err := storeutils.PinlWithLatestStats(ctx, pinls, monpkgs, stats, taggables, pinl.ID); err == nil {
		pinl = out
	} else {
		return nil, nil, err
	}

	return pinl, monl, nil
}

func (s *Server) pinlCreateHandler(w http.ResponseWriter, r *http.Request) {
	var in pinlBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}
	if err := validatePinlBody(in); err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		pinl   *model.Pinl
		monl   *model.Monl
		code   int
		outerr error
	)
	pinl = &model.Pinl{
		UserID:      user.ID,
		URL:         in.URL,
		Title:       in.Title,
		Description: in.Description,
	}

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		pinl2, monl2, err := savePinl(ctx, s.Pinls, s.Images, s.Tags, s.Taggables, s.Monls, s.Monpkgs, s.Stats, pinl, user.ID, in.Tags)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}
		pinl, monl = pinl2, monl2
		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}

	if s.ExchangeEnabled && monl.FetchedAt.Time().IsZero() {
		s.Queue.Add(job.NewFetchMonl(monl.ID))
	}
	s.Pubsub.Broadcast(message.NewPinlUpdated(pinl))
	response.JSON(w, pinl, http.StatusOK)
}

func (s *Server) pinlUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var in pinlBody
	err := request.JSON(r, &in)
	if err != nil {
		response.JSON(w, nil, http.StatusBadRequest)
		return
	}
	if err := validatePinlBody(in); err != nil {
		response.JSON(w, err, http.StatusBadRequest)
		return
	}

	var (
		ctx    = r.Context()
		user   = request.AuthedFrom(ctx)
		pinl   = request.PinlFrom(ctx)
		monl   *model.Monl
		code   int
		outerr error
	)
	pinl.MonlID = ""
	pinl.URL = in.URL
	pinl.Title = in.Title
	pinl.Description = in.Description

	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		var err error
		pinl2, monl2, err := savePinl(ctx, s.Pinls, s.Images, s.Tags, s.Taggables, s.Monls, s.Monpkgs, s.Stats, pinl, user.ID, in.Tags)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}
		pinl, monl = pinl2, monl2
		return true
	})

	if outerr != nil || response.IsError(code) {
		response.JSON(w, outerr, code)
		return
	}

	if s.ExchangeEnabled && monl.FetchedAt.Time().IsZero() {
		s.Queue.Add(job.NewFetchMonl(monl.ID))
	}
	s.Pubsub.Broadcast(message.NewPinlUpdated(pinl))
	response.JSON(w, pinl, http.StatusOK)
}

func (s *Server) pinlDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		pinl   = request.PinlFrom(ctx)
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		_, err := s.Taggables.DeleteByTarget(ctx, pinl)
		if err != nil {
			outerr, code = err, http.StatusInternalServerError
			return false
		}

		_, err = s.Pinls.Delete(ctx, pinl.ID)
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
	s.Pubsub.Broadcast(message.NewPinlDeleted(pinl))
	response.JSON(w, nil, http.StatusNoContent)
}

func (s *Server) pinlUploadImageHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		pinl   = request.PinlFrom(ctx)
		image  *model.Image
		code   int
		outerr error
	)
	s.Txer.TxFunc(ctx, func(ctx context.Context) bool {
		image2, code2, err := common.ImageUpload(ctx, r, s.Images, pinl, 1<<20, true)
		if err != nil {
			outerr, code = err, code2
			return false
		}
		image = image2

		pinl.ImageID = image.ID
		err = s.Pinls.Update(ctx, pinl)
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

	if out, err := storeutils.PinlWithLatestStats(ctx, s.Pinls, s.Monpkgs, s.Stats, s.Taggables, pinl.ID); err == nil {
		s.Pubsub.Broadcast(message.NewPinlUpdated(out))
	} else {
		response.JSON(w, err, http.StatusInternalServerError)
	}
	response.JSON(w, image, http.StatusOK)
}
