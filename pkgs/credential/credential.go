package credential

import "github.com/spf13/viper"

type Store struct {
	*viper.Viper
	creds Credential
}

type Credential struct {
	User struct {
		ID          string
		AccessToken string
	}

	Machine struct {
		ID          string
		AccessToken string
	}
}

func NewStore() (*Store, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".credential")

	v.SetEnvPrefix("PMC")

	return &Store{Viper: v}, nil
}

func (s *Store) Read() error {
	return s.Viper.Unmarshal(&s.creds)
}

func (s *Store) Write() error {
	return s.Viper.WriteConfig()
}
