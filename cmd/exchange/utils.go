package main

import (
	"github.com/pinmonl/pinmonl/pkgs/generate"
	"github.com/sirupsen/logrus"
)

func catchErr(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

func generateKey() []byte {
	key := generate.AlphaNum(128)
	return []byte(key)
}
