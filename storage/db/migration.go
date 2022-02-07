package db

import (
	//"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const _upMigrationCounter = 1
func UpPostgresMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	println("path - ", os.Getenv("PATH_MIGRATIONS"))
	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("PATH_MIGRATIONS"),
		_postgresDriverName, driver,
	)
	if err != nil {
		return err
	}

	m.Steps(_upMigrationCounter)
	return nil
}
 
const _downMigrationCounter = -1
func DownPostgresMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	println("path - ", os.Getenv("PATH_MIGRATIONS"))
	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("PATH_MIGRATIONS"),
		_postgresDriverName, driver,
	)
	if err != nil {
		return err
	}

	m.Steps(_downMigrationCounter)
	return nil
} 

func UpdatePostgresDbScheme(db *sqlx.DB) error {
	if err := DownPostgresMigrations(db); err != nil {
		return err
	} 

	if err := UpPostgresMigrations(db); err != nil {
		return err
	} 

	return nil
}