package routes

import (
	"database/sql"
	"net/http"

	"github.com/bank-api/internal/api/handlers"
	"github.com/bank-api/internal/api/middleware"
	"github.com/bank-api/internal/config"
	"github.com/bank-api/internal/repository"
	"github.com/bank-api/internal/services"
	"github.com/gorilla/mux"
)

type Router struct {
	accountHandler     *handlers.AccountHandler
	authHandler        *handlers.AuthHandler
	transactionHandler *handlers.TransactionHandler
	authMiddleware     func(http.Handler) http.Handler
}

func NewRouter(db *sql.DB, cfg *config.Config) *Router {
	// Initialize repositories
	accountRepo := repository.NewPostgresAccountRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)
	
	// Initialize services
	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo)
	
	// Initialize handlers
	accountHandler := handlers.NewAccountHandler(accountService)
	authHandler := handlers.NewAuthHandler(accountService, cfg.JWT.Secret, cfg.JWT.ExpiresIn)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	
	// Initialize middleware
	authMiddleware := middleware.JWTAuthMiddleware(accountRepo, cfg.JWT.Secret)
	
	return &Router{
		accountHandler:     accountHandler,
		authHandler:        authHandler,
		transactionHandler: transactionHandler,
		authMiddleware:     authMiddleware,
	}
}

func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	
	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.ContentTypeMiddleware)
	
	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Health check
	api.HandleFunc("/health", r.healthCheck).Methods("GET")
	
	// Authentication routes (no auth required)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", r.authHandler.Login).Methods("POST")
	auth.HandleFunc("/logout", r.authHandler.Logout).Methods("POST")
	auth.HandleFunc("/refresh", r.authHandler.RefreshToken).Methods("POST")
	
	// Account routes
	accounts := api.PathPrefix("/accounts").Subrouter()
	
	// Public account routes (no auth required)
	accounts.HandleFunc("", r.accountHandler.CreateAccount).Methods("POST")
	
	// Protected account routes (auth required)
	protectedAccounts := accounts.PathPrefix("").Subrouter()
	protectedAccounts.Use(r.authMiddleware)
	protectedAccounts.HandleFunc("", r.accountHandler.GetAccounts).Methods("GET")
	protectedAccounts.HandleFunc("/{id:[0-9]+}", r.accountHandler.GetAccount).Methods("GET")
	protectedAccounts.HandleFunc("/{id:[0-9]+}", r.accountHandler.UpdateAccount).Methods("PUT")
	protectedAccounts.HandleFunc("/{id:[0-9]+}", r.accountHandler.DeleteAccount).Methods("DELETE")
	protectedAccounts.HandleFunc("/{id:[0-9]+}/status", r.accountHandler.UpdateAccountStatus).Methods("PATCH")
	protectedAccounts.HandleFunc("/{accountNumber}/balance", r.accountHandler.GetAccountBalance).Methods("GET")
	
	// Transaction routes (all require auth)
	transactions := api.PathPrefix("/transactions").Subrouter()
	transactions.Use(r.authMiddleware)
	transactions.HandleFunc("/transfer", r.transactionHandler.Transfer).Methods("POST")
	transactions.HandleFunc("/deposit", r.transactionHandler.Deposit).Methods("POST")
	transactions.HandleFunc("/withdraw", r.transactionHandler.Withdraw).Methods("POST")
	transactions.HandleFunc("/history", r.transactionHandler.GetTransactionHistory).Methods("GET")
	transactions.HandleFunc("/{transactionId}", r.transactionHandler.GetTransaction).Methods("GET")
	
	return router
}

func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK", "message": "Bank API is running"}`))
}
