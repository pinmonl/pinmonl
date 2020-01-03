package config

import "github.com/spf13/viper"

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
		Endpoint string
	}
	Queue struct {
		Parallel int
		Interval string
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

	v.AutomaticEnv()

	v.SetDefault("db.driver", "sqlite3")
	v.SetDefault("db.dsn", "file:pinmonl.db?cache=shared")
	v.SetDefault("github.token", "")
	v.SetDefault("http.endpoint", ":8080")
	v.SetDefault("loglevel", "debug")
	v.SetDefault("queue.interval", "1s")
	v.SetDefault("queue.parallel", 1)

	return v
}
