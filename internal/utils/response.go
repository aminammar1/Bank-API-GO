package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bank-api/internal/models"
)

// WriteJSON writes a JSON response to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response to the http.ResponseWriter
func WriteError(w http.ResponseWriter, status int, message string) error {
	errorResponse := models.ErrorResponse{
		Error:     message,
		Timestamp: time.Now().UTC(),
	}
	return WriteJSON(w, status, errorResponse)
}

// WriteSuccess writes a success response to the http.ResponseWriter
func WriteSuccess(w http.ResponseWriter, status int, message string, data interface{}) error {
	successResponse := models.SuccessResponse{
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	return WriteJSON(w, status, successResponse)
}

// ParseJSON parses JSON request body into the provided interface
func ParseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
