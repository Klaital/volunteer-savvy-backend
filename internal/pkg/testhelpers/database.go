package testhelpers

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CleanupTestDb truncates tables that are usually populated with fixtures.
// The set of tables is hardcoded.
func CleanupTestDb(db *sqlx.DB) error {
	// Delete all existing data
	tables := []string{
		"organizations",
		"users",
		"roles",
		"sites",
		"site_coordinators",
		"daily_schedules",
	}
	//sqlStmt := db.Rebind("DROP TABLE ?")
	for _, tableName := range tables {
		// TODO: handle this error more gracefully to discern between "Table not found" and other DB connection issues.
		//_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
		//if err != nil {
		//	return err
		//}
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	}

	return nil
}

// LoadFixtures scans the given directory for .sql files and runs all of them on the given DB handle.
// They will run in filesystem order, not sorted.
func LoadFixtures(db *sqlx.DB, fixturesDirectory string) error {
	// Search for any .sql files in the given testdata directory and run them
	files, err := ioutil.ReadDir(fixturesDirectory)
	if err != nil {
		wd, _ := os.Getwd()
		log.WithField("cwd", wd).WithError(err).Error("Failed to open fixtures directory")
		return err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlBytes, err := ioutil.ReadFile(filepath.Join(fixturesDirectory, f.Name()))
			if err != nil {

				return err
			}
			sqlStmt := db.Rebind(string(sqlBytes))
			_, err = db.Exec(sqlStmt)
			if err != nil {
				log.WithError(err).Error("failed to run fixture")
				return err
			}
		}
	}

	// Success!
	return nil
}

// DropAllTables enumerates all tables and drops them.
func DropAllTables(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	driver.Drop()
	return nil
}

// InitializeDatabase drops all tables in the provided database, then runs the migrations
// found in the migrationsDir, then loads the fixtures found in fixturesDir.
func InitializeDatabase(db *sqlx.DB, migrationsDir, fixturesDir string) error {
	dropDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.WithError(err).Error("Failed to create driver for dropping tables")
		return err
	}
	dropMigrator, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", dropDriver)
	if err != nil {
		wd, _ := os.Getwd()
		log.WithError(err).WithField("wd", wd).WithField("MigrationsDir", migrationsDir).Error("Failed to init drop migrator")
		return err
	}
	dropMigrator.Drop()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
	if err != nil {
		wd, _ := os.Getwd()
		log.WithError(err).WithField("wd", wd).WithField("MigrationsDir", migrationsDir).Error("Failed to init migrator")
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.WithError(err).Error("Failed to migrate schema")
		return err
	}
	err = LoadFixtures(db, fixturesDir)
	if err != nil {
		log.WithError(err).Error("Failed to load fixtures")
	}
	return err
}

// CountTable counts all rows in a table. Useful for testing whether a create or delete succeeded.
func CountTable(tableName string, db *sqlx.DB) int {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	row := db.QueryRow(sql)
	if row == nil {
		return 0
	}
	var count int
	row.Scan(&count)
	return count
}
