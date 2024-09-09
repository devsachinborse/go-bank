package main

import (
	"math/rand"        // Import the rand package for generating random numbers
	"time"             // Import the time package for time-related operations
	"golang.org/x/crypto/bcrypt" // Import bcrypt for password hashing and comparison
)

// LoginResponse represents the response structure for login requests
type LoginResponse struct {
	Number int64  `json:"number"` // Account number
	Token  string `json:"token"`  // JWT token for authentication
}

// LoginRequest represents the structure of a login request
type LoginRequest struct {
	Number   int64  `json:"number"`   // Account number
	Password string `json:"password"` // Password for authentication
}

// TransferRequest represents the structure of a transfer request
type TransferRequest struct {
	ToAccount int `json:"toAccount"` // Account number to which the amount is transferred
	Amount    int `json:"amount"`    // Amount to be transferred
}

// CreateAccountRequest represents the structure of a create account request
type CreateAccountRequest struct {
	FirstName string `json:"firstName"` // First name of the account holder
	LastName  string `json:"lastName"`  // Last name of the account holder
	Password  string `json:"password"`  // Password for the new account
}

// Account represents an individual account's details
type Account struct {
	ID                int       `json:"id"`                // Unique identifier for the account
	FirstName         string    `json:"firstName"`         // First name of the account holder
	LastName          string    `json:"lastName"`          // Last name of the account holder
	Number            int64     `json:"number"`            // Account number
	EncryptedPassword string    `json:"-"`                 // Encrypted password (not included in JSON serialization)
	Balance           int64     `json:"balance"`           // Account balance
	CreatedAt         time.Time `json:"createdAt"`         // Account creation timestamp
}

// ValidPassword checks if the provided password matches the stored encrypted password
func (a *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

// NewAccount creates a new account with a hashed password and random account number
func NewAccount(firstName, lastName, password string) (*Account, error) {
	// Hash the password using bcrypt
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // Return the error if password hashing fails
	}

	// Create and return a new Account object
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encpw),
		Number:            int64(rand.Intn(1000000)), // Generate a random account number
		CreatedAt:         time.Now().UTC(),          // Set the account creation time to the current UTC time
	}, nil
}
