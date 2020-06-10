package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/server"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	logger := log.WithFields(log.Fields{
		"package":   "volunteer-savvy-backend",
		"operation": "main",
	})

	log.SetLevel(log.DebugLevel)
	cfg, err := config.GetServiceConfig()
	if err != nil {
		logger.Fatalf("Unable to load service config: %v", err)
	} else {
		logger.Debugf("Loaded service config: %+v", cfg)
	}

	// Set up the database
	db := cfg.GetDbConn()
	for db == nil {
		log.Warnf("Waiting for database to come up")
		time.Sleep(2000 * time.Millisecond)
		db = cfg.GetDbConn()
	}
	defer db.Close()

	log.Debugf("Migrating database...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatalf("Failed to configure postgres driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/",
		"postgres",
		driver)
	if err != nil {
		logger.Fatalf("Failed to generate a migrator: %v", err)
	}
	err = m.Up()
	if err != nil {
		if err.Error() != "no change" {
			logger.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Initialize the server
	s, err := server.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to create server struct: %v", err)
	}
	// Actually start the application
	s.Serve()
}
