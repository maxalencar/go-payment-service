package emulator

import (
	"fmt"
	"sync"
	"time"
)

// Repository defines the methods for transaction data access
type Repository interface {
	Create(tx *Transaction) error
	GetByID(id string) (*Transaction, error)
	List() []*Transaction
	Update(tx *Transaction) error
}

// memoryRepository represents an in-memory repository for transactions.
type memoryRepository struct {
	mu           sync.RWMutex
	transactions map[string]*Transaction
}

// newMemoryRepository creates a new in-memory transaction repository.
func newMemoryRepository() Repository {
	return &memoryRepository{
		transactions: make(map[string]*Transaction),
	}
}

// Create adds a new transaction to the repository.
func (r *memoryRepository) Create(tx *Transaction) error {
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
func (r *memoryRepository) GetByID(id string) (*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, fmt.Errorf("transaction %s not found", id)
	}

	return tx, nil
}

// List returns all transactions.
func (r *memoryRepository) List() []*Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var txList []*Transaction
	for _, tx := range r.transactions {
		txList = append(txList, tx)
	}

	return txList
}

// Update updates the transaction.
func (r *memoryRepository) Update(tx *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[tx.ID]; !exists {
		return fmt.Errorf("transaction %s not found", tx.ID)
	}

	tx.UpdatedAt = time.Now()
	r.transactions[tx.ID] = tx

	return nil
}
