package emulator

import (
	"net/http/httptest"
)

type server struct {
	httpServer *httptest.Server
}

func newServer(handler *handler) server {
	return server{
		httpServer: httptest.NewServer(handler.mux),
	}
}
