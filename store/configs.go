package store

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Configs struct {
	*sync.Mutex
	viper      *viper.Viper
	configFile string
	envPrefix  string
}

func NewConfigs() *Configs {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/pinmonl")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.BindEnv("machine.expireAt")
	v.BindEnv("machine.token")
	v.BindEnv("user.expireAt")
	v.BindEnv("user.login")
	v.BindEnv("user.token")
	v.BindEnv("user.linked.userId")
	v.BindEnv("user.defaultUserId")

	c := &Configs{
		Mutex: &sync.Mutex{},
		viper: v,
	}
	c.SetConfigName("pmdata")
	c.SetEnvPrefix("PMDATA")

	return c
}

func (c *Configs) SetConfigName(name string) {
	c.configFile = name
	c.viper.SetConfigName(name)
	c.viper.ReadInConfig()
}

func (c *Configs) SetEnvPrefix(prefix string) {
	c.envPrefix = prefix
	c.viper.SetEnvPrefix(prefix)
}

func (c *Configs) Save() error {
	c.Lock()
	defer c.Unlock()

	if err := c.viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return c.viper.WriteConfigAs(c.configFile + ".yaml")
}

func (c *Configs) GetMachineToken() string {
	return c.viper.GetString("machine.token")
}

func (c *Configs) GetMachineExpireAt() time.Time {
	return c.viper.GetTime("machine.expireAt")
}

func (c *Configs) GetUserToken() string {
	return c.viper.GetString("user.token")
}

func (c *Configs) GetUserLogin() string {
	return c.viper.GetString("user.login")
}

func (c *Configs) GetUserExpireAt() time.Time {
	return c.viper.GetTime("user.expireAt")
}

func (c *Configs) GetUserLinkedUserID() string {
	return c.viper.GetString("user.linked.userId")
}

func (c *Configs) GetUserDefaultUserID() string {
	return c.viper.GetString("user.defaultUserId")
}

func (c *Configs) SetMachineToken(token string) {
	c.viper.Set("machine.token", token)
}

func (c *Configs) SetMachineExpireAt(t time.Time) {
	c.viper.Set("machine.expireAt", t)
}

func (c *Configs) SetUserToken(token string) {
	c.viper.Set("user.token", token)
}

func (c *Configs) SetUserLogin(login string) {
	c.viper.Set("user.login", login)
}

func (c *Configs) SetUserExpireAt(t time.Time) {
	c.viper.Set("user.expireAt", t)
}

func (c *Configs) SetUserLinkedUserID(userID string) {
	c.viper.Set("user.linked.userId", userID)
}

func (c *Configs) SetUserDefaultUserID(userID string) {
	c.viper.Set("user.defaultUserId", userID)
}
