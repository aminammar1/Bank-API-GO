package repository

import (
	"database/sql"
	"fmt"

	"github.com/bank-api/internal/config"
	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

// InitializeDatabase creates the necessary tables
func InitializeDatabase(db *sql.DB) error {
	// Drop existing tables to ensure clean schema
	if err := dropTables(db); err != nil {
		return fmt.Errorf("failed to drop existing tables: %w", err)
	}
	
	if err := createAccountsTable(db); err != nil {
		return fmt.Errorf("failed to create accounts table: %w", err)
	}
	
	if err := createTransactionsTable(db); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	
	return nil
}

func dropTables(db *sql.DB) error {
	queries := []string{
		"DROP TABLE IF EXISTS transactions CASCADE;",
		"DROP TABLE IF EXISTS accounts CASCADE;",
	}
	
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute drop query %s: %w", query, err)
		}
	}
	
	return nil
}

func createAccountsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		customer_id VARCHAR(50) UNIQUE NOT NULL,
		account_number VARCHAR(20) UNIQUE NOT NULL,
		iban VARCHAR(34) UNIQUE NOT NULL,
		bic VARCHAR(11) NOT NULL,
		account_type VARCHAR(20) NOT NULL,
		currency VARCHAR(3) NOT NULL,
		balance BIGINT NOT NULL DEFAULT 0,
		available_balance BIGINT NOT NULL DEFAULT 0,
		hold_amount BIGINT NOT NULL DEFAULT 0,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(20) NOT NULL,
		date_of_birth DATE NOT NULL,
		street VARCHAR(255),
		city VARCHAR(100),
		postal_code VARCHAR(20),
		country VARCHAR(100),
		state VARCHAR(100),
		hash_password VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
		last_login_at TIMESTAMP WITH TIME ZONE
	);
	
	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_accounts_customer_id ON accounts(customer_id);
	CREATE INDEX IF NOT EXISTS idx_accounts_account_number ON accounts(account_number);
	CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);
	CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status);
	`
	
	_, err := db.Exec(query)
	return err
}

func createTransactionsTable(db *sql.DB) error {
	query := `	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		transaction_id VARCHAR(50) UNIQUE NOT NULL,
		from_account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
		to_account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
		from_account_number VARCHAR(20),
		to_account_number VARCHAR(20),
		amount BIGINT NOT NULL,
		currency VARCHAR(3) NOT NULL,
		exchange_rate DECIMAL(10,6) DEFAULT 1.0,
		converted_amount BIGINT NOT NULL,
		transaction_type VARCHAR(20) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		description TEXT,
		reference VARCHAR(100),
		fee BIGINT NOT NULL DEFAULT 0,
		processed_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
		failure_reason TEXT,
		
		-- Constraints
		CONSTRAINT chk_amount_positive CHECK (amount > 0),
		CONSTRAINT chk_valid_transaction_type CHECK (
			transaction_type IN ('TRANSFER', 'DEPOSIT', 'WITHDRAWAL', 'PAYMENT', 'FEE', 'INTEREST')
		),
		CONSTRAINT chk_valid_status CHECK (
			status IN ('PENDING', 'COMPLETED', 'FAILED', 'CANCELLED')
		),		CONSTRAINT chk_valid_currency CHECK (
			currency IN ('TND', 'EUR', 'USD')
		)
	);
	
	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_transactions_transaction_id ON transactions(transaction_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_from_account_id ON transactions(from_account_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_to_account_id ON transactions(to_account_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_from_account_number ON transactions(from_account_number);
	CREATE INDEX IF NOT EXISTS idx_transactions_to_account_number ON transactions(to_account_number);
	CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(transaction_type);
	CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
	`
	
	_, err := db.Exec(query)
	return err
}
