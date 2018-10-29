package config

import (
	"github.com/caarlos0/env"
)

type ServiceConfig struct {
	BasePath string `env:"BASE_PATH"`

	DatabaseHost string `env:"DB_HOST"`
	DatabaseDriver string `env:"DB_DRIVER"`
	DatabaseUser string `env:"DB_USER"`
	DatabasePassword string `env:"PGPASSWORD"`
	DatabaseDSN string // to be constructed after parsing the env variables
}

var serviceConfig ServiceConfig

func GetServiceConfig() (config *ServiceConfig, err error) {
	config = &ServiceConfig{}
	err = env.Parse(config)
	if err != nil {
		return nil, err
	}

	serviceConfig = *config
	return config, nil
}
