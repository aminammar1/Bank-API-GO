package handlers

import (
	"net/http"
	"strconv"

	"github.com/bank-api/internal/models"
	"github.com/bank-api/internal/services"
	"github.com/bank-api/internal/utils"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	account, err := h.accountService.CreateAccount(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusCreated, "Account created successfully", account)
}

// GetAccount handles GET /accounts/{id}
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Account ID is required")
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}
	
	account, err := h.accountService.GetAccountByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Account retrieved successfully", account)
}

// GetAccounts handles GET /accounts
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	
	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	
	accounts, err := h.accountService.GetAllAccounts(limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Accounts retrieved successfully", accounts)
}

// UpdateAccount handles PUT /accounts/{id}
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Account ID is required")
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}
	
	var req models.UpdateAccountRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := h.accountService.UpdateAccount(id, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Account updated successfully", nil)
}

// DeleteAccount handles DELETE /accounts/{id}
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Account ID is required")
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}
	
	if err := h.accountService.DeleteAccount(id); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Account deleted successfully", nil)
}

// GetAccountBalance handles GET /accounts/{accountNumber}/balance
func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountNumber, exists := vars["accountNumber"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Account number is required")
		return
	}
	
	balance, err := h.accountService.GetAccountBalance(accountNumber)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Balance retrieved successfully", balance)
}

// UpdateAccountStatus handles PATCH /accounts/{id}/status
func (h *AccountHandler) UpdateAccountStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		utils.WriteError(w, http.StatusBadRequest, "Account ID is required")
		return
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}
	
	var req struct {
		Status string `json:"status"`
	}
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	
	if err := h.accountService.UpdateAccountStatus(id, req.Status); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.WriteSuccess(w, http.StatusOK, "Account status updated successfully", nil)
}
