package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	"github.com/markbates/pkger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pinmonl/pinmonl/cmd/pinmonl-server/config"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/handler/serverapi"
	"github.com/pinmonl/pinmonl/store"
)

func main() {
	pkger.Include("/migrations/server")
	c := config.Read()
	db, _ := database.NewDB(&database.DBOpts{
		Driver: c.DB.Driver,
		DSN:    c.DB.DSN,
	})
	fmt.Println(db.Migrate.Down())
	fmt.Println(db.Migrate.Up())

	s := newStores(db)
	r := newRouter(s)
	http.ListenAndServe(c.Address, r)
}

type stores struct {
	store *store.Store
	users *store.Users
}

func newStores(db *database.DB) stores {
	s := store.NewStore(db)
	return stores{
		store: s,
		users: store.NewUsers(s),
	}
}

func newRouter(s stores) http.Handler {
	r := chi.NewRouter()
	serverapi := serverapi.NewServer(&serverapi.ServerOpts{
		Users: s.users,
	})
	r.Mount("/api", serverapi.Handler())
	return r
}
