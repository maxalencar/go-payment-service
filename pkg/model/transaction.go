package model

import (
	"time"
)

// TransactionType represents the type of transaction: Deposit or Withdrawal
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

// Transaction represents a financial transaction (deposit or withdrawal)
type Transaction struct {
	ID             string            `json:"id"`
	Amount         Money             `json:"amount"`
	CardDetails    CardDetails       `json:"cardDetails"`
	GatewayDetails GatewayDetails    `json:"gatewayDetails"`
	Type           TransactionType   `json:"type"`
	Status         TransactionStatus `json:"status"`
	ExternalID     string            `json:"externalId"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}
