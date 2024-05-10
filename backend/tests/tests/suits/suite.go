package suits

import (
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"net/http"
	"os"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client *http.Client
}

func NewDefault(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	err := os.Setenv("MONGO_CONNECTION_URI", "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	cfg := config.MustLoadByPath("../configs/local.yaml")

	client := http.DefaultClient
	client.Timeout = cfg.HTTPServer.Timeout

	return &Suite{
		T:      t,
		Cfg:    cfg,
		Client: client,
	}
}
