// Package httpserver provides HTTP response helpers.
// All handlers should use these helpers to ensure consistent JSON responses.
package httpserver

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse is the standard envelope for successful API responses.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// ErrorResponse is the standard envelope for error API responses.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// JSON writes a JSON response with the given status code and payload.
func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success sends a successful JSON response.
func Success(w http.ResponseWriter, statusCode int, message string, data any) {
	JSON(w, statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error JSON response.
func Error(w http.ResponseWriter, statusCode int, message string, err string) {
	JSON(w, statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
