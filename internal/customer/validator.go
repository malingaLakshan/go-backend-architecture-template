// Package customer – validator.go validates incoming request data.
// Validation is called from the handler before passing data to the service.
package customer

import (
	"fmt"
	"net/mail"
	"strings"
)

// ValidateCreateRequest checks the CreateRequest fields for correctness.
func ValidateCreateRequest(req *CreateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("email format is invalid")
	}
	return nil
}

// ValidateUpdateRequest checks the UpdateRequest fields for correctness.
func ValidateUpdateRequest(req *UpdateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("email format is invalid")
	}
	if req.Status != "" && req.Status != StatusActive && req.Status != StatusInactive {
		return fmt.Errorf("status must be ACTIVE or INACTIVE")
	}
	return nil
}
