package main

import (
	"flag"
	"fmt"
	"log"
)

// seedAccount creates and stores a new account with the given details
func seedAccount(store Storage, fname, lname, pw string) *Account {
	// Create a new account with the provided details
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	// Store the newly created account in the storage
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	// Print the account number of the newly created account
	fmt.Println("new account => ", acc.Number)

	// Return the created account
	return acc
}

// seedAccounts seeds the database with predefined accounts
func seedAccounts(s Storage) {
	// Add a specific account to the database
	seedAccount(s, "anthony", "GG", "hunter88888")
}

func main() {
	// Define a command-line flag to indicate whether to seed the database
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	// Create a new instance of the Postgres store
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the Postgres store (e.g., create tables, setup schema)
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Check if the seed flag is set; if so, seed the database with accounts
	if *seed {
		fmt.Println("seeding the database")
		seedAccounts(store)
	}

	// Create and run the API server
	server := NewAPIServer(":3000", store)
	server.Run()
}
