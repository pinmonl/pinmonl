package main

import (
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/cmd/exchange/version"
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
		viper.SetConfigName("exchange")
	}

	viper.SetEnvPrefix("PINMONL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("address", ":8080")
	viper.SetDefault("db.driver", "postgres")
	viper.SetDefault("db.dsn", "postgres://pinmonl:pinmonl@pg:5432/pinmonl?sslmode=disable")
	viper.SetDefault("git.dev", false)
	viper.SetDefault("github.tokens", []string{})
	viper.SetDefault("youtube.tokens", []string{})
	viper.SetDefault("jwt.expire", "168h")
	viper.SetDefault("jwt.issuer", "pinmonl-exchange")
	viper.SetDefault("jwt.secret", string(generateKey()))
	viper.SetDefault("queue.job", 1)
	viper.SetDefault("queue.worker", 1)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var rootCmd = &cobra.Command{
	Use:     "pinmonl-exchange",
	Short:   "Pinmonl exchange server",
	Long:    `Pinmonl exchange server handles resources heavy tasks and hosts shared bookmarks.`,
	Version: version.Version.String(),
}
