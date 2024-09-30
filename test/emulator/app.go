package emulator

import (
	"log/slog"
	"net/http/httptest"
	"os"

	paymenthttp "go-payment-service/pkg/http"
)

func Start() *httptest.Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	client := paymenthttp.NewResilientHTTPClient()
	memoryRepository := newMemoryRepository()
	service := newService(client, memoryRepository)
	handler := newHandler(service)
	server := newServer(handler)

	return server.httpServer
}
