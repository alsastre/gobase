package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// List of viper keys, env variables and default values
const (
	LogLevelKey     = "loglevel"
	LogLevelEnv     = "SERVER_LOG_LEVEL"
	LogLevelDefault = "debug"
)

func initializeConfig() (*viper.Viper, error) {
	viperCfg := viper.New()
	viperCfg.AddConfigPath(".")
	viperCfg.SetConfigName("config")

	viperCfg.SetDefault(LogLevelKey, LogLevelDefault)
	viperCfg.BindEnv(LogLevelKey, LogLevelEnv)

	viperCfg.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viperCfg.ReadInConfig(); err != nil {
		fmt.Println(err.Error())
	}
	// We return the reference to the viper config
	// We will use the viper API to retrieve settings instead of loading them into a struct
	return viperCfg, nil
}

// New creates a new configuration
func New() (*viper.Viper, error) {
	return initializeConfig()
}
