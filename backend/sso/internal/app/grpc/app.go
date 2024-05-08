package grpcapp

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"

	"log/slog"

	authgrpc "github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/grpc/auth"
	permissionsgrpc "github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/grpc/permissions"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/services/auth"
	"github.com/coddmeistr/quizzify-online-tests/backend/sso/internal/services/permissions"
	"google.golang.org/grpc"

	gw "github.com/coddmeistr/quizzify-online-tests/backend/protos/proto/sso"
)

type App struct {
	log        *slog.Logger
	authSrv    *auth.Auth
	permSrv    *permissions.Permissions
	gRPCServer *grpc.Server
	port       int
}

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "0.0.0.0:8000", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterAuthHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		fmt.Println("Failed to register gRPC gateway:", err)
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8001", mux)
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

	go run()

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
