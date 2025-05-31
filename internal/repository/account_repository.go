package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bank-api/internal/models"
	_ "github.com/lib/pq"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id int) (*models.Account, error)
	GetByAccountNumber(accountNumber string) (*models.Account, error)
	GetByCustomerID(customerID string) ([]*models.Account, error)
	GetAll(limit, offset int) ([]*models.Account, error)
	Update(id int, account *models.Account) error
	UpdateBalance(accountNumber string, balance int64) error
	UpdateStatus(id int, status string) error
	Delete(id int) error
	AccountExists(accountNumber string) (bool, error)
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) AccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (
			customer_id, account_number, iban, bic, account_type, currency,
			balance, available_balance, hold_amount, first_name, last_name,
			email, phone, date_of_birth, street, city, postal_code, country,
			state, hash_password, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23
		) RETURNING id`
	
	err := r.db.QueryRow(
		query,
		account.CustomerID, account.AccountNumber, account.IBAN, account.BIC,
		account.AccountType, account.Currency, account.Balance, account.AvailableBalance,
		account.HoldAmount, account.FirstName, account.LastName, account.Email,
		account.Phone, account.DateOfBirth, account.Address.Street, account.Address.City,
		account.Address.PostalCode, account.Address.Country, account.Address.State,
		account.HashPassword, account.Status, account.CreatedAt, account.UpdatedAt,
	).Scan(&account.ID)
	
	return err
}

func (r *PostgresAccountRepository) GetByID(id int) (*models.Account, error) {
	query := `
		SELECT id, customer_id, account_number, iban, bic, account_type, currency,
			   balance, available_balance, hold_amount, first_name, last_name,
			   email, phone, date_of_birth, street, city, postal_code, country,
			   state, hash_password, status, created_at, updated_at, last_login_at
		FROM accounts WHERE id = $1`
	
	account := &models.Account{}
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.CustomerID, &account.AccountNumber, &account.IBAN,
		&account.BIC, &account.AccountType, &account.Currency, &account.Balance,
		&account.AvailableBalance, &account.HoldAmount, &account.FirstName,
		&account.LastName, &account.Email, &account.Phone, &account.DateOfBirth,
		&account.Address.Street, &account.Address.City, &account.Address.PostalCode,
		&account.Address.Country, &account.Address.State, &account.HashPassword,
		&account.Status, &account.CreatedAt, &account.UpdatedAt, &lastLoginAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with id %d not found", id)
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		account.LastLoginAt = &lastLoginAt.Time
	}
	
	return account, nil
}

func (r *PostgresAccountRepository) GetByAccountNumber(accountNumber string) (*models.Account, error) {
	query := `
		SELECT id, customer_id, account_number, iban, bic, account_type, currency,
			   balance, available_balance, hold_amount, first_name, last_name,
			   email, phone, date_of_birth, street, city, postal_code, country,
			   state, hash_password, status, created_at, updated_at, last_login_at
		FROM accounts WHERE account_number = $1`
	
	account := &models.Account{}
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRow(query, accountNumber).Scan(
		&account.ID, &account.CustomerID, &account.AccountNumber, &account.IBAN,
		&account.BIC, &account.AccountType, &account.Currency, &account.Balance,
		&account.AvailableBalance, &account.HoldAmount, &account.FirstName,
		&account.LastName, &account.Email, &account.Phone, &account.DateOfBirth,
		&account.Address.Street, &account.Address.City, &account.Address.PostalCode,
		&account.Address.Country, &account.Address.State, &account.HashPassword,
		&account.Status, &account.CreatedAt, &account.UpdatedAt, &lastLoginAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with number %s not found", accountNumber)
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		account.LastLoginAt = &lastLoginAt.Time
	}
	
	return account, nil
}

func (r *PostgresAccountRepository) GetByCustomerID(customerID string) ([]*models.Account, error) {
	query := `
		SELECT id, customer_id, account_number, iban, bic, account_type, currency,
			   balance, available_balance, hold_amount, first_name, last_name,
			   email, phone, date_of_birth, street, city, postal_code, country,
			   state, hash_password, status, created_at, updated_at, last_login_at
		FROM accounts WHERE customer_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		var lastLoginAt sql.NullTime
		
		err := rows.Scan(
			&account.ID, &account.CustomerID, &account.AccountNumber, &account.IBAN,
			&account.BIC, &account.AccountType, &account.Currency, &account.Balance,
			&account.AvailableBalance, &account.HoldAmount, &account.FirstName,
			&account.LastName, &account.Email, &account.Phone, &account.DateOfBirth,
			&account.Address.Street, &account.Address.City, &account.Address.PostalCode,
			&account.Address.Country, &account.Address.State, &account.HashPassword,
			&account.Status, &account.CreatedAt, &account.UpdatedAt, &lastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		
		if lastLoginAt.Valid {
			account.LastLoginAt = &lastLoginAt.Time
		}
		
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}

func (r *PostgresAccountRepository) GetAll(limit, offset int) ([]*models.Account, error) {
	query := `
		SELECT id, customer_id, account_number, iban, bic, account_type, currency,
			   balance, available_balance, hold_amount, first_name, last_name,
			   email, phone, date_of_birth, street, city, postal_code, country,
			   state, hash_password, status, created_at, updated_at, last_login_at
		FROM accounts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var accounts []*models.Account
	for rows.Next() {
		account := &models.Account{}
		var lastLoginAt sql.NullTime
		
		err := rows.Scan(
			&account.ID, &account.CustomerID, &account.AccountNumber, &account.IBAN,
			&account.BIC, &account.AccountType, &account.Currency, &account.Balance,
			&account.AvailableBalance, &account.HoldAmount, &account.FirstName,
			&account.LastName, &account.Email, &account.Phone, &account.DateOfBirth,
			&account.Address.Street, &account.Address.City, &account.Address.PostalCode,
			&account.Address.Country, &account.Address.State, &account.HashPassword,
			&account.Status, &account.CreatedAt, &account.UpdatedAt, &lastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		
		if lastLoginAt.Valid {
			account.LastLoginAt = &lastLoginAt.Time
		}
		
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}

func (r *PostgresAccountRepository) Update(id int, account *models.Account) error {
	query := `
		UPDATE accounts SET
			first_name = $1, last_name = $2, email = $3, phone = $4,
			street = $5, city = $6, postal_code = $7, country = $8, state = $9,
			updated_at = $10
		WHERE id = $11`
	
	account.UpdatedAt = time.Now().UTC()
	
	_, err := r.db.Exec(
		query,
		account.FirstName, account.LastName, account.Email, account.Phone,
		account.Address.Street, account.Address.City, account.Address.PostalCode,
		account.Address.Country, account.Address.State, account.UpdatedAt, id,
	)
	
	return err
}

func (r *PostgresAccountRepository) UpdateBalance(accountNumber string, balance int64) error {
	query := `UPDATE accounts SET balance = $1, available_balance = $1, updated_at = $2 WHERE account_number = $3`
	_, err := r.db.Exec(query, balance, time.Now().UTC(), accountNumber)
	return err
}

func (r *PostgresAccountRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE accounts SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now().UTC(), id)
	return err
}

func (r *PostgresAccountRepository) Delete(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("account with id %d not found", id)
	}
	
	return nil
}

func (r *PostgresAccountRepository) AccountExists(accountNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE account_number = $1)`
	var exists bool
	err := r.db.QueryRow(query, accountNumber).Scan(&exists)
	return exists, err
}
