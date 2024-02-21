package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/config"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/helpers/user"
	testsservice "github.com/maxik12233/quizzify-online-tests/backend/tests/internal/service/tests"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/internal/storage/mongo"
	testshandlers "github.com/maxik12233/quizzify-online-tests/backend/tests/internal/transport/http/tests"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api/logging"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api/paginate"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/api/sort"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/cors"
	"github.com/maxik12233/quizzify-online-tests/backend/tests/pkg/metrics"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	log          *zap.Logger
	cfg          *config.Config
	testHandlers *testshandlers.Handlers
	testService  *testsservice.Service
	storage      *mongo.Storage
	server       *http.Server
}

func New(log *zap.Logger, storage *mongo.Storage, cfg *config.Config) *App {
	testServ := testsservice.New(log, cfg, storage)
	testHandlers := testshandlers.New(log, cfg, testServ)

	return &App{
		log:          log,
		cfg:          cfg,
		testHandlers: testHandlers,
		testService:  testServ,
		storage:      storage,
	}
}

func (a *App) MustRun() {
	a.log.Info("configuring rest api server")

	// TODO: move server configuration
	router := mux.NewRouter()
	router.Use(
		cors.Middleware,
		logging.RequestLogger(a.log),
		logging.ResponseLogger(a.log),
		paginate.Middleware(0, 5),  // TODO: create config values for paginate and sort middleware
		sort.Middleware("", "ASC"), // TODO: apply paginate and sort middleware only for needed routes
		user.Middleware(),
	)

	// Registering metrics endpoints
	m := metrics.New(a.log)
	m.Register(router)

	apiRouter := router.PathPrefix("/api").Subrouter()
	a.testHandlers.Register(apiRouter)

	addr := fmt.Sprintf("%s:%s", a.cfg.HTTPServer.Host, a.cfg.HTTPServer.Port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
		Handler:      router,
	}
	a.server = srv

	a.log.Info("starting server on port " + a.cfg.HTTPServer.Port)
	if err := srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			a.log.Info("server shutdown")
			return
		}
		panic(err)
	}
}

func (a *App) Stop() error {
	if a.server == nil {
		a.log.Error("trying to shutdown server that is nil")
		return errors.New("trying to shutdown server that is nil")
	}

	a.log.Info("starting rest api server gracefully shut down")
	if err := a.server.Shutdown(context.TODO()); err != nil {
		a.log.Error("failed to shutdown rest api server", zap.Error(err))
		return fmt.Errorf("failed to shutdown rest api server: %w", err)
	}

	a.log.Info("rest api server shutdown successful")
	return nil
}
