package configs

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

type Config struct {
	IsDebug        bool
	AppName        string
	FlushInterval  int
	StorageDirPath string
	Server         ServerConfig
}

func NewConfig() *Config {
	flushInterval := 3  // in minutes // added this variable to prevent 'magic number detected' lint error (gomnd)
	defaultPort := 8080 // added this variable to prevent 'magic number detected' lint error (gomnd)
	viper.SetDefault("IS_DEBUG", true)
	viper.SetDefault("APP_NAME", "kv-store")
	viper.SetDefault("STORAGE_DIR_PATH", "storage")
	viper.SetDefault("FLUSH_INTERVAL", flushInterval)
	viper.SetDefault("PORT", defaultPort)
	viper.AutomaticEnv()

	config := &Config{
		IsDebug:        viper.GetBool("IS_DEBUG"),
		AppName:        viper.GetString("APP_NAME"),
		StorageDirPath: viper.GetString("STORAGE_DIR_PATH"),
		FlushInterval:  viper.GetInt("FLUSH_INTERVAL"),
		Server: ServerConfig{
			Port: viper.GetInt("PORT"),
		},
	}
	return config
}

func (c *Config) Print() {
	_, _ = pp.Println(c)
}
