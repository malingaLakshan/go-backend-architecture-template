// Package customer – errors.go defines customer-specific error types.
// These errors are returned by the service layer and interpreted by the handler.
package customer

import "errors"

var (
	// ErrNotFound is returned when a customer cannot be found by ID.
	ErrNotFound = errors.New("customer not found")

	// ErrDuplicateEmail is returned when a customer with the same email already exists.
	ErrDuplicateEmail = errors.New("a customer with this email already exists")

	// ErrInvalidStatus is returned when an invalid status value is provided.
	ErrInvalidStatus = errors.New("invalid status: must be ACTIVE or INACTIVE")
)
