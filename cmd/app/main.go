package main

import (
	"log/slog"
	"os"

	"go-payment-service/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := slog.LevelInfo
	if os.Getenv("APP_ENV") == "development" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	slog.SetDefault(logger)

	srv := app.NewServer()

	slog.Info("starting app", slog.Any("mode", logLevel))

	// Start server
	if err := srv.Start(port); err != nil {
		logger.Error("error starting server", slog.Any("error", err))
	}
}
