package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// APIServer struct holds the server's listening address and the storage interface
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates and returns a new APIServer instance with the given address and storage
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the HTTP server with all defined routes
func (s *APIServer) Run() {
	// Create a new router
	router := mux.NewRouter()

	// Define routes and their handlers
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	// Log the server start message
	log.Println("JSON API server running on port: ", s.listenAddr)

	// Start the HTTP server
	http.ListenAndServe(s.listenAddr, router)
}

// handleLogin handles the login request, verifies the credentials, and returns a JWT token
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// Only allow POST method
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	// Decode the login request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Retrieve the account by account number
	acc, err := s.store.GetAccountByNumber(int(req.Number))
	if err != nil {
		return err
	}

	// Verify the provided password
	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	// Create a JWT token for the authenticated account
	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	// Send the token and account number as the response
	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

// handleAccount handles both GET and POST requests for accounts
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Handle GET and POST requests
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	// Return an error if the method is not allowed
	return fmt.Errorf("method not allowed %s", r.Method)
}

// handleGetAccount retrieves all accounts and sends them as a response
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// Retrieve all accounts from the storage
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	// Send the accounts as JSON response
	return WriteJSON(w, http.StatusOK, accounts)
}

// handleGetAccountByID retrieves an account by ID or deletes it if DELETE method is used
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	// Handle GET method for fetching an account by ID
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	// Handle DELETE method for deleting an account
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	// Return an error if the method is not allowed
	return fmt.Errorf("method not allowed %s", r.Method)
}

// handleCreateAccount creates a new account and stores it
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// Decode the request body to create an account
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// Create a new account
	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}
	// Store the account in the storage
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	// Send the created account as JSON response
	return WriteJSON(w, http.StatusOK, account)
}

// handleDeleteAccount deletes an account by its ID
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Get the account ID from the URL
	id, err := getID(r)
	if err != nil {
		return err
	}

	// Delete the account from the storage
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	// Send a confirmation response
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// handleTransfer handles the transfer request and sends the transfer details as the response
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	// Decode the transfer request body
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	// Send the transfer details as JSON response
	return WriteJSON(w, http.StatusOK, transferReq)
}

// WriteJSON sends a JSON response with the specified status and value
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// createJWT creates a JWT token for the given account
func createJWT(account *Account) (string, error) {
	// Define the JWT claims
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	// Retrieve the secret key from environment variables
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	return token.SignedString([]byte(secret))
}

// permissionDenied sends a permission denied response
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

// withJWTAuth is a middleware that checks JWT authentication for the given handler function
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		// Retrieve the token from the request header
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		// Get the user ID from the request
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		// Retrieve the account associated with the user ID
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		// Validate the token claims against the account number
		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		// Call the next handler function
		handlerFunc(w, r)
	}
}

// validateJWT parses and validates a JWT token
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	// Parse the token and verify the signing method
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for token verification
		return []byte(secret), nil
	})
}

// apiFunc is a type alias for functions that handle HTTP requests and return an error
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents an error response
type ApiError struct {
	Error string `json:"error"`
}

// makeHTTPHandleFunc wraps an apiFunc to handle HTTP requests and send error responses
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
