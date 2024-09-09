package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert" // Import the testify package for assertions
)

// TestNewAccount tests the NewAccount function for creating a new account
func TestNewAccount(t *testing.T) {
	// Create a new account with given first name, last name, and password
	acc, err := NewAccount("a", "b", "hunter")

	// Assert that there is no error during account creation
	assert.Nil(t, err)

	// Print the created account details for debugging purposes
	fmt.Printf("%+v\n", acc)
}