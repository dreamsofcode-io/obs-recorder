package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

type Middleware func(next http.Handler) http.Handler

func loggingMw(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(wrapped, r)

		logger.Info(
			"request handled",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("duration", time.Since(start).String()),
		)
	})
}

func Logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return loggingMw(logger, next)
	}
}
