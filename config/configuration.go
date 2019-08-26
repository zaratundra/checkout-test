package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Server ServerConfig
	Data   DataConfig
}

type DataConfig struct {
	Products   string
	Promotions string
}

type ServerConfig struct {
	Port int
}

func LoadConfiguration(configPath, configFileName string) (Configuration, error) {
	var configuration Configuration

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)

	if err := viper.ReadInConfig(); err != nil {
		return Configuration{}, err
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
