package main

import (
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/store"
)

type stores struct {
	store     store.Store
	images    store.ImageStore
	jobs      store.JobStore
	monls     store.MonlStore
	pinls     store.PinlStore
	pinmons   store.PinmonStore
	shares    store.ShareStore
	sharetags store.ShareTagStore
	stats     store.StatStore
	taggables store.TaggableStore
	tags      store.TagStore
	users     store.UserStore
}

func initStores(db *database.DB) stores {
	s := store.NewStore(db)
	return stores{
		store:     s,
		images:    store.NewImageStore(s),
		jobs:      store.NewJobStore(s),
		monls:     store.NewMonlStore(s),
		pinls:     store.NewPinlStore(s),
		pinmons:   store.NewPinmonStore(s),
		shares:    store.NewShareStore(s),
		sharetags: store.NewShareTagStore(s),
		stats:     store.NewStatStore(s),
		taggables: store.NewTaggableStore(s),
		tags:      store.NewTagStore(s),
		users:     store.NewUserStore(s),
	}
}