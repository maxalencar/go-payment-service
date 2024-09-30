package model

type Money struct {
	Amount   float64 `json:"amount" xml:"amount" validate:"gt=0"`
	Currency string  `json:"currency" xml:"currency" validate:"iso4217"`
}
