package metrics

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	URL = "/heartbeat"
)

type Metric struct {
	Logger *zap.Logger
}

func (m *Metric) Register(router *gin.Engine) {
	router.GET(URL, m.Heartbeat)
}

func (m *Metric) Heartbeat(c *gin.Context) {
	m.Logger.Info("Health check OK")
	c.Writer.WriteHeader(204)
}
