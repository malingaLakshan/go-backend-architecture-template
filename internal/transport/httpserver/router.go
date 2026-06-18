// Package httpserver – router.go sets up the HTTP router and registers all routes.
package httpserver

import (
	"log/slog"
	"net/http"

	"go-backend-architecture-template/internal/customer"
)

// NewRouter creates the application router with all middleware and routes.
func NewRouter(logger *slog.Logger, customerHandler *customer.Handler) http.Handler {
	mux := http.NewServeMux()

	// ── Health check ────────────────────────────────────────────────────
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		Success(w, http.StatusOK, "Server is running", nil)
	})

	// ── Customer routes ─────────────────────────────────────────────────
	mux.HandleFunc("POST /api/customers", customerHandler.Create)
	mux.HandleFunc("GET /api/customers", customerHandler.List)
	mux.HandleFunc("GET /api/customers/{id}", customerHandler.GetByID)
	mux.HandleFunc("PUT /api/customers/{id}", customerHandler.Update)
	mux.HandleFunc("DELETE /api/customers/{id}", customerHandler.Delete)

	// ── Apply middleware (outermost runs first) ─────────────────────────
	var handler http.Handler = mux
	handler = CORS()(handler)
	handler = RequestLogger(logger)(handler)
	handler = PanicRecovery(logger)(handler)

	return handler
}
