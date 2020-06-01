package cfgio

import "github.com/spf13/viper"

type IO struct {
	*viper.Viper
}

func New() *IO {
	v := viper.NewWithOptions()
	return &IO{Viper: v}
}
