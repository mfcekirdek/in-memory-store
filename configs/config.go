package configs

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

type Config struct {
	IsDebug bool
	AppName string
	Server  ServerConfig
}

func NewConfig() *Config {
	defaultPort := 8080 // added this variable to prevent 'magic number detected' lint error (gomnd)
	viper.SetDefault("IS_DEBUG", true)
	viper.SetDefault("APP_NAME", "kv-store")
	viper.SetDefault("PORT", defaultPort)
	viper.AutomaticEnv()

	config := &Config{
		IsDebug: viper.GetBool("IS_DEBUG"),
		AppName: viper.GetString("APP_NAME"),
		Server: ServerConfig{
			Port: viper.GetInt("PORT"),
		},
	}
	return config
}

func (c *Config) Print() {
	_, _ = pp.Println(c)
}
