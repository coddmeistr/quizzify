package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/internal/helpers/user"
	testsservice "github.com/coddmeistr/quizzify/backend/tests/internal/service/tests"
	"github.com/coddmeistr/quizzify/backend/tests/internal/service/tests/validation"
	"github.com/coddmeistr/quizzify/backend/tests/internal/storage/mongo"
	testshandlers "github.com/coddmeistr/quizzify/backend/tests/internal/transport/http/tests"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api/logging"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api/paginate"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/api/sort"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/cors"
	"github.com/coddmeistr/quizzify/backend/tests/pkg/metrics"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	testValidator := validation.NewValidation(cfg, log)
	testServ := testsservice.New(log, cfg, storage, testValidator)
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
	a.log.Info("attempting to run rest api application")

	server := a.createServer()

	a.log.Info("starting server on port " + a.cfg.HTTPServer.Port)
	if err := server.ListenAndServe(); err != nil {
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

func (a *App) createServer() *http.Server {
	a.log.Info("creating rest api server")

	router := mux.NewRouter()
	router.Use(
		cors.Middleware,
		logging.RequestLogger(a.log),
		paginate.Middleware(a.cfg.Other.DefaultPage, a.cfg.Other.DefaultPerPage),
		sort.Middleware(a.cfg.Other.DefaultSortField, a.cfg.Other.DefaultSortOrder),
		user.Middleware(),
	)

	// Registering metrics endpoints
	m := metrics.New(a.log)
	m.Register(router)

	apiRouter := router.PathPrefix("/api").Subrouter()
	a.testHandlers.Register(apiRouter)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Auth-User-Info"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})

	addr := fmt.Sprintf("%s:%s", a.cfg.HTTPServer.Host, a.cfg.HTTPServer.Port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.IdleTimeout,
		Handler:      handlers.CORS(headersOk, originsOk, methodsOk)(router),
	}
	a.server = srv

	a.log.Info("rest api server created successfully")
	return srv
}
