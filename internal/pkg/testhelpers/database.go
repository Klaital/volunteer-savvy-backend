package testhelpers

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// LoadFixtures scans the given directory for .sql files and runs all of them on the given DB handle.
// They will run in filesystem order, not sorted.
func LoadFixtures(db *sqlx.DB, fixturesDirectory string) error {
	// Search for any .sql files in the given testdata directory and run them
	files, err := ioutil.ReadDir(fixturesDirectory)
	if err != nil {
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
				return err
			}
		}
	}

	// Success!
	return nil
}

// DropAllTables enumerates all tables and drops them.
func DropAllTables(db *sqlx.DB) error {
	tables := make([]string, 0)
	sqlStmt := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'`
	err := db.Select(&tables, sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = db.Rebind(`DROP TABLE IF EXISTS ? CASCADE`)
	for _, tableName := range tables {
		_, err = db.Exec(sqlStmt, tableName)
		if err != nil {
			return err
		}
	}

	// Success!
	return nil
}

// ResetFixtures truncates all tables, then reruns all .sql files in the fixtures directory.
func ResetFixtures(db *sqlx.DB, fixturesDir string) error {
	// Scan for all table names
	tables := make([]string, 0)
	sqlStmt := `SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'`
	err := db.Select(&tables, sqlStmt)
	if err != nil {
		return err
	}

	// Truncate each of them
	sqlStmt = db.Rebind(`TRUNCATE ? CASCADE`)
	for _, tableName := range tables {
		_, err = db.Exec(sqlStmt, tableName)
		if err != nil {
			return err
		}
	}

	// Reload the fixture data
	return LoadFixtures(db, fixturesDir)
}

// InitializeDatabase drops all tables in the provided database, then runs the migrations
// found in the migrationsDir, then loads the fixtures found in fixturesDir.
func InitializeDatabase(db *sqlx.DB, migrationsDir, fixturesDir string) error {
	_ = DropAllTables(db)
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsDir, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	err = LoadFixtures(db, fixturesDir)
	return err
}
