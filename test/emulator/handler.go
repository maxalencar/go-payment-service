package emulator

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"

	paymenthttp "go-payment-service/pkg/http"
)

type handler struct {
	mux      *http.ServeMux
	service  Service
	validate *validator.Validate
}

func newHandler(service Service) *handler {
	h := handler{
		service:  service,
		validate: validator.New(),
	}

	h.registerRoutes()

	return &h
}

func (h *handler) registerRoutes() {
	// Initialize HTTP request multiplexer
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /process", h.process)
	mux.HandleFunc("GET /{id}", h.getTransaction)

	h.mux = mux
}

func (h *handler) process(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(paymenthttp.HeaderContentType)

	// add content type to context
	// to handle multiple formats in the callback
	ctx := context.WithValue(r.Context(), ContextKey(ContextKeyContentType), contentType)

	// decode request
	var req ProcessRequest
	if err := paymenthttp.Decode(r.Body, contentType, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process request
	resp, err := h.service.Process(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode response
	if err := paymenthttp.Encode(w, contentType, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	// get transaction ID from URL
	id := r.PathValue("id")

	// get transaction
	tx, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode response
	if err := paymenthttp.Encode(w, paymenthttp.MIMETypeJSON, tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
