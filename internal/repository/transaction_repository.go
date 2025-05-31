package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bank-api/internal/models"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id int) (*models.Transaction, error)
	GetByTransactionID(transactionID string) (*models.Transaction, error)
	GetByAccountNumber(accountNumber string, limit, offset int) ([]*models.Transaction, error)
	GetByDateRange(accountNumber string, startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error)
	UpdateStatus(transactionID string, status string) error
	GetPendingTransactions() ([]*models.Transaction, error)
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) TransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			transaction_id, from_account_id, to_account_id, from_account_number,
			to_account_number, amount, currency, exchange_rate, converted_amount,
			transaction_type, status, description, reference, fee, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING id`
	
	// Handle nullable foreign key references
	var fromAccountID, toAccountID interface{}
	if transaction.FromAccountID != 0 {
		fromAccountID = transaction.FromAccountID
	}
	if transaction.ToAccountID != 0 {
		toAccountID = transaction.ToAccountID
	}
	
	err := r.db.QueryRow(
		query,
		transaction.TransactionID, fromAccountID, toAccountID,
		transaction.FromAccountNumber, transaction.ToAccountNumber, transaction.Amount,
		transaction.Currency, transaction.ExchangeRate, transaction.ConvertedAmount,
		transaction.TransactionType, transaction.Status, transaction.Description,
		transaction.Reference, transaction.Fee, transaction.CreatedAt, transaction.UpdatedAt,
	).Scan(&transaction.ID)
	
	return err
}

func (r *PostgresTransactionRepository) GetByID(id int) (*models.Transaction, error) {
	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, from_account_number,
			   to_account_number, amount, currency, exchange_rate, converted_amount,
			   transaction_type, status, description, reference, fee, processed_at,
			   created_at, updated_at, failure_reason
		FROM transactions WHERE id = $1`
	
	transaction := &models.Transaction{}
	var fromAccountID, toAccountID sql.NullInt32
	var processedAt sql.NullTime
	var failureReason sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&transaction.ID, &transaction.TransactionID, &fromAccountID,
		&toAccountID, &transaction.FromAccountNumber, &transaction.ToAccountNumber,
		&transaction.Amount, &transaction.Currency, &transaction.ExchangeRate,
		&transaction.ConvertedAmount, &transaction.TransactionType, &transaction.Status,
		&transaction.Description, &transaction.Reference, &transaction.Fee,
		&processedAt, &transaction.CreatedAt, &transaction.UpdatedAt, &failureReason,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", id)
		}
		return nil, err
	}
	
	// Handle nullable fields
	if fromAccountID.Valid {
		transaction.FromAccountID = int(fromAccountID.Int32)
	}
	if toAccountID.Valid {
		transaction.ToAccountID = int(toAccountID.Int32)
	}
	if processedAt.Valid {
		transaction.ProcessedAt = &processedAt.Time
	}
	if failureReason.Valid {
		transaction.FailureReason = failureReason.String
	}
	
	return transaction, nil
}

func (r *PostgresTransactionRepository) GetByTransactionID(transactionID string) (*models.Transaction, error) {
	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, from_account_number,
			   to_account_number, amount, currency, exchange_rate, converted_amount,
			   transaction_type, status, description, reference, fee, processed_at,
			   created_at, updated_at, failure_reason
		FROM transactions WHERE transaction_id = $1`
	
	transaction := &models.Transaction{}
	var fromAccountID, toAccountID sql.NullInt32
	var processedAt sql.NullTime
	var failureReason sql.NullString
	
	err := r.db.QueryRow(query, transactionID).Scan(
		&transaction.ID, &transaction.TransactionID, &fromAccountID,
		&toAccountID, &transaction.FromAccountNumber, &transaction.ToAccountNumber,
		&transaction.Amount, &transaction.Currency, &transaction.ExchangeRate,
		&transaction.ConvertedAmount, &transaction.TransactionType, &transaction.Status,
		&transaction.Description, &transaction.Reference, &transaction.Fee,
		&processedAt, &transaction.CreatedAt, &transaction.UpdatedAt, &failureReason,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with ID %s not found", transactionID)
		}
		return nil, err
	}
	
	// Handle nullable fields
	if fromAccountID.Valid {
		transaction.FromAccountID = int(fromAccountID.Int32)
	}
	if toAccountID.Valid {
		transaction.ToAccountID = int(toAccountID.Int32)
	}
	if processedAt.Valid {
		transaction.ProcessedAt = &processedAt.Time
	}
	if failureReason.Valid {
		transaction.FailureReason = failureReason.String
	}
	
	return transaction, nil
}

func (r *PostgresTransactionRepository) GetByAccountNumber(accountNumber string, limit, offset int) ([]*models.Transaction, error) {
	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, from_account_number,
			   to_account_number, amount, currency, exchange_rate, converted_amount,
			   transaction_type, status, description, reference, fee, processed_at,
			   created_at, updated_at, failure_reason
		FROM transactions 
		WHERE from_account_number = $1 OR to_account_number = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	
	rows, err := r.db.Query(query, accountNumber, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		var fromAccountID, toAccountID sql.NullInt32
		var processedAt sql.NullTime
		var failureReason sql.NullString
		
		err := rows.Scan(
			&transaction.ID, &transaction.TransactionID, &fromAccountID,
			&toAccountID, &transaction.FromAccountNumber, &transaction.ToAccountNumber,
			&transaction.Amount, &transaction.Currency, &transaction.ExchangeRate,
			&transaction.ConvertedAmount, &transaction.TransactionType, &transaction.Status,
			&transaction.Description, &transaction.Reference, &transaction.Fee,
			&processedAt, &transaction.CreatedAt, &transaction.UpdatedAt, &failureReason,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if fromAccountID.Valid {
			transaction.FromAccountID = int(fromAccountID.Int32)
		}
		if toAccountID.Valid {
			transaction.ToAccountID = int(toAccountID.Int32)
		}
		if processedAt.Valid {
			transaction.ProcessedAt = &processedAt.Time
		}
		if failureReason.Valid {
			transaction.FailureReason = failureReason.String
		}
		
		transactions = append(transactions, transaction)
	}
	
	return transactions, nil
}

func (r *PostgresTransactionRepository) GetByDateRange(accountNumber string, startDate, endDate time.Time, limit, offset int) ([]*models.Transaction, error) {
	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, from_account_number,
			   to_account_number, amount, currency, exchange_rate, converted_amount,
			   transaction_type, status, description, reference, fee, processed_at,
			   created_at, updated_at, failure_reason
		FROM transactions 
		WHERE (from_account_number = $1 OR to_account_number = $1)
		  AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC 
		LIMIT $4 OFFSET $5`
	
	rows, err := r.db.Query(query, accountNumber, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		var fromAccountID, toAccountID sql.NullInt32
		var processedAt sql.NullTime
		var failureReason sql.NullString
		
		err := rows.Scan(
			&transaction.ID, &transaction.TransactionID, &fromAccountID,
			&toAccountID, &transaction.FromAccountNumber, &transaction.ToAccountNumber,
			&transaction.Amount, &transaction.Currency, &transaction.ExchangeRate,
			&transaction.ConvertedAmount, &transaction.TransactionType, &transaction.Status,
			&transaction.Description, &transaction.Reference, &transaction.Fee,
			&processedAt, &transaction.CreatedAt, &transaction.UpdatedAt, &failureReason,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if fromAccountID.Valid {
			transaction.FromAccountID = int(fromAccountID.Int32)
		}
		if toAccountID.Valid {
			transaction.ToAccountID = int(toAccountID.Int32)
		}
		if processedAt.Valid {
			transaction.ProcessedAt = &processedAt.Time
		}
		if failureReason.Valid {
			transaction.FailureReason = failureReason.String
		}
		
		transactions = append(transactions, transaction)
	}
	
	return transactions, nil
}

func (r *PostgresTransactionRepository) UpdateStatus(transactionID string, status string) error {
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE transaction_id = $3`
	
	now := time.Now().UTC()
	_, err := r.db.Exec(query, status, now, transactionID)
	
	if status == models.TransactionStatusCompleted {
		processedQuery := `UPDATE transactions SET processed_at = $1 WHERE transaction_id = $2`
		_, err = r.db.Exec(processedQuery, now, transactionID)
	}
	
	return err
}

func (r *PostgresTransactionRepository) GetPendingTransactions() ([]*models.Transaction, error) {
	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, from_account_number,
			   to_account_number, amount, currency, exchange_rate, converted_amount,
			   transaction_type, status, description, reference, fee, processed_at,
			   created_at, updated_at, failure_reason
		FROM transactions 
		WHERE status = $1
		ORDER BY created_at ASC`
	
	rows, err := r.db.Query(query, models.TransactionStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		var fromAccountID, toAccountID sql.NullInt32
		var processedAt sql.NullTime
		var failureReason sql.NullString
		
		err := rows.Scan(
			&transaction.ID, &transaction.TransactionID, &fromAccountID,
			&toAccountID, &transaction.FromAccountNumber, &transaction.ToAccountNumber,
			&transaction.Amount, &transaction.Currency, &transaction.ExchangeRate,
			&transaction.ConvertedAmount, &transaction.TransactionType, &transaction.Status,
			&transaction.Description, &transaction.Reference, &transaction.Fee,
			&processedAt, &transaction.CreatedAt, &transaction.UpdatedAt, &failureReason,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable fields
		if fromAccountID.Valid {
			transaction.FromAccountID = int(fromAccountID.Int32)
		}
		if toAccountID.Valid {
			transaction.ToAccountID = int(toAccountID.Int32)
		}
		if processedAt.Valid {
			transaction.ProcessedAt = &processedAt.Time
		}
		if failureReason.Valid {
			transaction.FailureReason = failureReason.String
		}
		
		transactions = append(transactions, transaction)
	}
	
	return transactions, nil
}
