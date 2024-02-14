package grpcapp

import (
	"fmt"
	"net"

	"log/slog"

	authgrpc "github.com/maxik12233/quizzify-online-tests/backend/sso/internal/grpc/auth"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/services/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	authSrv    *auth.Auth
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, auth *auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)

	return &App{
		log:        log,
		authSrv:    auth,
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

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))
	log.Info("Stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
