package models

import (
	"errors"
	"time"
)

// Transaction represents a financial transaction following international standards
type Transaction struct {
	ID                    int       `json:"id" db:"id"`
	TransactionID         string    `json:"transaction_id" db:"transaction_id"`
	FromAccountID         int       `json:"from_account_id" db:"from_account_id"`
	ToAccountID           int       `json:"to_account_id" db:"to_account_id"`
	FromAccountNumber     string    `json:"from_account_number" db:"from_account_number"`
	ToAccountNumber       string    `json:"to_account_number" db:"to_account_number"`
	Amount                int64     `json:"amount" db:"amount"` // Amount in cents
	Currency              string    `json:"currency" db:"currency"`
	ExchangeRate          float64   `json:"exchange_rate" db:"exchange_rate"`
	ConvertedAmount       int64     `json:"converted_amount" db:"converted_amount"`
	TransactionType       string    `json:"transaction_type" db:"transaction_type"`
	Status                string    `json:"status" db:"status"`
	Description           string    `json:"description" db:"description"`
	Reference             string    `json:"reference" db:"reference"`
	Fee                   int64     `json:"fee" db:"fee"`
	ProcessedAt           *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	FailureReason         string    `json:"failure_reason,omitempty" db:"failure_reason"`
}

// Transaction type constants
const (
	TransactionTypeTransfer    = "TRANSFER"
	TransactionTypeDeposit     = "DEPOSIT"
	TransactionTypeWithdrawal  = "WITHDRAWAL"
	TransactionTypePayment     = "PAYMENT"
	TransactionTypeFee         = "FEE"
	TransactionTypeInterest    = "INTEREST"
)

// Transaction status constants
const (
	TransactionStatusPending   = "PENDING"
	TransactionStatusCompleted = "COMPLETED"
	TransactionStatusFailed    = "FAILED"
	TransactionStatusCancelled = "CANCELLED"
)

// ValidateTransaction validates transaction data
func (t *Transaction) ValidateTransaction() error {
	if t.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}
	
	if t.Currency == "" {
		return errors.New("currency is required")
	}
	
	if t.TransactionType == "" {
		return errors.New("transaction type is required")
	}
	
	if t.TransactionType == TransactionTypeTransfer {
		if t.FromAccountID == t.ToAccountID {
			return errors.New("cannot transfer to the same account")
		}
		if t.ToAccountID == 0 {
			return errors.New("destination account is required for transfers")
		}
	}
	
	return nil
}

// IsCompleted checks if transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsPending checks if transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// CanBeCancelled checks if transaction can be cancelled
func (t *Transaction) CanBeCancelled() bool {
	return t.Status == TransactionStatusPending
}
