// Package customer – service.go contains all business logic for customers.
// The service sits between the handler and the repository.
package customer

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service contains business logic for customer operations.
type Service struct {
	repo   *Repository
	logger *slog.Logger
}

// NewService creates a new customer service with its dependencies.
func NewService(repo *Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create validates input, builds a customer, and persists it.
func (s *Service) Create(ctx context.Context, req *CreateRequest) (*Customer, error) {
	now := time.Now().UTC()

	customer := &Customer{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		s.logger.Error("failed to create customer", "error", err, "email", req.Email)
		return nil, err
	}

	s.logger.Info("customer created", "id", customer.ID, "email", customer.Email)
	return customer, nil
}

// GetByID retrieves a single customer by their ID.
func (s *Service) GetByID(ctx context.Context, id string) (*Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer", "error", err, "id", id)
		return nil, err
	}
	return customer, nil
}

// List retrieves all customers.
func (s *Service) List(ctx context.Context) ([]Customer, error) {
	customers, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("failed to list customers", "error", err)
		return nil, err
	}
	return customers, nil
}

// Update modifies an existing customer with the provided data.
func (s *Service) Update(ctx context.Context, id string, req *UpdateRequest) (*Customer, error) {
	// Fetch the existing customer.
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to find customer for update", "error", err, "id", id)
		return nil, err
	}

	// Apply changes.
	existing.Name = req.Name
	existing.Email = req.Email
	existing.Phone = req.Phone
	existing.UpdatedAt = time.Now().UTC()

	// Only update status if provided; otherwise keep the existing status.
	if req.Status != "" {
		existing.Status = req.Status
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("failed to update customer", "error", err, "id", id)
		return nil, err
	}

	s.logger.Info("customer updated", "id", id)
	return existing, nil
}

// Delete removes a customer by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete customer", "error", err, "id", id)
		return err
	}

	s.logger.Info("customer deleted", "id", id)
	return nil
}
