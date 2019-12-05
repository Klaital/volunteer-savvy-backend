package helpers

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

func CleanupTestDb(db *sqlx.DB) error {
	// Delete all existing data
	tables := []string{
		"organizations",
	}
	//sqlStmt := db.Rebind("DROP TABLE ?")
	for _, tableName := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s ", tableName))
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadFixtures(db *sqlx.DB) error {
	// Search for any .sql files in the testdata directory, and run them
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		logrus.WithError(err).Fatal("failed to load fixture data")
		return err
	}

	for _, f := range files {
		sqlBytes, err := ioutil.ReadFile(filepath.Join("testdata", f.Name()))
		if err != nil {
			logrus.WithError(err).Fatal("Failed to read fixture file")
			return err
		}
		sqlStmt := db.Rebind(string(sqlBytes))
		_, err = db.Exec(sqlStmt)
	}

	return nil
}

func InitializeTestDb() (*sqlx.DB, error) {
	// Launch the test database

	// Connect to a local db
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "vstester", "vstester", "localhost", 5556, "vstest")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Run migrations to sync the schema
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "",
		DatabaseName:    "vstest",
	})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../db/migrations/",
		"vstest",
		driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return nil, err
		}
	}

	// TODO: load fixture data

	return db, nil
}

func CountTable(table string, db *sqlx.DB) int {
	var count int
	sqlStmt := db.Rebind("SELECT COUNT(*) FROM ?")
	row := db.QueryRowx(sqlStmt, table)
	err := row.Scan(&count)
	if err != nil {
		fmt.Printf("Failed to count table: %v", err)
		panic(err)
	}
	return count
}