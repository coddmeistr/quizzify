package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/app"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/config"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("COULDN'T LOAD ENVS FROM .ENV FILE")
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Doing migrations")
	output, err := exec.Command("/sso/migrator",
		fmt.Sprintf("--postgres-url=%s", cfg.PostgresUrl),
		fmt.Sprintf("--migrations-path=%s", "/sso/migrations")).
		Output()
	if err != nil {
		log.Error("Failed migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("Migrations done", slog.String("output", string(output)))

	log.Info("Starting application", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.PostgresUrl, cfg.TokenTTL)
	go application.GRPCApp.MustRun()

	// Init services

	// Start tests-server

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	log.Info("stopping application", slog.String("signal", sig.String()))

	application.GRPCApp.Stop()

	log.Info("gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
