package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bank-api/internal/repository"
	"github.com/bank-api/internal/utils"
)

type contextKey string

const (
	AccountNumberKey contextKey = "account_number"
	CustomerIDKey    contextKey = "customer_id"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(accountRepo repository.AccountRepository, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}
			
			// Check if it's a Bearer token
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}
			
			tokenString := tokenParts[1]
			
			// Verify token
			claims, err := utils.VerifyJWT(tokenString, secret)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
				return
			}
			
			// Verify account exists and is active
			account, err := accountRepo.GetByAccountNumber(claims.AccountNumber)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "Account not found")
				return
			}
			
			if !account.IsActive() {
				utils.WriteError(w, http.StatusUnauthorized, "Account is not active")
				return
			}
			
			// Add account info to context
			ctx := context.WithValue(r.Context(), AccountNumberKey, claims.AccountNumber)
			ctx = context.WithValue(ctx, CustomerIDKey, claims.CustomerID)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAccountNumberFromContext retrieves account number from request context
func GetAccountNumberFromContext(ctx context.Context) (string, bool) {
	accountNumber, ok := ctx.Value(AccountNumberKey).(string)
	return accountNumber, ok
}

// GetCustomerIDFromContext retrieves customer ID from request context
func GetCustomerIDFromContext(ctx context.Context) (string, bool) {
	customerID, ok := ctx.Value(CustomerIDKey).(string)
	return customerID, ok
}
