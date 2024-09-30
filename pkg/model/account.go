package model

import "time"

// Account represents a user account in the system
type Account struct {
	ID        string    `json:"id"`
	Balance   float64   `json:"balance"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
