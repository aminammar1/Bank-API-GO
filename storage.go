package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
	GetAccountByNumber(int) (*Account, error)
	UpdateAccount(int, *Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}


func NewPostgresStore() (*PostgresStore, error) {
	connStr := "host=localhost port=5433 user=bankgo dbname=bankdb password=testbank sslmode=disable"
	
	db, err := sql.Open("postgres", connStr)
	
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, err
	
	}
	
	return &PostgresStore{db: db}, nil
}


func (s *PostgresStore) init() error {
	return s.createAccountTable()
}



func (s *PostgresStore) createAccountTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		hashpassword VARCHAR(60) NOT NULL,
		number BIGINT UNIQUE NOT NULL,
		balance BIGINT NOT NULL DEFAULT 0, 
		CreatedAt TIMESTAMP 
	);
	`
	_, err := s.db.Exec(query)
	
	return err
}



func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `
	INSERT INTO accounts (first_name, last_name, number, hashpassword, balance, CreatedAt)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query, 
		account.FirstName,
		account.LastName,
		account.Number,
		account.HashPassword,
		account.Balance,
		account.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}


func (s *PostgresStore) UpdateAccount(id int, account *Account) error {
	query := `
	UPDATE accounts
	SET first_name = $1, last_name = $2, number = $3, hashpassword = $4, balance = $5
	WHERE id = $6`
	_, err := s.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.HashPassword,
		account.Balance,
		id)
		
		if err != nil {
			return  err
		}
	return nil
	
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int)(*Account, error) {
	rows , err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil , fmt.Errorf("account with id %d not found", id)
}




func (s *PostgresStore) GetAccountByNumber (number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE number = $1", number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account with iban %s not found", number)
}


func (s *PostgresStore) GetAllAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}




func scanIntoAccount (rows *sql.Rows) (*Account, error) {
	account := new (Account)
		
	err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.HashPassword,
			&account.Number,
			&account.Balance,
			&account.CreatedAt)
	
	return account, err
}

