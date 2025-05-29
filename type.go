package main

import (
	"fmt"
	"math/rand"
)

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IBAN      string `json:"iban"`
	Balance   int64  `json:"balance"`
}

func generateIBAN() string {
	// Generate a random IBAN for TN (Tunisia)
	// Format: TN + 2 check digits + 20 digit account number
	accountNumber := rand.Int63n(1e18) // Generate up to 18-digit number (max for int64)
	return fmt.Sprintf("TN%02d%020d", rand.Intn(100), accountNumber)
}


func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000), // Random ID for demonstration
		FirstName: firstName,
		LastName:  lastName,
		IBAN:      generateIBAN(),
	}
}