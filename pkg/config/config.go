package config

import "github.com/spf13/viper"

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Paswword string `mapstructure:"password"`
}

type Port struct {
	SvcPort string `mapstructure:"port"`
}

type Config struct {
	Postgres DBConfig `mapstructure:"db"`
	Port     Port     `mapstructure:"svc-port"`
}

var config Config

func LoadConfig() (Config, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("pkg/config/")

	if err := vp.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := vp.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

