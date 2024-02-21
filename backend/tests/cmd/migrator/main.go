package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func main() {
	var option string
	flag.StringVar(&option, "option", "", "option for migrator")

	var connectionString, migrationsPath, migrationsTable string

	if err := godotenv.Load(); err != nil {
		fmt.Println("couldn't load envs from .env file")
	}

	connectionString = os.Getenv("MONGO_CONNECTION_URI")
	flag.StringVar(&migrationsPath, "migrations-path", "", "where migrations files")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of the migrations table")
	flag.Parse()

	if connectionString == "" {
		panic("connection string is empty")
	}

	if migrationsPath == "" {
		panic("migrations-path is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(fmt.Errorf("failed to connect to mongoapp: %w", err))
	}
	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: "quizzify-tests", MigrationsCollection: migrationsTable})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mongodb", driver)
	if err != nil {
		panic(fmt.Errorf("failed to create migrator: %w", err))
	}

	switch option {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no up migrations to apply")
				return
			}
			panic(fmt.Errorf("failed to perform up migrations: %w", err))
		}
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no down migrations to apply")
				return
			}
			panic(fmt.Errorf("failed to perform down migrations: %w", err))
		}
	default:
		panic("unknown option")
	}

	fmt.Println(fmt.Sprintf("%s migrations were applied successfully", option))
}
