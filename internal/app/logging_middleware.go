package app

import (
	"log/slog"
	"net/http"
	"time"

	paymenthttp "go-payment-service/pkg/http"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{w, http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		slog.InfoContext(
			r.Context(),
			"operation completed",
			slog.Any("status", wrapped.statusCode),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.Any("content-type", r.Header.Get(paymenthttp.HeaderContentType)),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
