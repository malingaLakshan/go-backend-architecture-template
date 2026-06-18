// Package customer – repository.go handles all SQLite queries for customers.
// The repository is the only layer that interacts with the database.
package customer

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

// Repository provides access to customer storage.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new customer repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// timeFormat is the layout used to store timestamps as TEXT in SQLite.
const timeFormat = "2006-01-02T15:04:05Z"

// Create inserts a new customer into the database.
func (r *Repository) Create(ctx context.Context, c *Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID,
		c.Name,
		c.Email,
		c.Phone,
		c.Status,
		c.CreatedAt.Format(timeFormat),
		c.UpdatedAt.Format(timeFormat),
	)
	if err != nil {
		// Check for unique constraint violation on email.
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

// GetByID retrieves a single customer by their ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*Customer, error) {
	query := `SELECT id, name, email, phone, status, created_at, updated_at FROM customers WHERE id = ?`

	var c Customer
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	c.CreatedAt, _ = time.Parse(timeFormat, createdAt)
	c.UpdatedAt, _ = time.Parse(timeFormat, updatedAt)
	return &c, nil
}

// List retrieves all customers from the database ordered by creation date.
func (r *Repository) List(ctx context.Context) ([]Customer, error) {
	query := `SELECT id, name, email, phone, status, created_at, updated_at FROM customers ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		var createdAt, updatedAt string

		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		c.CreatedAt, _ = time.Parse(timeFormat, createdAt)
		c.UpdatedAt, _ = time.Parse(timeFormat, updatedAt)
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

// Update modifies an existing customer in the database.
func (r *Repository) Update(ctx context.Context, c *Customer) error {
	query := `
		UPDATE customers
		SET name = ?, email = ?, phone = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		c.Name, c.Email, c.Phone, c.Status,
		c.UpdatedAt.Format(timeFormat),
		c.ID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrDuplicateEmail
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a customer from the database by ID.
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM customers WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
