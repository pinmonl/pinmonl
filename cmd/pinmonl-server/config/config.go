package config

import (
	"strings"

	"github.com/pinmonl/pinmonl/pkgs/cfgio"
)

var cfg *cfgio.IO

func init() {
	cfg = cfgio.New()
	cfg.SetConfigName("server")
	cfg.AddConfigPath("/etc/pinmonl")
	cfg.AddConfigPath("/pinmonl")
	cfg.ReadInConfig()

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()

	cfg.SetDefault("address", ":8080")
	cfg.SetDefault("db.driver", "sqlite3")
	cfg.SetDefault("db.dsn", "server.db")
}

type Config struct {
	Address string

	DB struct {
		Driver string
		DSN    string
	}

	Github struct {
		Tokens []string
	}
}

func Read() *Config {
	var c Config
	cfg.Unmarshal(&c)
	return &c
}
