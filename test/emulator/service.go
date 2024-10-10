package emulator

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	paymenthttp "go-payment-service/pkg/http"
	"go-payment-service/pkg/model"
)

type Service interface {
	Process(ctx context.Context, req ProcessRequest) (ProcessResponse, error)
	GetTransaction(ctx context.Context, id string) (*Transaction, error)
}

type service struct {
	client     paymenthttp.HTTPClient
	repository Repository
}

func newService(client paymenthttp.HTTPClient, repository Repository) Service {
	return &service{
		client:     client,
		repository: repository,
	}
}

func (s *service) Process(ctx context.Context, req ProcessRequest) (ProcessResponse, error) {
	// Create transaction
	tx := Transaction{
		ID:          uuid.New().String(),
		OrderID:     req.OrderID,
		Amount:      req.Amount,
		CardDetails: req.CardDetails,
		CallbackURL: req.CallbackURL,
		Type:        req.Type,
		Status:      model.Succeeded,
		CreatedAt:   time.Now(),
		RequestedAt: req.RequestedAt,
	}

	// Save transaction
	if err := s.repository.Create(&tx); err != nil {
		return ProcessResponse{}, err
	}

	// Send response via callback URL asynchronously
	s.sendTransactionUpdate(ctx, tx)

	return ProcessResponse{
		TransactionID: tx.ID,
		Status:        model.Succeeded,
		ProcessedAt:   tx.CreatedAt,
		Data:          req,
	}, nil
}

func (s *service) sendTransactionUpdate(ctx context.Context, tx Transaction) {
	if tx.CallbackURL == "" {
		return
	}

	slog.Info("emulator: sending transaction update",
		slog.Any("transaction-id", tx.ID),
		slog.Any("status", tx.Status),
		slog.Any("callback-url", tx.CallbackURL),
	)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)

	go func() {
		defer cancel()

		select {
		case <-ctx.Done():
			return
		default:
			tsu := TransactionStatusUpdate{
				TransactionID: tx.ID,
				Status:        tx.Status,
				ReceivedAt:    tx.UpdatedAt,
				Details:       "transaction processed",
			}

			contentType := ctx.Value(ContextKey(ContextKeyContentType)).(string)

			b, err := paymenthttp.Marshal(contentType, tsu)
			if err != nil {
				slog.Error("emulator: failed to marshal transaction status update", slog.Any("error", err))
				return
			}

			slog.Info("emulator: sending request", slog.Any("payload", string(b)))

			req, err := http.NewRequest(http.MethodPost, tx.CallbackURL, bytes.NewBuffer(b))
			if err != nil {
				slog.Error("emulator: failed to create HTTP request", slog.Any("error", err))
				return
			}

			req.Header.Add(paymenthttp.HeaderContentType, contentType)

			resp, err := s.client.Do(req)
			if err != nil {
				slog.Error("emulator: failed to send HTTP request", slog.Any("error", err))
				return
			}

			if resp.StatusCode != http.StatusOK {
				slog.Error("emulator: unexpected status code", slog.Any("status_code", resp.StatusCode))
				return
			}

			slog.Info("emulator: transaction update sent",
				slog.Any("transaction-id", tx.ID),
				slog.Any("status", tx.Status),
				slog.Any("callback-url", tx.CallbackURL),
			)
		}
	}()
}

func (s *service) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	return s.repository.GetByID(id)
}
