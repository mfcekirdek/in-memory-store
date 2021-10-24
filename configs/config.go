// Package configs contains the application configurations.
package configs

import (
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

// Config struct defines the application's configuration variables.
// IsDebug is used for enabling HTTP logger.
// AppName is the application name.
// SaveToFileInterval defines the time period (as minutes) during which the store content will be written to a json file.
// StorageDirPath is the folder where the TIMESTAMP-data.json files are located.
// Server includes the Port variable to set which port the application will run on.
type Config struct {
	IsDebug            bool
	AppName            string
	SaveToFileInterval int
	StorageDirPath     string
	Server             ServerConfig
}

// Creates a Config instance by loading values from environment variables.
func NewConfig() *Config {
	saveToFileInterval := 10 // in minutes // added this variable to prevent 'magic number detected' lint error (gomnd)
	defaultPort := 8080      // added this variable to prevent 'magic number detected' lint error (gomnd)
	viper.SetDefault("IS_DEBUG", true)
	viper.SetDefault("APP_NAME", "kv-store")
	viper.SetDefault("STORAGE_DIR_PATH", "storage")
	viper.SetDefault("SAVE_TO_FILE_INTERVAL", saveToFileInterval)
	viper.SetDefault("PORT", defaultPort)
	viper.AutomaticEnv()

	config := &Config{
		IsDebug:            viper.GetBool("IS_DEBUG"),
		AppName:            viper.GetString("APP_NAME"),
		StorageDirPath:     viper.GetString("STORAGE_DIR_PATH"),
		SaveToFileInterval: viper.GetInt("SAVE_TO_FILE_INTERVAL"),
		Server: ServerConfig{
			Port: viper.GetInt("PORT"),
		},
	}
	return config
}

// Prints all config variables in a good format.
func (c *Config) Print() {
	_, _ = pp.Println(c)
}
