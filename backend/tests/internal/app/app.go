package app

import (
	"errors"
	"github.com/coddmeistr/quizzify/backend/tests/internal/app/mongoapp"
	"github.com/coddmeistr/quizzify/backend/tests/internal/app/rest"
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/internal/storage/mongo"
	errs "github.com/coddmeistr/quizzify/backend/tests/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	cfg        *config.Config
	log        *zap.Logger
	restAPIApp *rest.App
	mongoDbApp *mongoapp.App
}

func New(log *zap.Logger, cfg *config.Config) *App {
	mongoApp := mongoapp.New(log, cfg.MongoDB.ConnectionURI)

	return &App{
		cfg:        cfg,
		log:        log,
		mongoDbApp: mongoApp,
	}
}

func (a *App) MustRun() {
	a.log.Info("creating mongo application")
	done := a.mongoDbApp.MustRunConcurrently()
	<-done
	a.log.Info("mongo application now running")
	a.log.Info("creating rest api application")
	mongoStorage := mongo.New(a.mongoDbApp.Client().Database(a.cfg.MongoDB.DatabaseName))
	a.restAPIApp = rest.New(a.log, mongoStorage, a.cfg)
	go a.restAPIApp.MustRun()
}

func (a *App) Stop() error {
	if !errs.NoErrors(
		a.restAPIApp.Stop(),
		a.mongoDbApp.Stop(),
	) {
		a.log.Error("application was not gracefully shut down")
		return errors.New("graceful shutdown failed")
	}

	a.log.Info("application was gracefully shut down")
	return nil
}
