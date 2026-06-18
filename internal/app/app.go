// Package app wires all dependencies together and starts the HTTP server.
// This is the composition root of the application.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-architecture-template/internal/config"
	"go-backend-architecture-template/internal/customer"
	"go-backend-architecture-template/internal/database"
	"go-backend-architecture-template/internal/logger"
	"go-backend-architecture-template/internal/transport/httpserver"
)

// Run initialises all components and starts the server.
// This function blocks until the server is shut down gracefully.
func Run() error {
	// ── 1. Load configuration ───────────────────────────────────────────
	cfg := config.Load()

	// ── 2. Setup logger ─────────────────────────────────────────────────
	log, err := logger.Setup(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}
	log.Info("logger initialised", "level", cfg.Logger.Level)

	// ── 3. Connect to database ──────────────────────────────────────────
	db, err := database.Connect(cfg.Database.FilePath, log)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// ── 4. Run database migrations ──────────────────────────────────────
	if err := database.RunMigrations(db, log); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// ── 5. Wire dependencies ────────────────────────────────────────────
	// Customer module
	customerRepo := customer.NewRepository(db)
	customerService := customer.NewService(customerRepo, log)
	customerHandler := customer.NewHandler(customerService)

	// ── 6. Create router ────────────────────────────────────────────────
	router := httpserver.NewRouter(log, customerHandler)

	// ── 7. Create HTTP server ───────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// ── 8. Graceful shutdown ────────────────────────────────────────────
	// Start server in a goroutine.
	errCh := make(chan error, 1)
	go func() {
		log.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("shutdown signal received", "signal", sig.String())
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	// Give outstanding requests 10 seconds to complete.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info("server stopped gracefully")
	return nil
}
