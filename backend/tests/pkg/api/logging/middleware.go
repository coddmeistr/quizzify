package logging

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lrw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			next.ServeHTTP(&lrw, r)
			latency := time.Since(t)
			message := fmt.Sprintf("HTTP Method=%s Uri=%s Proto=%s Code=%d ResponseLength=%d Latency=%s\n",
				r.Method,
				r.RequestURI,
				r.Proto,
				responseData.status,
				responseData.size,
				latency)
			logger.Info(message)
		})
	}
}
