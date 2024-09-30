package emulator

import (
	"time"

	"go-payment-service/pkg/model"
)

type Transaction struct {
	ID          string                  `json:"id"`
	OrderID     string                  `json:"orderId"`
	Amount      model.Money             `json:"amount"`
	CardDetails model.CardDetails       `json:"cardDetails"`
	CallbackURL string                  `json:"callbackUrl"` // URL where the payment gateway sends transaction updates
	Type        model.TransactionType   `json:"type"`
	Status      model.TransactionStatus `json:"status"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
	RequestedAt time.Time               `json:"requestedAt"`
}

// ProcessRequest represents a request to a payment gateway
type ProcessRequest struct {
	OrderID     string                `json:"orderId" xml:"orderId" validate:"required"`
	Amount      model.Money           `json:"amount" xml:"amount" validate:"required"`
	CardDetails model.CardDetails     `json:"cardDetails" xml:"cardDetails" validate:"required"`
	CallbackURL string                `json:"callbackUrl" xml:"callbackUrl"` // URL where the payment gateway sends transaction updates
	Type        model.TransactionType `json:"type" xml:"type" validate:"required"`
	RequestedAt time.Time             `json:"requestedAt" xml:"requestedAt"`
}

// ProcessResponse represents the response from a payment gateway
type ProcessResponse struct {
	TransactionID string                  `json:"transactionId" xml:"transactionId"`
	Data          ProcessRequest          `json:"data,omitempty" xml:"data,omitempty"`
	Message       string                  `json:"message,omitempty" xml:"message,omitempty"`
	Status        model.TransactionStatus `json:"status" xml:"status"`
	ProcessedAt   time.Time               `json:"processedAt" xml:"processedAt"`
}

// TransactionStatusUpdate represents an update to a transaction's status
type TransactionStatusUpdate struct {
	ID            string                  `json:"id"`
	TransactionID string                  `json:"transactionId"`
	Status        model.TransactionStatus `json:"status"`
	ReceivedAt    time.Time               `json:"receivedAt"`
	Details       any                     `json:"details"`
}
