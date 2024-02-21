package logging

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			next.ServeHTTP(w, r)
			latency := time.Since(t)
			message := fmt.Sprintf("REQUEST:: Method: %s; URI: %s; Proto:  %s; Latency: %s;\n",
				r.Method,
				r.RequestURI,
				r.Proto,
				latency)
			logger.Info(message)
		})
	}
}

func ResponseLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			logger.Warn("unimplemented response logging")
			// TODO: log response
			//message := fmt.Sprintf("RESPONSE:: Status: %d; Method: %s; URI: %s;\n",
			//	w.Writer().Status(),
			//	w.Request.Method,
			//	w.Request.RequestURI)
		})
	}
}
