package main

import (
	"time"

	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/session"
)

type sessions struct {
	cookie *session.CookieStore
}

func initSessionStore(cfg *config.Config) sessions {
	s := sessions{
		cookie: initSessionCookieStore(cfg),
	}
	return s
}

func initSessionCookieStore(cfg *config.Config) *session.CookieStore {
	age, err := time.ParseDuration(cfg.Cookie.MaxAge)
	if err != nil {
		logx.Panicln("Config.Cookie.MaxAge parse error", err)
	}

	return session.NewCookieStore(session.CookieStoreOpts{
		Name:     cfg.Cookie.Name,
		MaxAge:   age,
		Secure:   cfg.Cookie.Secure,
		HashKey:  []byte(cfg.Cookie.HashKey),
		BlockKey: []byte(cfg.Cookie.BlockKey),
	})
}
