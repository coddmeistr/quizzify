package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/app"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/logger/zaplogger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("couldn't load envs from .env file " + err.Error())
	}

	cfg := config.MustLoad()

	log, err := zaplogger.New(cfg.Env)
	if err != nil {
		panic(err)
	}

	log.Debug("config and logger initialized", zap.Any("config", cfg))
	log.Info("starting application's services")

	a := app.New(log, cfg)

	a.MustRun()

	// gracefully shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	log.Info("starting gracefully shutdown", zap.Any("signal", <-exit))

	if err := a.Stop(); err != nil {
		log.Error("program was not gracefully shut down")
		os.Exit(1)
		return
	}

	log.Info("program was gracefully shut down")
	os.Exit(0)
}
