// Package main is the entry point for the application.
// It should only start the application — all wiring happens in internal/app.
package main

import (
	"log"

	"go-backend-architecture-template/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
