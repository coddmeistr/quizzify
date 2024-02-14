package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var connectionString, migrationsPath, migrationsTable string

	flag.StringVar(&connectionString, "postgres-url", "", "connection string")
	flag.StringVar(&migrationsPath, "migrations-path", "", "where migrations files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of the migrations table")
	flag.Parse()

	if connectionString == "" {
		panic("postgres-url is empty")
	}

	if migrationsPath == "" {
		panic("migrations-path is empty")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: migrationsTable, DatabaseName: "sso"})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	//m.Force(1)
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations were applied sucesessfully")
}
