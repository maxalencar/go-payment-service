package model

import "time"

// TransactionStatusUpdate represents an update to a transaction's status
type TransactionStatusUpdate struct {
	ID            string            `json:"id"`
	TransactionID string            `json:"transactionId"`
	Status        TransactionStatus `json:"status"`
	ReceivedAt    time.Time         `json:"receivedAt"`
	Details       string            `json:"details"`
}
