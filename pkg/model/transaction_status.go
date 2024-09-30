package model

// TransactionStatus represents the current status of a transaction
type TransactionStatus string

const (
	Failed     TransactionStatus = "failed"
	Pending    TransactionStatus = "pending"
	Processing TransactionStatus = "processing"
	Succeeded  TransactionStatus = "succeeded"
)
