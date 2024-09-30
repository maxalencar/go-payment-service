package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"go-payment-service/pkg/model"
)

type TransactionService interface {
	Deposit(ctx context.Context, req model.DepositRequest) (model.DepositResponse, error)
	Withdrawal(ctx context.Context, req model.WithdrawalRequest) (model.WithdrawalResponse, error)
	UpdateStatus(ctx context.Context, req model.TransactionStatusUpdate) error
	GetByID(ctx context.Context, id string) (*model.Transaction, error)
}

type transactionService struct {
	gateways   map[string]PaymentGateway
	repository TransactionRepository
	wg         *sync.WaitGroup
}

// newTransactionService creates a new transaction service
func newTransactionService(wg *sync.WaitGroup, gateways map[string]PaymentGateway, repo TransactionRepository) TransactionService {
	return &transactionService{
		gateways:   gateways,
		repository: repo,
		wg:         wg,
	}
}

func (s *transactionService) Deposit(ctx context.Context, req model.DepositRequest) (model.DepositResponse, error) {
	tx, err := s.create(req.BaseRequest, model.Deposit)
	if err != nil {
		slog.Debug("deposit: failed to create transaction", slog.Any("error", err))
		return model.DepositResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Process deposit transaction asynchronously
	errChain := s.process(ctx, tx)
	for err := range errChain {
		if err != nil {
			slog.Debug("deposit: failed to process transaction", slog.Any("error", err))
			return model.DepositResponse{}, err
		}
	}

	return model.DepositResponse{
		GatewayResponse: model.GatewayResponse{
			TransactionID: tx.ID,
			Status:        tx.Status,
			ProcessedAt:   tx.UpdatedAt,
		},
	}, nil
}

func (s *transactionService) Withdrawal(ctx context.Context, req model.WithdrawalRequest) (model.WithdrawalResponse, error) {
	tx, err := s.create(req.BaseRequest, model.Withdrawal)
	if err != nil {
		slog.Debug("withdrawal: failed to create transaction", slog.Any("error", err))
		return model.WithdrawalResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Process withdrawal transaction asynchronously
	errChain := s.process(ctx, tx)
	for err := range errChain {
		if err != nil {
			slog.Debug("withdrawal: failed to process transaction", slog.Any("error", err))
			return model.WithdrawalResponse{}, err
		}
	}

	return model.WithdrawalResponse{
		GatewayResponse: model.GatewayResponse{
			TransactionID: tx.ID,
			Status:        tx.Status,
			ProcessedAt:   tx.UpdatedAt,
		},
	}, nil
}

func (s *transactionService) UpdateStatus(ctx context.Context, req model.TransactionStatusUpdate) error {
	tx, err := s.repository.GetByExternalID(req.TransactionID)
	if err != nil {
		slog.Debug("update status: could not find transaction", slog.Any("error", err))
		return fmt.Errorf("could not find transaction. err: %w", err)
	}

	tx.Status = req.Status
	tx.UpdatedAt = time.Now()

	if err := s.repository.Update(tx); err != nil {
		slog.Debug("update status: could not update transaction", slog.Any("error", err))
		return fmt.Errorf("could not update transaction. err: %w", err)
	}

	return nil
}

func (s *transactionService) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	return s.repository.GetByID(id)
}

func (s *transactionService) create(req model.BaseRequest, transactionType model.TransactionType) (model.Transaction, error) {
	if _, exists := s.gateways[req.GatewayDetails.ID]; !exists {
		slog.Debug("create: unsupported payment gateway", slog.String("gateway", req.GatewayDetails.ID))
		return model.Transaction{}, errors.New("unsupported payment gateway")
	}

	tx := model.Transaction{
		ID:             uuid.New().String(),
		Amount:         req.Amount,
		CardDetails:    req.CardDetails,
		Type:           transactionType,
		Status:         model.Pending,
		GatewayDetails: req.GatewayDetails,
	}

	if err := s.repository.Create(&tx); err != nil {
		slog.Debug("create: could not create transaction", slog.Any("error", err))
		return model.Transaction{}, fmt.Errorf("could not create transaction: %w", err)
	}

	return tx, nil
}

func (s *transactionService) process(ctx context.Context, tx model.Transaction) <-chan error {
	errChan := make(chan error)

	gateway, exists := s.gateways[tx.GatewayDetails.ID]
	if !exists {
		slog.Debug("process: payment gateway not registered", slog.String("gateway", tx.GatewayDetails.ID))
		errChan <- errors.New("payment gateway not registered")
		close(errChan)

		return errChan
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	s.wg.Add(1)

	go func() {
		defer close(errChan)
		defer cancel()
		defer s.wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			res, err := gateway.ProcessTransaction(tx)
			if err != nil {
				tx.Status = model.Failed

				if err := s.repository.Update(&tx); err != nil {
					slog.Debug("process: could not update transaction", slog.Any("error", err))
					errChan <- fmt.Errorf("could not update transaction. err: %w", err)
					return
				}

				slog.Debug("process: could not process transaction", slog.Any("error", err))
				errChan <- fmt.Errorf("could not process transaction. err: %w", err)
				return
			}

			slog.Info("process: transaction processed", slog.Any("response", res))

			// Update transaction with external ID and status
			tx.ExternalID = res.TransactionID
			tx.Status = res.Status
			tx.UpdatedAt = time.Now()

			if err := s.repository.Update(&tx); err != nil {
				slog.Debug("process: could not update transaction", slog.Any("error", err))
				errChan <- fmt.Errorf("could not update transaction. err: %w", err)
				return
			}
		}
	}()

	return errChan
}
