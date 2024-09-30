package model

// DepositRequest represents a deposit request
type DepositRequest struct {
	BaseRequest
}

// WithdrawalRequest represents a withdrawal request
type WithdrawalRequest struct {
	BaseRequest
}

type BaseRequest struct {
	Amount         Money          `json:"amount" xml:"amount" validate:"required"`
	CardDetails    CardDetails    `json:"cardDetails" xml:"cardDetails" validate:"required"`
	GatewayDetails GatewayDetails `json:"gatewayDetails" xml:"gatewayDetails" validate:"required"`
}

// CallbackRequest represents a callback request from a payment gateway
type CallbackRequest struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Details       string `json:"details"`
}

// GatewayRequest represents a request to a payment gateway
type GatewayRequest struct {
	OrderID     string          `json:"orderId" xml:"orderId" validate:"required"`
	Amount      Money           `json:"amount" xml:"amount" validate:"required"`
	CardDetails CardDetails     `json:"cardDetails" xml:"cardDetails" validate:"required"`
	CallbackURL string          `json:"callbackUrl" xml:"callbackUrl"` // URL where the payment gateway sends transaction updates
	Type        TransactionType `json:"type" xml:"type" validate:"required"`
}
