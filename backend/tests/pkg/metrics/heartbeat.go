package metrics

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

const (
	URL = "/heartbeat"
)

type Metric struct {
	log *zap.Logger
}

func New(logger *zap.Logger) *Metric {
	return &Metric{
		log: logger,
	}
}

func (m *Metric) Register(router *mux.Router) {
	router.Methods("GET").Path(URL).HandlerFunc(m.Heartbeat)
}

func (m *Metric) Heartbeat(w http.ResponseWriter, _ *http.Request) {
	m.log.Info("Health check OK")
	w.WriteHeader(204)
}
