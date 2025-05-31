package services

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/repository"
)

type TransactionService interface {
	Transfer(req *models.TransferRequest) (*models.Transaction, error)
	Deposit(req *models.DepositRequest) (*models.Transaction, error)
	Withdraw(req *models.WithdrawalRequest) (*models.Transaction, error)
	GetTransaction(transactionID string) (*models.Transaction, error)
	GetTransactionHistory(req *models.TransactionHistoryRequest) ([]*models.Transaction, error)
	ProcessPendingTransactions() error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, accountRepo repository.AccountRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *transactionService) Transfer(req *models.TransferRequest) (*models.Transaction, error) {
	// Validate request
	if req.Amount <= 0 {
		return nil, fmt.Errorf("transfer amount must be positive")
	}
	
	if req.FromAccountNumber == req.ToAccountNumber {
		return nil, fmt.Errorf("cannot transfer to the same account")
	}
	
	// Get source account
	fromAccount, err := s.accountRepo.GetByAccountNumber(req.FromAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("source account not found")
	}
	
	if !fromAccount.IsActive() {
		return nil, fmt.Errorf("source account is not active")
	}
	
	// Get destination account
	toAccount, err := s.accountRepo.GetByAccountNumber(req.ToAccountNumber)
	if err != nil {
		return nil, fmt.Errorf("destination account not found")
	}
	
	if !toAccount.IsActive() {
		return nil, fmt.Errorf("destination account is not active")
	}
	
	// Check sufficient balance
	if !fromAccount.HasSufficientBalance(req.Amount) {
		return nil, fmt.Errorf("insufficient balance")
	}
	
	// Calculate fee (simplified - 0.1% of transfer amount, minimum $1)
	fee := s.calculateTransferFee(req.Amount)
	totalAmount := req.Amount + fee
	
	if !fromAccount.HasSufficientBalance(totalAmount) {
		return nil, fmt.Errorf("insufficient balance including fees")
	}
	
	// Create transaction
	transaction := &models.Transaction{
		TransactionID:     s.generateTransactionID(),
		FromAccountID:     fromAccount.ID,
		ToAccountID:       toAccount.ID,
		FromAccountNumber: req.FromAccountNumber,
		ToAccountNumber:   req.ToAccountNumber,
		Amount:            req.Amount,
		Currency:          req.Currency,
		ExchangeRate:      1.0, // Simplified - same currency
		ConvertedAmount:   req.Amount,
		TransactionType:   models.TransactionTypeTransfer,
		Status:            models.TransactionStatusPending,
		Description:       req.Description,
		Reference:         req.Reference,
		Fee:               fee,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	
	// Validate transaction
	if err := transaction.ValidateTransaction(); err != nil {
		return nil, err
	}
	
	// Save transaction
	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	
	// Process transaction immediately (in real system, this might be async)
	if err := s.processTransfer(transaction); err != nil {
		// Update transaction status to failed
		s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusFailed)
		return nil, err
	}
	
	return transaction, nil
}

func (s *transactionService) Deposit(req *models.DepositRequest) (*models.Transaction, error) {
	// Validate request
	if req.Amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be positive")
	}
	
	// Get account
	account, err := s.accountRepo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	
	if !account.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}
	
	// Create transaction
	transaction := &models.Transaction{
		TransactionID:     s.generateTransactionID(),
		ToAccountID:       account.ID,
		ToAccountNumber:   req.AccountNumber,
		Amount:            req.Amount,
		Currency:          req.Currency,
		ExchangeRate:      1.0,
		ConvertedAmount:   req.Amount,
		TransactionType:   models.TransactionTypeDeposit,
		Status:            models.TransactionStatusPending,
		Description:       req.Description,
		Reference:         req.Reference,
		Fee:               0, // No fee for deposits
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	
	// Save transaction
	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	
	// Process deposit
	if err := s.processDeposit(transaction); err != nil {
		s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusFailed)
		return nil, err
	}
	
	return transaction, nil
}

func (s *transactionService) Withdraw(req *models.WithdrawalRequest) (*models.Transaction, error) {
	// Validate request
	if req.Amount <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be positive")
	}
	
	// Get account
	account, err := s.accountRepo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	
	if !account.IsActive() {
		return nil, fmt.Errorf("account is not active")
	}
	
	// Calculate fee (simplified - $2 per withdrawal)
	fee := s.calculateWithdrawalFee()
	totalAmount := req.Amount + fee
	
	// Check sufficient balance
	if !account.HasSufficientBalance(totalAmount) {
		return nil, fmt.Errorf("insufficient balance")
	}
	
	// Create transaction
	transaction := &models.Transaction{
		TransactionID:     s.generateTransactionID(),
		FromAccountID:     account.ID,
		FromAccountNumber: req.AccountNumber,
		Amount:            req.Amount,
		Currency:          req.Currency,
		ExchangeRate:      1.0,
		ConvertedAmount:   req.Amount,
		TransactionType:   models.TransactionTypeWithdrawal,
		Status:            models.TransactionStatusPending,
		Description:       req.Description,
		Reference:         req.Reference,
		Fee:               fee,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}
	
	// Save transaction
	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	
	// Process withdrawal
	if err := s.processWithdrawal(transaction); err != nil {
		s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusFailed)
		return nil, err
	}
	
	return transaction, nil
}

func (s *transactionService) GetTransaction(transactionID string) (*models.Transaction, error) {
	return s.transactionRepo.GetByTransactionID(transactionID)
}

func (s *transactionService) GetTransactionHistory(req *models.TransactionHistoryRequest) ([]*models.Transaction, error) {
	// Validate account exists
	_, err := s.accountRepo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	
	// Set default limit if not provided
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	
	// Get transactions by date range if provided
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		return s.transactionRepo.GetByDateRange(req.AccountNumber, req.StartDate, req.EndDate, req.Limit, req.Offset)
	}
	
	// Get all transactions for account
	return s.transactionRepo.GetByAccountNumber(req.AccountNumber, req.Limit, req.Offset)
}

func (s *transactionService) ProcessPendingTransactions() error {
	transactions, err := s.transactionRepo.GetPendingTransactions()
	if err != nil {
		return err
	}
	
	for _, transaction := range transactions {
		switch transaction.TransactionType {
		case models.TransactionTypeTransfer:
			s.processTransfer(transaction)
		case models.TransactionTypeDeposit:
			s.processDeposit(transaction)
		case models.TransactionTypeWithdrawal:
			s.processWithdrawal(transaction)
		}
	}
	
	return nil
}

// Helper methods

func (s *transactionService) processTransfer(transaction *models.Transaction) error {
	// Get accounts
	fromAccount, err := s.accountRepo.GetByID(transaction.FromAccountID)
	if err != nil {
		return err
	}
	
	toAccount, err := s.accountRepo.GetByID(transaction.ToAccountID)
	if err != nil {
		return err
	}
	
	// Update balances
	totalAmount := transaction.Amount + transaction.Fee
	newFromBalance := fromAccount.Balance - totalAmount
	newToBalance := toAccount.Balance + transaction.Amount
	
	// Update from account balance
	if err := s.accountRepo.UpdateBalance(fromAccount.AccountNumber, newFromBalance); err != nil {
		return err
	}
	
	// Update to account balance
	if err := s.accountRepo.UpdateBalance(toAccount.AccountNumber, newToBalance); err != nil {
		// Rollback from account
		s.accountRepo.UpdateBalance(fromAccount.AccountNumber, fromAccount.Balance)
		return err
	}
	
	// Update transaction status
	return s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusCompleted)
}

func (s *transactionService) processDeposit(transaction *models.Transaction) error {
	// Get account
	account, err := s.accountRepo.GetByID(transaction.ToAccountID)
	if err != nil {
		return err
	}
	
	// Update balance
	newBalance := account.Balance + transaction.Amount
	
	if err := s.accountRepo.UpdateBalance(account.AccountNumber, newBalance); err != nil {
		return err
	}
	
	// Update transaction status
	return s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusCompleted)
}

func (s *transactionService) processWithdrawal(transaction *models.Transaction) error {
	// Get account
	account, err := s.accountRepo.GetByID(transaction.FromAccountID)
	if err != nil {
		return err
	}
	
	// Update balance
	totalAmount := transaction.Amount + transaction.Fee
	newBalance := account.Balance - totalAmount
	
	if err := s.accountRepo.UpdateBalance(account.AccountNumber, newBalance); err != nil {
		return err
	}
	
	// Update transaction status
	return s.transactionRepo.UpdateStatus(transaction.TransactionID, models.TransactionStatusCompleted)
}

func (s *transactionService) calculateTransferFee(amount int64) int64 {
	// 0.1% of transfer amount, minimum 100 cents ($1)
	fee := amount / 1000
	if fee < 100 {
		fee = 100
	}
	return fee
}

func (s *transactionService) calculateWithdrawalFee() int64 {
	// Fixed $2 fee for withdrawals
	return 200
}

func (s *transactionService) generateTransactionID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	
	return fmt.Sprintf("TXN%d%x", timestamp, randomBytes)
}
