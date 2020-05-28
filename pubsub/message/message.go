package message

import (
	"context"

	"github.com/pinmonl/pinmonl/handler/api/apibody"
	"github.com/pinmonl/pinmonl/handler/api/apiutils"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/store"
)

func newPinlMessage(pinl apibody.Pinl, topic string) *pubsub.Message {
	return pubsub.NewMessage(pinl.UserID, topic, pinl)
}

func NewPinlCreateMessage(pinl apibody.Pinl) *pubsub.Message {
	return newPinlMessage(pinl, "pinl.create")
}

func NewPinlUpdateMessage(pinl apibody.Pinl) *pubsub.Message {
	return newPinlMessage(pinl, "pinl.update")
}

func NewPinlDeleteMessage(pinl apibody.Pinl) *pubsub.Message {
	return newPinlMessage(pinl, "pinl.delete")
}

func NotifyPkgUser(
	ctx context.Context,
	ws *pubsub.Server,
	pinlStore store.PinlStore,
	monpkgStore store.MonpkgStore,
	taggableStore store.TaggableStore,
	statStore store.StatStore,
	pkg model.Pkg,
) error {
	monlMap, err := monpkgStore.ListMonls(ctx, &store.MonpkgOpts{
		PkgIDs: []string{pkg.ID},
	})
	if err != nil {
		return err
	}

	pinls, err := pinlStore.List(ctx, &store.PinlOpts{
		MonlIDs: model.MonlList(monlMap[pkg.ID]).Keys(),
	})
	if err != nil {
		return err
	}

	return NotifyPinlUser(ctx, ws, monpkgStore, taggableStore, statStore, pinls...)
}

func NotifyPinlUser(
	ctx context.Context,
	ws *pubsub.Server,
	monpkgStore store.MonpkgStore,
	taggableStore store.TaggableStore,
	statStore store.StatStore,
	pinls ...model.Pinl,
) error {
	if len(pinls) == 0 {
		return nil
	}

	pkgMap, statMap, err := apiutils.ListPinlStats(ctx, monpkgStore, statStore, pinls...)
	if err != nil {
		return err
	}

	tagMap, err := taggableStore.ListTags(ctx, &store.TaggableOpts{
		Targets: model.MustBeMorphables(pinls),
	})
	if err != nil {
		return err
	}

	for _, p := range pinls {
		body := apibody.NewPinl(p).
			WithTags(tagMap[p.ID]).
			WithPkgs(pkgMap[p.MonlID], statMap)
		msg := NewPinlUpdateMessage(body)
		ws.Publish(msg)
	}

	return nil
}
