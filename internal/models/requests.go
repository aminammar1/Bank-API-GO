package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// CreateAccountRequest represents the request payload for creating a new account
type CreateAccountRequest struct {
	FirstName    string    `json:"first_name" validate:"required,min=2,max=50"`
	LastName     string    `json:"last_name" validate:"required,min=2,max=50"`
	Email        string    `json:"email" validate:"required,email"`
	Phone        string    `json:"phone" validate:"required"`
	Password     string    `json:"password" validate:"required,min=8"`
	DateOfBirth  time.Time `json:"date_of_birth" validate:"required"`
	AccountType  string    `json:"account_type" validate:"required"`
	Currency     string    `json:"currency" validate:"required"`
	Address      Address   `json:"address" validate:"required"`
}

// UpdateAccountRequest represents the request payload for updating an account
type UpdateAccountRequest struct {
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	Email       string  `json:"email,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Address     Address `json:"address,omitempty"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	AccountNumber string `json:"account_number" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token         string    `json:"token"`
	AccountNumber string    `json:"account_number"`
	CustomerID    string    `json:"customer_id"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// TransferRequest represents a transfer request payload
type TransferRequest struct {
	FromAccountNumber string `json:"from_account_number" validate:"required"`
	ToAccountNumber   string `json:"to_account_number" validate:"required"`
	Amount            int64  `json:"amount" validate:"required,min=1"`
	Currency          string `json:"currency" validate:"required"`
	Description       string `json:"description,omitempty"`
	Reference         string `json:"reference,omitempty"`
}

// DepositRequest represents a deposit request payload
type DepositRequest struct {
	AccountNumber string `json:"account_number" validate:"required"`
	Amount        int64  `json:"amount" validate:"required,min=1"`
	Currency      string `json:"currency" validate:"required"`
	Description   string `json:"description,omitempty"`
	Reference     string `json:"reference,omitempty"`
}

// WithdrawalRequest represents a withdrawal request payload
type WithdrawalRequest struct {
	AccountNumber string `json:"account_number" validate:"required"`
	Amount        int64  `json:"amount" validate:"required,min=1"`
	Currency      string `json:"currency" validate:"required"`
	Description   string `json:"description,omitempty"`
	Reference     string `json:"reference,omitempty"`
}

// TransactionHistoryRequest represents request for transaction history
type TransactionHistoryRequest struct {
	AccountNumber string    `json:"account_number" validate:"required"`
	StartDate     time.Time `json:"start_date,omitempty"`
	EndDate       time.Time `json:"end_date,omitempty"`
	Limit         int       `json:"limit,omitempty"`
	Offset        int       `json:"offset,omitempty"`
	Type          string    `json:"type,omitempty"`
}

// BalanceResponse represents account balance response
type BalanceResponse struct {
	AccountNumber    string `json:"account_number"`
	Balance          int64  `json:"balance"`
	AvailableBalance int64  `json:"available_balance"`
	Currency         string `json:"currency"`
	HoldAmount       int64  `json:"hold_amount"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      string    `json:"code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Validate validates the CreateAccountRequest
func (r *CreateAccountRequest) Validate() error {
	if r.FirstName == "" || len(r.FirstName) < 2 {
		return errors.New("first name is required and must be at least 2 characters")
	}
	
	if r.LastName == "" || len(r.LastName) < 2 {
		return errors.New("last name is required and must be at least 2 characters")
	}
	
	if err := ValidateEmail(r.Email); err != nil {
		return err
	}
	
	if err := ValidatePhone(r.Phone); err != nil {
		return err
	}
	
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	
	if r.AccountType == "" {
		return errors.New("account type is required")
	}
	
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	
	// Check if account type is valid
	validAccountTypes := map[string]bool{
		AccountTypeChecking: true,
		AccountTypeSavings:  true,
		AccountTypeBusiness: true,
	}
	
	if !validAccountTypes[r.AccountType] {
		return errors.New("invalid account type")
	}
		// Check if currency is valid for Tunisian banking
	validCurrencies := map[string]bool{
		CurrencyTND: true,
		CurrencyEUR: true, 
		CurrencyUSD: true,
	}
	
	if !validCurrencies[r.Currency] {
		return errors.New("invalid currency")
	}
	
	return nil
}

// NewAccount creates a new account from CreateAccountRequest
func NewAccount(req *CreateAccountRequest, customerID, accountNumber, iban, bic string) (*Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	account := &Account{
		CustomerID:       customerID,
		AccountNumber:    accountNumber,
		IBAN:            iban,
		BIC:             bic,
		AccountType:     req.AccountType,
		Currency:        req.Currency,
		Balance:         0,
		AvailableBalance: 0,
		HoldAmount:      0,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Phone:           req.Phone,
		DateOfBirth:     req.DateOfBirth,
		Address:         req.Address,
		HashPassword:    string(hashedPassword),
		Status:          AccountStatusActive,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	
	return account, nil
}
