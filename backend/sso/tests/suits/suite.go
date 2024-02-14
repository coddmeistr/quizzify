package suits

import (
	"fmt"
	"net"
	"testing"

	ssov1 "github.com/maxik12233/quizzify-online-tests/backend/protos/gen/go/sso"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg         *config.Config
	AuthClient  ssov1.AuthClient
	PermsClient ssov1.PermissionClient
}

func NewDefault(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_test.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(ctx, grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc connection failed: %v", err)
	}

	return ctx, &Suite{
		T:           t,
		Cfg:         cfg,
		AuthClient:  ssov1.NewAuthClient(cc),
		PermsClient: ssov1.NewPermissionClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, fmt.Sprint(cfg.GRPC.Port))
}
