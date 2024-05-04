package grpcapp

import (
	"fmt"
	"net"

	"log/slog"

	authgrpc "github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/grpc/auth"
	permissionsgrpc "github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/grpc/permissions"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/services/auth"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/services/permissions"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	authSrv    *auth.Auth
	permSrv    *permissions.Permissions
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, auth *auth.Auth, perm *permissions.Permissions, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)
	permissionsgrpc.Register(gRPCServer, perm)

	return &App{
		log:        log,
		authSrv:    auth,
		permSrv:    perm,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC tests-server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))
	log.Info("Stopping gRPC tests-server")

	a.gRPCServer.GracefulStop()
}
