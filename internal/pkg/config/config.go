package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type ServiceConfig struct {
	BasePath string `env:"BASE_PATH"`

	ServiceVersion string `env:"SERIVCE_VERSION" envDefault:"0.0.0"`

	DatabaseHost       string   `env:"DB_HOST"`
	DatabaseDriver     string   `env:"DB_DRIVER"`
	DatabaseUser       string   `env:"DB_USER"`
	DatabasePassword   string   `env:"PGPASSWORD"`
	DatabasePort       int64    `env:"DB_PORT"`
	DatabaseName       string   `env:"DB_NAME"` // the actual database name to connect to
	databaseConnection *sqlx.DB // to be set at runtime after main connects to the database
	FixturesPath       string   `env:"FIXTURES_PATH" envDefault:"testdata"` // only used in test
	MigrationsPath     string   `env:"MIGRATIONS_PATH" envDefault:"db/migrations"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	LogStyle string `env:"LOG_STYLE" envDefault:"prettyjson"`

	Port int64 `env:"PORT" envDefault:"8080"`

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

	// Oauth
	BcryptCost          int    `env:"BCRYPT_COST" envDefault:"5"`
	TokenExpirationTime string `env:"JWT_EXPIRATION_DURATION" envDefault:"4h"`
	JwtPrivateKey       string `env:"OAUTH_JWT_PRIVATE_KEY"`
	jwtPrivateKey       *rsa.PrivateKey
	JwtPublicKey        string `env:"OAUTH_JWT_PUBLIC_KEY"`
	jwtPublicKey        *rsa.PublicKey
}

// GetTokenExpirationDuration converts the JWT_EXPIRATION_DURATION environment
// variable to a time.Duration, substituting a safe default if the env var is
// malformed or missing.
func (cfg *ServiceConfig) GetTokenExpirationDuration() time.Duration {
	d, err := time.ParseDuration(cfg.TokenExpirationTime)
	if err == nil {
		// Default to 1 hour for JWT expiration
		return 1 * time.Hour
	}
	return d
}

var serviceConfig *ServiceConfig

func (cfg *ServiceConfig) getDatabaseDsn() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", cfg.DatabaseDriver, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
}
func GetServiceConfig() (config *ServiceConfig, err error) {
	if serviceConfig != nil {
		return serviceConfig, nil
	}

	config = &ServiceConfig{}
	err = env.Parse(config)
	if err != nil {
		return nil, err
	}

	// instantiate the logger with some default fields
	logger := logrus.New()
	// Configure the service logger's level
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logger.WithError(err).Errorf("Invalid log level '%s' specified. Defaulting to %s", config.LogLevel, logrus.DebugLevel)
		logLevel = logrus.DebugLevel
		logger.SetLevel(logLevel)
	} else {
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
		dbConn, err := sqlx.Connect(cfg.DatabaseDriver, cfg.getDatabaseDsn())
		if err != nil {
			cfg.Logger.WithField("driver", cfg.DatabaseDriver).WithError(err).Error("Failed to connect to database")
			return nil
		}
		cfg.databaseConnection = dbConn
	}

	return cfg.databaseConnection
}

func (cfg *ServiceConfig) GetPublicKey() *rsa.PublicKey {
	if cfg.jwtPublicKey == nil {
		cfg.GetJWTKeys()
	}
	return cfg.jwtPublicKey
}

func (cfg *ServiceConfig) GetJWTKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	if cfg.jwtPublicKey == nil {
		// Try to base64 decode it once, since that's how we have to handle it for local run and testing
		tmpKey, err := base64.StdEncoding.DecodeString(cfg.JwtPublicKey)
		if err != nil {
			tmpKey = []byte(cfg.JwtPublicKey)
		}
		block, _ := pem.Decode(tmpKey)
		if block == nil {
			cfg.Logger.WithField("JWTPublicKey", string(tmpKey)).Error("ssh: no public key found")
			return nil, nil
		}
		publicKeyTmp, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			cfg.Logger.WithError(err).Error("Failed to parse public key")
			return nil, nil
		}
		publicKey, ok := publicKeyTmp.(*rsa.PublicKey)
		if !ok {
			cfg.Logger.Errorf("Failed to cast public key to rsa pointer: %v", reflect.TypeOf(publicKeyTmp))
			return nil, nil
		}
		cfg.jwtPublicKey = publicKey
	}

	if cfg.jwtPrivateKey == nil {
		// Try to base64 decode it once, since that's how we have to handle it for local run and testing
		tmpKey, err := base64.StdEncoding.DecodeString(cfg.JwtPrivateKey)
		if err != nil {
			tmpKey = []byte(cfg.JwtPrivateKey)
		}
		block, _ := pem.Decode(tmpKey)
		if block == nil {
			cfg.Logger.Error("ssh: no private key found")
			return nil, nil
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			cfg.Logger.WithError(err).Error("Failed to parse private key")
			return nil, nil
		}
		cfg.jwtPrivateKey = privateKey
	}

	// Success!
	return cfg.jwtPrivateKey, cfg.jwtPublicKey
}
