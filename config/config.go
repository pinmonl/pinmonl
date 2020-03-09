package config

import (
	"strings"

	"github.com/pinmonl/pinmonl/pkg/generate"
	"github.com/spf13/viper"
)

// Config stores configuration of the server.
type Config struct {
	LogLevel string

	DB struct {
		Driver string
		DSN    string
	}
	Github struct {
		Token string
	}
	HTTP struct {
		Endpoint  string
		DevServer string
	}
	Queue struct {
		Parallel int
		Interval string
	}
	Cookie struct {
		Name     string
		MaxAge   string
		Secure   bool
		HashKey  string
		BlockKey string
	}
}

// Read returns configuration from config file and env.
func Read() *Config {
	var c Config
	v := newViper()
	v.Unmarshal(&c)
	return &c
}

func newViper() *viper.Viper {
	v := viper.NewWithOptions()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.ReadInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("db.driver", "sqlite3")
	v.SetDefault("db.dsn", "file:pinmonl.db?cache=shared")
	v.SetDefault("github.token", "")
	v.SetDefault("http.endpoint", ":8080")
	v.SetDefault("http.devserver", "")
	v.SetDefault("loglevel", "info")
	v.SetDefault("queue.interval", "1s")
	v.SetDefault("queue.parallel", 1)
	v.SetDefault("cookie.name", "session")
	v.SetDefault("cookie.maxage", "360h")
	v.SetDefault("cookie.secure", false)
	v.SetDefault("cookie.hashkey", generate.RandomString(32))
	v.SetDefault("cookie.blockkey", generate.RandomString(32))

	return v
}
