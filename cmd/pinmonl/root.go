package main

import (
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/cmd/pinmonl/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config", "c", "", "path to config file")
	pflags.IntP("v", "v", 0, "log level verbosity")

	viper.BindPFlag("verbose", pflags.Lookup("v"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/pinmonl")
		viper.SetConfigName("client")
	}

	viper.SetEnvPrefix("PINMONL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("address", ":3399")
	viper.SetDefault("db.driver", "sqlite3")
	viper.SetDefault("db.dsn", "client.db")
	viper.SetDefault("exchange.address", "https://pinmonl.com")
	viper.SetDefault("exchange.enabled", true)
	viper.SetDefault("jwt.expire", "24h")
	viper.SetDefault("jwt.issuer", "pinmonl")
	viper.SetDefault("jwt.secret", string(generateKey()))
	viper.SetDefault("queue.job", 1)
	viper.SetDefault("queue.worker", 1)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var rootCmd = &cobra.Command{
	Use:     "pinmonl",
	Short:   "Pinmonl bookmark client",
	Long:    `Pinmonl bookmark client`,
	Version: version.Version.String(),
}
