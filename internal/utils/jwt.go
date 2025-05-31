package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/bank-api/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	AccountNumber string `json:"account_number"`
	CustomerID    string `json:"customer_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for the given account
func GenerateJWT(account *models.Account, secret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		AccountNumber: account.AccountNumber,
		CustomerID:    account.CustomerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bank-api",
			Subject:   account.CustomerID,
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// VerifyJWT verifies and parses a JWT token
func VerifyJWT(tokenString string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

// GetJWTSecret gets the JWT secret from environment variable
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return secret
}
