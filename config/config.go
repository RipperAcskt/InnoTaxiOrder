package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ELASTIC_DB_USERNAME string `mapstructure:"ELASTIC_DB_USERNAME"`
	ELASTIC_DB_PASSWORD string `mapstructure:"ELASTIC_DB_PASSWORD"`
	ELASTIC_DB_HOST     string `mapstructure:"ELASTIC_DB_HOST"`
	ELASTIC_DB_NAME     string `mapstructure:"ELASTIC_DB_NAME"`

	SERVER_HOST string `mapstructure:"SERVER_HOST"`

	GRPC_ANALYST_SERVICE_HOST string `mapstructure:"GRPC_ANALYST_SERVICE_HOST"`
	GRPC_DIVER_SERVICE_HOST   string `mapstructure:"GRPC_DIVER_SERVICE_HOST"`

	HS256_SECRET string `mapstructure:"HS256_SECRET"`

	SYNC_TIME int `mapstructure:"SYNC_TIME"`

	GRPC_HOST string `mapstructure:"GRPC_HOST"`

	ADMIN_LOGIN string `mapstructure:"ADMIN_LOGIN"`
	ADMIN_PASS  string `mapstructure:"ADMIN_PASS"`

	BROKER_HOST string `mapstructure:"BROKER_HOST"`
}

func New() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return config, nil
}
