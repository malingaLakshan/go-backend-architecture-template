// Package customer – dto.go defines request and response data transfer objects.
// DTOs decouple the API contract from the internal domain model.
package customer

// CreateRequest is the expected JSON body for creating a customer.
type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UpdateRequest is the expected JSON body for updating a customer.
type UpdateRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

// Response is the JSON representation returned to API callers.
type Response struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToResponse converts a domain Customer to an API Response.
func ToResponse(c *Customer) *Response {
	return &Response{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ToResponseList converts a slice of domain Customers to a slice of API Responses.
func ToResponseList(customers []Customer) []Response {
	responses := make([]Response, len(customers))
	for i, c := range customers {
		responses[i] = *ToResponse(&c)
	}
	return responses
}
