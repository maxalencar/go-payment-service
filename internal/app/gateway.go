package app

import "go-payment-service/pkg/model"

// PaymentGateway represents an extensible payment gateway interface (protocol-agnostic)
// that can process transactions
type PaymentGateway interface {
	ProcessTransaction(tx model.Transaction) (model.GatewayResponse, error)
}
