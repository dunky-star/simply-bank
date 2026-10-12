package util

import (
	"github.com/spf13/viper"
)

// COnfig store all configuration of the application
// The values are read by viper from a config file or environment variables

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".app.env")
	viper.SetConfigType("env") // We can also use json, yaml, etc.

	viper.AutomaticEnv()
	viper.BindEnv("DB_DRIVER", "DB_DRIVER")
	viper.BindEnv("DB_SOURCE", "DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS", "SERVER_ADDRESS")

	// Config file is optional; env vars (AutomaticEnv) suffice in CI/prod.
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}
