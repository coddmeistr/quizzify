package app

import (
	"errors"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/app/mongoapp"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/app/rest"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage/mongo"
	errs "github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/errors"
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

// MustRun TODO: this method creates new instances when called MustRun, fix this strange behavior
func (a *App) MustRun() {
	a.log.Info("starting mongo application")
	done := a.mongoDbApp.MustRunConcurrently()
	<-done
	a.log.Info("mongo application now running")
	a.log.Info("starting rest api application")
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
