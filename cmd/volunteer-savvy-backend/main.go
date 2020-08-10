package main

import (
	"github.com/emicklei/go-restful"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/config"
	oServer "github.com/klaital/volunteer-savvy-backend/internal/pkg/organizations/server"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/server"
	sServer "github.com/klaital/volunteer-savvy-backend/internal/pkg/sites/server"
	uServer "github.com/klaital/volunteer-savvy-backend/internal/pkg/users/server"
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
	orgServer := oServer.New(cfg)
	sitesServer := sServer.New(cfg)
	authServer := uServer.New(cfg)
	usersServer := uServer.New(cfg)

	services := []*restful.WebService{
		orgServer.GetOrganizationsAPI(),
		sitesServer.GetSitesAPI(),
		authServer.GetAuthAPI(),
		usersServer.GetUsersAPI(),
	}
	s, err := server.New(cfg, services)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create server struct")
	}
	// Actually start the application
	s.Serve()
}
