package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bank-api/internal/api/routes"
	"github.com/bank-api/internal/config"
	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/repository"
)

const (
	testDBName = "bankdb" // Use main database for now
)

var (
	testRouter *routes.Router
	testDB     *sql.DB
	testConfig *config.Config
)

func TestMain(m *testing.M) {
	// Setup test environment
	setup()
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	teardown()
	
	os.Exit(code)
}

func setup() {
	// Load test configuration
	testConfig = &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5433",
			User:     "bankgo",
			Password: "testbank",
			DBName:   testDBName,
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			Secret:    "test-secret-key",
			ExpiresIn: 24 * time.Hour,
			Issuer:    "bank-api-test",
		},
	}
	
	// Create test database connection
	var err error
	testDB, err = repository.NewPostgresDB(&testConfig.Database)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}
	
	// Initialize test database
	if err := repository.InitializeDatabase(testDB); err != nil {
		panic(fmt.Sprintf("Failed to initialize test database: %v", err))
	}
	
	// Create test router
	testRouter = routes.NewRouter(testDB, testConfig)
}

func teardown() {
	if testDB != nil {
		// Clean up test data
		testDB.Exec("TRUNCATE TABLE transactions CASCADE")
		testDB.Exec("TRUNCATE TABLE accounts CASCADE")
		testDB.Close()
	}
}

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	expected := `{"status": "OK", "message": "Bank API is running"}`
	if rr.Body.String() != expected {
		t.Errorf("Health check returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateAccount(t *testing.T) {
	createAccountReq := models.CreateAccountRequest{
		FirstName:   "Mohamed",
		LastName:    "Ben Ahmed",
		Email:       "mohamed.benahmed@example.tn",
		Phone:       "+21612345678",
		Password:    "motdepasse123",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		AccountType: models.AccountTypeChecking,
		Currency:    models.CurrencyTND,
		Address: models.Address{
			Street:     "Avenue Habib Bourguiba 123",
			City:       "Tunis",
			PostalCode: "1001",
			Country:    "Tunisia",
			State:      "Tunis",
		},
	}
	
	jsonData, _ := json.Marshal(createAccountReq)
	req, err := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Create account returned wrong status code: got %v want %v", status, http.StatusCreated)
		t.Errorf("Response body: %s", rr.Body.String())
	}
	
	var response models.SuccessResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}
	
	if response.Message != "Account created successfully" {
		t.Errorf("Unexpected response message: got %v", response.Message)
	}
}

func TestLogin(t *testing.T) {
	// First create an account
	account := createTestAccount(t)
		// Now test login
	loginReq := models.LoginRequest{
		AccountNumber: account.AccountNumber,
		Password:      "motdepasse123",
	}
	
	jsonData, _ := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Login returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Response body: %s", rr.Body.String())
	}
	
	var response models.SuccessResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}
	
	// Extract login response data
	loginData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Failed to extract login data from response")
	}
	
	if loginData["token"] == nil {
		t.Error("Login response should contain a token")
	}
	
	if loginData["account_number"] != account.AccountNumber {
		t.Errorf("Login response account number mismatch: got %v want %v", 
			loginData["account_number"], account.AccountNumber)
	}
}

func TestTransfer(t *testing.T) {
	// Create two test accounts
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)
	
	// Login to get token
	token := loginAndGetToken(t, fromAccount.AccountNumber)
		// First deposit money into the from account
	depositReq := models.DepositRequest{
		AccountNumber: fromAccount.AccountNumber,
		Amount:        100000, // 100 TND (in millimes)
		Currency:      models.CurrencyTND,
		Description:   "Dépôt de test",
	}
	
	jsonData, _ := json.Marshal(depositReq)
	req, err := http.NewRequest("POST", "/api/v1/transactions/deposit", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Deposit returned wrong status code: got %v want %v", status, http.StatusCreated)
		t.Logf("Deposit response body: %s", rr.Body.String())
		return // Skip transfer test if deposit fails
	}
		// Now test transfer
	transferReq := models.TransferRequest{
		FromAccountNumber: fromAccount.AccountNumber,
		ToAccountNumber:   toAccount.AccountNumber,
		Amount:            50000, // 50 TND (in millimes)
		Currency:          models.CurrencyTND,
		Description:       "Virement de test",
	}
	
	jsonData, _ = json.Marshal(transferReq)
	req, err = http.NewRequest("POST", "/api/v1/transactions/transfer", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Transfer returned wrong status code: got %v want %v", status, http.StatusCreated)
		t.Errorf("Response body: %s", rr.Body.String())
	}
}

func TestGetTransactionHistory(t *testing.T) {
	// Create test account and login
	account := createTestAccount(t)
	token := loginAndGetToken(t, account.AccountNumber)
		// Create some transactions first (deposit)
	depositReq := models.DepositRequest{
		AccountNumber: account.AccountNumber,
		Amount:        100000, // 100 TND
		Currency:      models.CurrencyTND,
		Description:   "Dépôt de test",
	}
	
	jsonData, _ := json.Marshal(depositReq)
	req, err := http.NewRequest("POST", "/api/v1/transactions/deposit", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	// Now get transaction history
	req, err = http.NewRequest("GET", "/api/v1/transactions/history", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Get transaction history returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Response body: %s", rr.Body.String())
	}
}

// Helper functions

func createTestAccount(t *testing.T) *models.Account {
	createAccountReq := models.CreateAccountRequest{
		FirstName:   fmt.Sprintf("Ahmed-%d", time.Now().UnixNano()),
		LastName:    "Trabelsi",
		Email:       fmt.Sprintf("ahmed.trabelsi-%d@example.tn", time.Now().UnixNano()),
		Phone:       "+21625123456",
		Password:    "motdepasse123",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		AccountType: models.AccountTypeChecking,
		Currency:    models.CurrencyTND,
		Address: models.Address{
			Street:     "Rue de la République 45",
			City:       "Sfax",
			PostalCode: "3000",
			Country:    "Tunisia",
			State:      "Sfax",
		},
	}
	
	jsonData, _ := json.Marshal(createAccountReq)
	req, err := http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Failed to create test account: status %v", status)
	}
	
	var response models.SuccessResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal create account response:", err)
	}
	
	// Convert response data to account
	accountData, _ := json.Marshal(response.Data)
	var account models.Account
	if err := json.Unmarshal(accountData, &account); err != nil {
		t.Fatal("Failed to unmarshal account data:", err)
	}
	
	return &account
}

func loginAndGetToken(t *testing.T, accountNumber string) string {
	loginReq := models.LoginRequest{
		AccountNumber: accountNumber,
		Password:      "motdepasse123",
	}
	
	jsonData, _ := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := testRouter.SetupRoutes()
	handler.ServeHTTP(rr, req)
	
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Failed to login: status %v", status)
	}
	
	var response models.SuccessResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to unmarshal login response:", err)
	}
	
	loginData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Failed to extract login data from response")
	}
	
	token, ok := loginData["token"].(string)
	if !ok {
		t.Fatal("Failed to extract token from login response")
	}
	
	return token
}
