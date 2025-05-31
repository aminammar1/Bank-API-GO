package handlers

import (
	"net/http"
	"time"

	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/services"
	"github.com/bank-api/internal/utils"
)

type AuthHandler struct {
	accountService services.AccountService
	jwtSecret      string
	jwtExpiresIn   time.Duration
}

func NewAuthHandler(accountService services.AccountService, jwtSecret string, jwtExpiresIn time.Duration) *AuthHandler {
	return &AuthHandler{
		accountService: accountService,
		jwtSecret:      jwtSecret,
		jwtExpiresIn:   jwtExpiresIn,
	}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	// Validate request
	if req.AccountNumber == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Account number and password are required")
		return
	}
	
	// Authenticate account
	account, err := h.accountService.AuthenticateAccount(req.AccountNumber, req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	
	// Generate JWT token
	token, err := utils.GenerateJWT(account, h.jwtSecret, h.jwtExpiresIn)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	
	// Create response
	response := models.LoginResponse{
		Token:         token,
		AccountNumber: account.AccountNumber,
		CustomerID:    account.CustomerID,
		ExpiresAt:     time.Now().Add(h.jwtExpiresIn),
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Login successful", response)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return a success response
	utils.WriteSuccess(w, http.StatusOK, "Logout successful", nil)
}

// RefreshToken handles POST /auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, "Authorization header required")
		return
	}
	
	// Extract token
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
		return
	}
	
	// Verify current token
	claims, err := utils.VerifyJWT(tokenString, h.jwtSecret)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	
	// Get account
	account, err := h.accountService.GetAccountByAccountNumber(claims.AccountNumber)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found")
		return
	}
	
	// Generate new token
	newToken, err := utils.GenerateJWT(account, h.jwtSecret, h.jwtExpiresIn)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	
	// Create response
	response := models.LoginResponse{
		Token:         newToken,
		AccountNumber: account.AccountNumber,
		CustomerID:    account.CustomerID,
		ExpiresAt:     time.Now().Add(h.jwtExpiresIn),
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Token refreshed successfully", response)
}
