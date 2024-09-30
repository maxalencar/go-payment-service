package app

import (
	"fmt"
	"sync"
	"time"

	"go-payment-service/pkg/model"
)

// TransactionRepository defines the methods for transaction data access
type TransactionRepository interface {
	Create(tx *model.Transaction) error
	GetByID(id string) (*model.Transaction, error)
	GetByExternalID(externalID string) (*model.Transaction, error)
	List() []*model.Transaction
	Update(tx *model.Transaction) error
}

// memoryTransactionRepository represents an in-memory repository for transactions.
type memoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[string]*model.Transaction
}

// newMemoryRepository creates a new in-memory transaction repository.
func newMemoryTransactionRepository() TransactionRepository {
	return &memoryTransactionRepository{
		transactions: make(map[string]*model.Transaction),
	}
}

// Create adds a new transaction to the repository.
func (r *memoryTransactionRepository) Create(tx *model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction %s already exists", tx.ID)
	}

	tx.CreatedAt = time.Now()
	r.transactions[tx.ID] = tx

	return nil
}

// GetByID retrieves a transaction by its ID.
func (r *memoryTransactionRepository) GetByID(id string) (*model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, fmt.Errorf("transaction %s not found", id)
	}

	return tx, nil
}

// GetByExternalID retrieves a transaction by its external ID.
func (r *memoryTransactionRepository) GetByExternalID(externalID string) (*model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tx := range r.transactions {
		if tx.ExternalID == externalID {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction with external ID %s not found", externalID)
}

// List returns all transactions.
func (r *memoryTransactionRepository) List() []*model.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var txList []*model.Transaction
	for _, tx := range r.transactions {
		txList = append(txList, tx)
	}

	return txList
}

// Update updates the transaction.
func (r *memoryTransactionRepository) Update(tx *model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[tx.ID]; !exists {
		return fmt.Errorf("transaction %s not found", tx.ID)
	}

	tx.UpdatedAt = time.Now()
	r.transactions[tx.ID] = tx

	return nil
}
