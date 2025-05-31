package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bank-api/internal/api/middleware"
	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/services"
	"github.com/bank-api/internal/utils"
	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Transfer handles POST /transactions/transfer
func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	// Get account number from context (set by auth middleware)
	accountNumber, ok := middleware.GetAccountNumberFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found in context")
		return
	}
	
	// Ensure the user can only transfer from their own account
	if req.FromAccountNumber != accountNumber {
		utils.WriteError(w, http.StatusForbidden, "You can only transfer from your own account")
		return
	}
	
	transaction, err := h.transactionService.Transfer(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusCreated, "Transfer initiated successfully", transaction)
}

// Deposit handles POST /transactions/deposit
func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req models.DepositRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	// Get account number from context (set by auth middleware)
	accountNumber, ok := middleware.GetAccountNumberFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found in context")
		return
	}
	
	// Ensure the user can only deposit to their own account
	if req.AccountNumber != accountNumber {
		utils.WriteError(w, http.StatusForbidden, "You can only deposit to your own account")
		return
	}
	
	transaction, err := h.transactionService.Deposit(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusCreated, "Deposit completed successfully", transaction)
}

// Withdraw handles POST /transactions/withdraw
func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req models.WithdrawalRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	// Get account number from context (set by auth middleware)
	accountNumber, ok := middleware.GetAccountNumberFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found in context")
		return
	}
	
	// Ensure the user can only withdraw from their own account
	if req.AccountNumber != accountNumber {
		utils.WriteError(w, http.StatusForbidden, "You can only withdraw from your own account")
		return
	}
	
	transaction, err := h.transactionService.Withdraw(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusCreated, "Withdrawal completed successfully", transaction)
}

// GetTransaction handles GET /transactions/{transactionId}
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, exists := vars["transactionId"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Transaction ID is required")
		return
	}
	
	transaction, err := h.transactionService.GetTransaction(transactionID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	
	// Get account number from context to ensure user can only see their own transactions
	accountNumber, ok := middleware.GetAccountNumberFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found in context")
		return
	}
	
	// Check if the user is authorized to see this transaction
	if transaction.FromAccountNumber != accountNumber && transaction.ToAccountNumber != accountNumber {
		utils.WriteError(w, http.StatusForbidden, "You are not authorized to view this transaction")
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// GetTransactionHistory handles GET /transactions/history
func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	// Get account number from context
	accountNumber, ok := middleware.GetAccountNumberFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Account not found in context")
		return
	}
	
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	transactionType := r.URL.Query().Get("type")
	
	req := models.TransactionHistoryRequest{
		AccountNumber: accountNumber,
		Type:          transactionType,
	}
	
	// Parse limit
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	
	// Parse offset
	if offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}
	
	// Parse start date
	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = startDate
		}
	}
	
	// Parse end date
	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = endDate
		}
	}
	
	transactions, err := h.transactionService.GetTransactionHistory(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Transaction history retrieved successfully", transactions)
}
