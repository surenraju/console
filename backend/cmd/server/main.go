// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envoyproxy/ai-gateway/console/backend/internal/router"
	"github.com/envoyproxy/ai-gateway/console/backend/internal/server"
)

func main() {
	log.Println("Starting Envoy AI Gateway Console Backend...")

	// Create server
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create router
	rt := router.NewRouter(srv)

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:         ":8081",
		Handler:      rt,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on :8081")
		log.Println("Available endpoints:")
		log.Println("  GET /api/v1/llm/providers      - List all LLM providers")
		log.Println("  POST /api/v1/llm/providers     - Create a new LLM provider")
		log.Println("  GET /api/v1/llm/providers/{name} - Get specific LLM provider")
		log.Println("  DELETE /api/v1/llm/providers/{name} - Delete an LLM provider")
		log.Println("  GET /health                     - Health check")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Gracefully shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
