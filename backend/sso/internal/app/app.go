package app

import (
	"time"

	"log/slog"

	grpcapp "github.com/maxik12233/quizzify-online-tests/backend/sso/internal/app/grpc"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/services/auth"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/services/permissions"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/storage/postgres"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, postgresURL string, tokenTTL time.Duration) *App {

	// Init storage
	storage, err := postgres.New(postgresURL)
	if err != nil {
		panic(err)
	}

	// Init auth service
	authSrv := auth.New(log, storage, storage, storage, storage, tokenTTL)

	// Init permissions service
	permSrv := permissions.New(log, storage)

	// Init gRPC app
	grpcApp := grpcapp.New(log, authSrv, permSrv, grpcPort)

	return &App{
		GRPCApp: grpcApp,
	}
}
