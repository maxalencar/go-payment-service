package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/test/emulator"
)

type Server interface {
	Start(port string) error
	StartTest() *httptest.Server
}

type server struct {
	handler *handler
	wg      *sync.WaitGroup
}

func NewServer() Server {
	s := &server{
		wg: &sync.WaitGroup{},
	}

	// Initialize HTTP client
	resilientHTTPClient := paymenthttp.NewResilientHTTPClient()

	// Initialize payment gateways
	gatewayEmulator := emulator.Start()

	gateways := map[string]PaymentGateway{
		"gatewayA": newGatewayAAdapter(resilientHTTPClient, gatewayEmulator.URL),
		"gatewayB": newGatewayBAdapter(resilientHTTPClient, gatewayEmulator.URL),
	}

	memoryRepository := newMemoryTransactionRepository()
	transactionService := newTransactionService(s.wg, gateways, memoryRepository)
	s.handler = newHandler(transactionService)

	return s
}

func (s *server) Start(port string) error {
	// Start HTTP server
	slog.Info("server: listening on", slog.String("port", port))

	// Setup signal catching
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           Logging(s.handler.mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Run server in a goroutine
	var err error
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server: ListenAndServe() error", slog.Any("error", err))
		}
	}()

	<-quit

	slog.Info("server: waiting for all transactions to complete...")
	s.wg.Wait()

	slog.Info("server: shutting down...")

	return err
}

func (s *server) StartTest() *httptest.Server {
	return httptest.NewServer(s.handler.mux)
}
