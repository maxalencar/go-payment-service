package app

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type handler struct {
	mux      *http.ServeMux
	service  TransactionService
	validate *validator.Validate
}

func newHandler(service TransactionService) *handler {
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
	mux.HandleFunc("POST /deposit", h.deposit)
	mux.HandleFunc("POST /withdrawal", h.withdrawal)
	mux.HandleFunc("POST /callback", h.callback)
	mux.HandleFunc("GET /transactions/{id}", h.getTransaction)

	h.mux = mux
}

func (h *handler) deposit(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(paymenthttp.HeaderContentType)

	// decode request
	var req model.DepositRequest
	if err := paymenthttp.Decode(r.Body, contentType, &req); err != nil {
		slog.Debug("failed to decode deposit request", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		slog.Debug("failed to validate deposit request", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// process request
	res, err := h.service.Deposit(r.Context(), req)
	if err != nil {
		slog.Debug("failed to process deposit", slog.Any("error", err))
		h.errorResponse(w, contentType, http.StatusInternalServerError, err.Error())
		return
	}

	// encode response
	if err := paymenthttp.Encode(w, contentType, res); err != nil {
		slog.Debug("failed to encode deposit response", slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *handler) withdrawal(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(paymenthttp.HeaderContentType)

	// decode request
	var req model.WithdrawalRequest
	if err := paymenthttp.Decode(r.Body, contentType, &req); err != nil {
		slog.Debug("failed to decode withdrawal request", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		slog.Debug("failed to validate withdrawal request", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process request
	res, err := h.service.Withdrawal(r.Context(), req)
	if err != nil {
		slog.Debug("failed to process withdrawal", slog.Any("error", err))
		h.errorResponse(w, contentType, http.StatusInternalServerError, err.Error())
		return
	}

	// encode response
	if err := paymenthttp.Encode(w, contentType, res); err != nil {
		slog.Debug("failed to encode withdrawal response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) callback(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(paymenthttp.HeaderContentType)

	// decode request
	var req model.TransactionStatusUpdate
	if err := paymenthttp.Decode(r.Body, contentType, &req); err != nil {
		slog.Debug("failed to decode transaction status update", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate request
	if err := h.validate.Struct(req); err != nil {
		slog.Debug("failed to validate transaction status update", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process request
	if err := h.service.UpdateStatus(r.Context(), req); err != nil {
		slog.Debug("failed to update transaction status", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get(paymenthttp.HeaderContentType)

	// get transaction ID from path
	id := r.PathValue("id")

	// get transaction
	tx, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		slog.Debug("failed to get transaction", slog.Any("error", err))
		h.errorResponse(w, contentType, http.StatusInternalServerError, err.Error())
		return
	}

	// encode response
	if err := paymenthttp.Encode(w, contentType, tx); err != nil {
		slog.Debug("failed to encode transaction", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) errorResponse(w http.ResponseWriter, contentType string, code int, message string) {
	er := model.ErrorResponse{
		Code:    code,
		Message: message,
	}

	b, err := paymenthttp.Marshal(contentType, er)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Error(w, string(b), code)
}
