package main

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger := log.WithFields(log.Fields{
		"package":   "volunteer-savvy-backend",
		"operation": "main",
	})

	logger.Info("Starting server...")

	cfg, err := config.GetServiceConfig()
	if err != nil {
		logger.Fatalf("Unable to load service config: %v", err)
	} else {
		logger.Debugf("Loaded service config: %v", cfg)
	}

	log.Debugf("Preparing DB connection with %s", cfg.DatabaseDriver)
	db, err := sqlx.Connect(cfg.DatabaseDriver, cfg.DatabaseDSN)
	for err != nil {
		log.Warnf("Waiting for database to come up: %v", err)
		time.Sleep(1000 * time.Millisecond)
		db, err = sqlx.Connect(cfg.DatabaseDriver, cfg.DatabaseDSN)
	}
	defer db.Close()
}
