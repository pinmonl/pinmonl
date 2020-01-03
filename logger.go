package main

import (
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/sirupsen/logrus"
)

func initLogger(cfg *config.Config) {
	lvl, _ := logrus.ParseLevel(cfg.LogLevel)
	logx.SetLevel(lvl)
}
