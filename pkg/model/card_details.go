package model

// CardDetails holds information about the card used in the transaction
type CardDetails struct {
	Name        string `json:"name" xml:"name" validate:"required"`
	Number      string `json:"number" xml:"number" validate:"credit_card"`
	Type        string `json:"type" xml:"type"` // e.g. Visa, MasterCard, etc.
	ExpiryMonth int    `json:"expiryMonth" xml:"expiryMonth" validate:"min=1,max=12"`
	ExpiryYear  int    `json:"expiryYear" xml:"expiryYear" validate:"min=2021,max=2040"`
	CVV         string `json:"cvv" xml:"cvv" validate:"min=3,max=4"`
}
