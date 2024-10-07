package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	SERVER_ADDRESS      string `mapstructure:"SERVER_ADDRESS"`
	OPENSEARCH_HOST     string `mapstructure:"OPENSEARCH_HOST"`
	OPENSEARCH_PORT     string `mapstructure:"OPENSEARCH_PORT"`
	OPENSEARCH_USERNAME string `mapstructure:"OPENSEARCH_USERNAME"`
	OPENSEARCH_PASSWORD string `mapstructure:"OPENSEARCH_PASSWORD"`
}

// LoadConfig parses env variables stored in .env file and returns them
//
// Parameters:
//   - path: a string representing .env file relative location
//
// Returns:
//   - config: a interface{} which gives variables stored in .env
//   - err: a error in case parsing environment variables failed
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
