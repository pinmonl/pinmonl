package config

import (
	"strings"

	"github.com/pinmonl/pinmonl/pkg/generate"
	"github.com/spf13/viper"
)

// Config stores configuration of the server.
type Config struct {
	LogLevel   string
	SingleUser bool

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
	Client struct {
		Host  string
		Token string
	}
	Oauth struct {
		PrivateKey     string
		PrivateKeyFile string
	}
}

// Read returns configuration from config file and env.
func Read() *Config {
	r := NewReader()
	return r.Config()
}

// Reader reads config values from environment and files.
type Reader struct {
	*viper.Viper
}

// NewReader creates config reader.
func NewReader() *Reader {
	v := viper.NewWithOptions()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.ReadInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("singleuser", true)
	v.SetDefault("db.driver", "sqlite3")
	v.SetDefault("db.dsn", "file:pinmonl.db?cache=shared")
	v.SetDefault("github.token", "")
	v.SetDefault("http.endpoint", ":3399")
	v.SetDefault("http.devserver", "")
	v.SetDefault("loglevel", "info")
	v.SetDefault("queue.interval", "1s")
	v.SetDefault("queue.parallel", 1)
	v.SetDefault("cookie.name", "session")
	v.SetDefault("cookie.maxage", "360h")
	v.SetDefault("cookie.secure", false)
	v.SetDefault("cookie.hashkey", generate.RandomString(32))
	v.SetDefault("cookie.blockkey", generate.RandomString(32))
	v.SetDefault("client.host", "http://localhost:3399")
	v.SetDefault("oauth.privatekey", "")
	v.SetDefault("oauth.privatekeyfile", "oauth.key")

	return &Reader{Viper: v}
}

// Config parses values into Config.
func (r *Reader) Config() *Config {
	var c Config
	r.Unmarshal(&c)
	return &c
}
