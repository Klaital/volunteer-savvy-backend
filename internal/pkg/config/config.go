package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
)

type ServiceConfig struct {
	BasePath string `env:"BASE_PATH"`

	DatabaseHost     string `env:"DB_HOST"`
	DatabaseDriver   string `env:"DB_DRIVER"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"PGPASSWORD"`
	DatabasePort     int64  `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"` // the actual database name to connect to
	DatabaseDSN      string // to be constructed after parsing the env variables
	DatabaseConnection *sqlx.DB // to be set at runtime after main connects to the database

	Debug            bool   `env:"DEBUG" envDefault:"false"`

	Port             int64  `env:"PORT" envDefault:"8080"`

	// Swagger
	SwaggerFilePath string `env:"SWAGGER_FILE_PATH"`
	APIPath         string `env:"SWAGGER_API_PATH"`
	SwaggerPath     string `env:"SWAGGER_PATH"`

	// HealthCheck
	HealthCheckPath string `env:"HEALTH_PATH" envDefault:"/GetServiceStatus"`
}

var serviceConfig *ServiceConfig


func GetServiceConfig() (config *ServiceConfig, err error) {
	if serviceConfig != nil {
		return serviceConfig, nil
	}

	config = &ServiceConfig{}
	err = env.Parse(config)
	if err != nil {
		return nil, err
	}

	// manually construct the DSN
	config.DatabaseDSN = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", config.DatabaseDriver, config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)

	serviceConfig = config
	return config, nil
}
