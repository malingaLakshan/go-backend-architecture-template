// Package customer – handler.go handles HTTP requests and responses.
// The handler parses input, calls the service, and writes JSON responses.
// It contains NO business logic.
package customer

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-backend-architecture-template/internal/transport/httpserver"
)

// Handler handles HTTP requests for customer operations.
type Handler struct {
	service *Service
}

// NewHandler creates a new customer handler with its service dependency.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /api/customers.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.Error(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := ValidateCreateRequest(&req); err != nil {
		httpserver.Error(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	customer, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	httpserver.Success(w, http.StatusCreated, "Customer created successfully", ToResponse(customer))
}

// List handles GET /api/customers.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.List(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	httpserver.Success(w, http.StatusOK, "Customers retrieved successfully", ToResponseList(customers))
}

// GetByID handles GET /api/customers/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpserver.Error(w, http.StatusBadRequest, "Invalid request", "customer id is required")
		return
	}

	customer, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	httpserver.Success(w, http.StatusOK, "Customer retrieved successfully", ToResponse(customer))
}

// Update handles PUT /api/customers/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpserver.Error(w, http.StatusBadRequest, "Invalid request", "customer id is required")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpserver.Error(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := ValidateUpdateRequest(&req); err != nil {
		httpserver.Error(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	customer, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	httpserver.Success(w, http.StatusOK, "Customer updated successfully", ToResponse(customer))
}

// Delete handles DELETE /api/customers/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpserver.Error(w, http.StatusBadRequest, "Invalid request", "customer id is required")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	httpserver.Success(w, http.StatusOK, "Customer deleted successfully", nil)
}

// handleServiceError maps domain errors to HTTP responses.
func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpserver.Error(w, http.StatusNotFound, "Not found", err.Error())
	case errors.Is(err, ErrDuplicateEmail):
		httpserver.Error(w, http.StatusConflict, "Conflict", err.Error())
	case errors.Is(err, ErrInvalidStatus):
		httpserver.Error(w, http.StatusBadRequest, "Validation failed", err.Error())
	default:
		httpserver.Error(w, http.StatusInternalServerError, "Internal server error", "an unexpected error occurred")
	}
}
