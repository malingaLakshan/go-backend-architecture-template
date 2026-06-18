// Package customer – model.go defines the domain model for the customer entity.
// This struct maps directly to the database table.
package customer

import "time"

// Status constants define the allowed customer statuses.
const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)

// Customer is the domain model representing a customer in the system.
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
