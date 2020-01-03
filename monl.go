package main

import (
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/monl"
	"github.com/pinmonl/pinmonl/monl/github"
)

func initMonl(cfg *config.Config, ss stores) *monl.Monl {
	ml := monl.New(monl.Opts{})
	ml.Register(newGithubVendor(cfg))

	return ml
}

func newGithubVendor(cfg *config.Config) monl.Vendor {
	return github.NewVendor(github.Opts{
		Token: cfg.Github.Token,
	})
}
