package logging

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)

		message := fmt.Sprintf("REQUEST:: Method: %s; URI: %s; Proto:  %s; Latency: %s;\n",
			c.Request.Method,
			c.Request.RequestURI,
			c.Request.Proto,
			latency)
		logger.Info(message)
	}
}

func ResponseLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		message := fmt.Sprintf("RESPONSE:: Status: %d; Method: %s; URI: %s;\n",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.RequestURI)
		logger.Info(message)
	}
}
