package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ServiceConfig struct {
	BasePath string `env:"BASE_PATH"`

	ServiceVersion string `env:"SERIVCE_VERSION" envDefault:"0.0.0"`

	DatabaseHost     string `env:"DB_HOST"`
	DatabaseDriver   string `env:"DB_DRIVER"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"PGPASSWORD"`
	DatabasePort     int64  `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"` // the actual database name to connect to
	databaseDSN      string // to be constructed after parsing the env variables
	databaseConnection *sqlx.DB // to be set at runtime after main connects to the database

	LogLevel         string   `env:"LOG_LEVEL" envDefault:"debug"`
	LogStyle         string `env:"LOG_STYLE" envDefault:"prettyjson"`

	Port             int64  `env:"PORT" envDefault:"8080"`

	// Configure the static fileserver
	StaticContentPath string `env:"STATIC_CONTENT_PATH"`

	// Swagger
	SwaggerFilePath string `env:"SWAGGER_FILE_PATH"`
	APIPath         string `env:"SWAGGER_API_PATH"`
	SwaggerPath     string `env:"SWAGGER_PATH"`

	// HealthCheck
	HealthCheckPath string `env:"HEALTH_PATH" envDefault:"/GetServiceStatus"`

	// Singleton logrus instance
	Logger *logrus.Entry
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
	config.databaseDSN = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", config.DatabaseDriver, config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)

	// instantiate the logger with some default fields
	logger := logrus.New()
	// Configure the service logger's level
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logger.SetLevel(logLevel)
	} else {
		logger.WithError(err).Errorf("Invalid log level '%s' specified. Defaulting to %s", config.LogLevel, logrus.DebugLevel)
		logLevel = logrus.DebugLevel
		logger.SetLevel(logLevel)
	}
	// Configure the service logger's formatter
	switch config.LogStyle {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	case "prettyjson":
		logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	// Configure the service logger to include the calling func
	// as a field in debug mode only
	logger.SetReportCaller(logger.GetLevel() == logrus.DebugLevel)

	// Set the configured logger with some default global fields
	config.Logger = logger.WithFields(logrus.Fields{
		"Service": "volunteer-savvy-backend",
		"Version": config.ServiceVersion,
	})

	serviceConfig = config
	return config, nil
}

func (cfg *ServiceConfig) GetDbConn() *sqlx.DB {
	// Construct the singleton db connection pool if needed
	if cfg.databaseConnection == nil {
		dbConn, err := sqlx.Connect(cfg.DatabaseDriver, cfg.databaseDSN)
		if err != nil {
			cfg.Logger.WithError(err).Error("Failed to connect to database")
			return nil
		}
		cfg.databaseConnection = dbConn
	}

	return cfg.databaseConnection
}