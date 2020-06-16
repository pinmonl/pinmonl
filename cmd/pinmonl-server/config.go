package main

import (
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Address string
	Verbose int

	JWT struct {
		Secret string
		Issuer string
		Expire time.Duration
	}

	DB struct {
		Driver string
		DSN    string
	}

	Github struct {
		Tokens []string
	}

	Queue struct {
		Job    int
		Worker int
	}
}

func readConfig() (*config, error) {
	var c config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
