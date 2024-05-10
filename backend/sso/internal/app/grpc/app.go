package grpcapp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"

	"log/slog"

	authgrpc "github.com/coddmeistr/quizzify/backend/sso/internal/grpc/auth"
	permissionsgrpc "github.com/coddmeistr/quizzify/backend/sso/internal/grpc/permissions"
	"github.com/coddmeistr/quizzify/backend/sso/internal/services/auth"
	"github.com/coddmeistr/quizzify/backend/sso/internal/services/permissions"
	"google.golang.org/grpc"

	gw "github.com/coddmeistr/quizzify/backend/protos/proto/sso"
)

const gatewayPort = ":8001"

type App struct {
	log        *slog.Logger
	authSrv    *auth.Auth
	permSrv    *permissions.Permissions
	gRPCServer *grpc.Server
	port       int
}

// REST Gateway
func run(log *slog.Logger, grpcAddr string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterAuthHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}
	err = gw.RegisterPermissionHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Info("gRPC Gateway is listening on port " + gatewayPort)
	return http.ListenAndServe(gatewayPort, mux)
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

	log.Info("Starting gRPC gateway")
	go func() {
		if err := run(log, l.Addr().String()); err != nil {
			log.Error("FATAL: gRPC gateway error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

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
