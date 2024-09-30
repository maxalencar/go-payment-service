package model

import "time"

type DepositResponse struct {
	GatewayResponse
}

type WithdrawalResponse struct {
	GatewayResponse
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// GatewayResponse represents the response from a payment gateway
type GatewayResponse struct {
	TransactionID string            `json:"transactionId" xml:"transactionId"`
	Data          any               `json:"data,omitempty" xml:"data,omitempty"`
	Message       string            `json:"message,omitempty" xml:"message,omitempty"`
	Status        TransactionStatus `json:"status" xml:"status"`
	ProcessedAt   time.Time         `json:"processedAt" xml:"processedAt"`
}
