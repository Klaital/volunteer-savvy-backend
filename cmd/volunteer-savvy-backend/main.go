package main

import (
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	log "github.com/sirupsen/logrus"
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

}
