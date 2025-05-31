package services

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/repository"
)

type AccountService interface {
	CreateAccount(req *models.CreateAccountRequest) (*models.Account, error)
	GetAccountByID(id int) (*models.Account, error)
	GetAccountByAccountNumber(accountNumber string) (*models.Account, error)
	GetAccountsByCustomerID(customerID string) ([]*models.Account, error)
	GetAllAccounts(limit, offset int) ([]*models.Account, error)
	UpdateAccount(id int, req *models.UpdateAccountRequest) error
	DeleteAccount(id int) error
	AuthenticateAccount(accountNumber, password string) (*models.Account, error)
	UpdateAccountStatus(id int, status string) error
	GetAccountBalance(accountNumber string) (*models.BalanceResponse, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (s *accountService) CreateAccount(req *models.CreateAccountRequest) (*models.Account, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}
		// Generate unique identifiers
	customerID := s.generateCustomerID()
	accountNumber := s.generateAccountNumber()
	iban := s.generateTunisianIBAN(accountNumber)
	bic := s.generateTunisianBIC()
	
	// Create account
	account, err := models.NewAccount(req, customerID, accountNumber, iban, bic)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	
	// Save to database
	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}
	
	// Clear password from response
	account.HashPassword = ""
	
	return account, nil
}

func (s *accountService) GetAccountByID(id int) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// Clear password from response
	account.HashPassword = ""
	return account, nil
}

func (s *accountService) GetAccountByAccountNumber(accountNumber string) (*models.Account, error) {
	account, err := s.accountRepo.GetByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	
	// Clear password from response
	account.HashPassword = ""
	return account, nil
}

func (s *accountService) GetAccountsByCustomerID(customerID string) ([]*models.Account, error) {
	accounts, err := s.accountRepo.GetByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	
	// Clear passwords from response
	for _, account := range accounts {
		account.HashPassword = ""
	}
	
	return accounts, nil
}

func (s *accountService) GetAllAccounts(limit, offset int) ([]*models.Account, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	
	accounts, err := s.accountRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}
	
	// Clear passwords from response
	for _, account := range accounts {
		account.HashPassword = ""
	}
	
	return accounts, nil
}

func (s *accountService) UpdateAccount(id int, req *models.UpdateAccountRequest) error {
	// Get existing account
	existingAccount, err := s.accountRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	
	// Update fields if provided
	if req.FirstName != "" {
		existingAccount.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingAccount.LastName = req.LastName
	}
	if req.Email != "" {
		if err := models.ValidateEmail(req.Email); err != nil {
			return err
		}
		existingAccount.Email = req.Email
	}
	if req.Phone != "" {
		if err := models.ValidatePhone(req.Phone); err != nil {
			return err
		}
		existingAccount.Phone = req.Phone
	}
	if req.Address.Street != "" || req.Address.City != "" || req.Address.PostalCode != "" || req.Address.Country != "" {
		existingAccount.Address = req.Address
	}
	
	existingAccount.UpdatedAt = time.Now().UTC()
	
	return s.accountRepo.Update(id, existingAccount)
}

func (s *accountService) DeleteAccount(id int) error {
	// Check if account exists
	account, err := s.accountRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	
	// Check account status (only allow deletion of inactive accounts)
	if account.Status == models.AccountStatusActive {
		return fmt.Errorf("cannot delete active account, please deactivate first")
	}
	
	return s.accountRepo.Delete(id)
}

func (s *accountService) AuthenticateAccount(accountNumber, password string) (*models.Account, error) {
	account, err := s.accountRepo.GetByAccountNumber(accountNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	if !account.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}
	
	if !account.ValidatePassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	// Clear password from response
	account.HashPassword = ""
	
	return account, nil
}

func (s *accountService) UpdateAccountStatus(id int, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		models.AccountStatusActive:    true,
		models.AccountStatusInactive:  true,
		models.AccountStatusSuspended: true,
		models.AccountStatusClosed:    true,
	}
	
	if !validStatuses[status] {
		return fmt.Errorf("invalid account status: %s", status)
	}
	
	return s.accountRepo.UpdateStatus(id, status)
}

func (s *accountService) GetAccountBalance(accountNumber string) (*models.BalanceResponse, error) {
	account, err := s.accountRepo.GetByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	
	return &models.BalanceResponse{
		AccountNumber:    account.AccountNumber,
		Balance:          account.Balance,
		AvailableBalance: account.GetAvailableBalance(),
		Currency:         account.Currency,
		HoldAmount:       account.HoldAmount,
	}, nil
}

// Helper methods for generating unique identifiers

func (s *accountService) generateCustomerID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	
	return fmt.Sprintf("CUST%d%x", timestamp, randomBytes)
}

func (s *accountService) generateAccountNumber() string {
	// Generate Tunisian account number format: 20 digits
	// Format: Bank code (3) + Branch code (3) + Account sequence (14)
	randomBytes := make([]byte, 10)
	rand.Read(randomBytes)
	
	var accountNumber strings.Builder
	for _, b := range randomBytes {
		accountNumber.WriteString(fmt.Sprintf("%02d", int(b)%100))
	}
	
	return accountNumber.String()
}

func (s *accountService) generateTunisianIBAN(accountNumber string) string {
	// Tunisian IBAN format: TN + 2 check digits + 20-digit account number
	countryCode := "TN"
	checkDigits := "59" // Simplified check digits for demo (should be calculated)
	
	// Ensure account number is exactly 20 digits
	if len(accountNumber) > 20 {
		accountNumber = accountNumber[:20]
	}
	for len(accountNumber) < 20 {
		accountNumber = "0" + accountNumber
	}
	
	return countryCode + checkDigits + accountNumber
}

func (s *accountService) generateTunisianBIC() string {
	// Common Tunisian bank BIC codes (for demo purposes)
	tunisianBICs := []string{
		"STBKTNTT", // STB (Société Tunisienne de Banque)
		"BIATTNTT", // BIAT (Banque Internationale Arabe de Tunisie)
		"BNTUTNTT", // BNA (Banque Nationale Agricole)
		"ABTNTNTT", // ATB (Arab Tunisian Bank)
		"UBCITNTT", // UBCI (Union Bancaire pour le Commerce et l'Industrie)
	}
	
	return tunisianBICs[0] // Using STB as default for demo
}
