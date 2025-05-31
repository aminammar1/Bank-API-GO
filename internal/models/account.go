package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account represents a Tunisian bank account following BCT (Central Bank of Tunisia) standards
type Account struct {
	ID              int       `json:"id" db:"id"`
	CustomerID      string    `json:"customer_id" db:"customer_id"`
	AccountNumber   string    `json:"account_number" db:"account_number"`
	IBAN            string    `json:"iban" db:"iban"`           // Tunisian IBAN format: TN59XXXXXXXXXXXXXXXXXXXX
	BIC             string    `json:"bic" db:"bic"`             // Tunisian bank BIC codes
	AccountType     string    `json:"account_type" db:"account_type"`
	Currency        string    `json:"currency" db:"currency"`  // TND primary, EUR/USD for foreign accounts
	Balance         int64     `json:"balance" db:"balance"`     // Amount in millimes (1 TND = 1000 millimes)
	AvailableBalance int64    `json:"available_balance" db:"available_balance"`
	HoldAmount      int64     `json:"hold_amount" db:"hold_amount"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Email           string    `json:"email" db:"email"`
	Phone           string    `json:"phone" db:"phone"`
	DateOfBirth     time.Time `json:"date_of_birth" db:"date_of_birth"`
	Address         Address   `json:"address" db:"address"`
	HashPassword    string    `json:"-" db:"hash_password"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
}

// Address represents customer address
type Address struct {
	Street     string `json:"street" db:"street"`
	City       string `json:"city" db:"city"`
	PostalCode string `json:"postal_code" db:"postal_code"`
	Country    string `json:"country" db:"country"`
	State      string `json:"state" db:"state"`
}

// Account status constants
const (
	AccountStatusActive    = "ACTIVE"
	AccountStatusInactive  = "INACTIVE"
	AccountStatusSuspended = "SUSPENDED"
	AccountStatusClosed    = "CLOSED"
)

// Tunisian account types following BCT regulations
const (
	AccountTypeChecking = "COMPTE_COURANT"    // Current account
	AccountTypeSavings  = "COMPTE_EPARGNE"    // Savings account  
	AccountTypeBusiness = "COMPTE_ENTREPRISE" // Business account
	AccountTypeForeign  = "COMPTE_DEVISES"    // Foreign currency account
)

// Supported currencies in Tunisia
const (
	CurrencyTND = "TND" // Tunisian Dinar (primary currency)
	CurrencyEUR = "EUR" // Euro (for foreign accounts)
	CurrencyUSD = "USD" // US Dollar (for foreign accounts)
)

// ValidatePassword checks if the provided password matches the hashed password
func (a *Account) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.HashPassword), []byte(password)) == nil
}

// IsActive checks if the account is active
func (a *Account) IsActive() bool {
	return a.Status == AccountStatusActive
}

// HasSufficientBalance checks if account has sufficient balance for a transaction
func (a *Account) HasSufficientBalance(amount int64) bool {
	return a.AvailableBalance >= amount
}

// GetAvailableBalance calculates available balance in millimes
func (a *Account) GetAvailableBalance() int64 {
	return a.Balance - a.HoldAmount
}

// ValidateTunisianIBAN validates Tunisian IBAN format (TN59 + 20 digits)
func ValidateTunisianIBAN(iban string) error {
	iban = strings.ReplaceAll(strings.ToUpper(iban), " ", "")
	
	// Tunisian IBAN: TN + 2 check digits + 20 digits = 24 characters total
	if len(iban) != 24 {
		return errors.New("IBAN tunisien doit contenir exactement 24 caractères")
	}
	
	// Must start with TN followed by 2 check digits
	if !strings.HasPrefix(iban, "TN") {
		return errors.New("IBAN tunisien doit commencer par 'TN'")
	}
	
	// Validate format: TN + 2 digits + 20 alphanumeric
	tunisianIBANRegex := regexp.MustCompile(`^TN[0-9]{2}[0-9]{20}$`)
	if !tunisianIBANRegex.MatchString(iban) {
		return errors.New("format IBAN tunisien invalide (TN + 2 chiffres de contrôle + 20 chiffres)")
	}
	
	return nil
}

// ValidateBIC validates BIC format
func ValidateBIC(bic string) error {
	bic = strings.ToUpper(bic)
	
	// BIC format: 4 letters (bank code) + 2 letters (country code) + 2 alphanumeric (location) + optional 3 alphanumeric (branch)
	bicRegex := regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	if !bicRegex.MatchString(bic) {
		return errors.New("invalid BIC format")
	}
	
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}
