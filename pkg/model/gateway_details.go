package model

// GatewayDetails holds information about the payment gateway used
type GatewayDetails struct {
	ID          string `json:"id" xml:"id" validate:"required"`
	Name        string `json:"name" xml:"name"`
	CallbackURL string `json:"callbackUrl" xml:"callbackUrl"` // URL where the payment gateway sends transaction updates
}
