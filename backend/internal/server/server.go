// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/envoyproxy/ai-gateway/console/backend/internal/service"
	"github.com/envoyproxy/ai-gateway/console/backend/pkg/client"
	"github.com/envoyproxy/ai-gateway/console/backend/pkg/llm"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

// Server holds the HTTP server and service dependencies
type Server struct {
	clientManager      *client.Manager
	llmProviderService *service.LLMProviderService
}

// NewServer creates a new server instance
func NewServer() (*Server, error) {
	// Create client manager
	config := client.Config{
		Logger: logr.Discard(), // Use a simple logger for now
	}

	clientManager, err := client.NewManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client manager: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := clientManager.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("kubernetes connection failed: %w", err)
	}

	server := &Server{
		clientManager:      clientManager,
		llmProviderService: service.NewLLMProviderService(clientManager),
	}

	return server, nil
}

// GetLLMProviders handles GET /api/v1/llm/providers
func (s *Server) GetLLMProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Getting LLM providers from namespace: %s", namespace)

	providers, err := s.llmProviderService.ListProviders(ctx, namespace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list LLM providers: %v", err), http.StatusInternalServerError)
		return
	}

	// Ensure we always return an array, never null
	if providers == nil {
		providers = []llm.LLMProvider{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// GetLLMProvider handles GET /api/v1/llm/providers/{name}
func (s *Server) GetLLMProvider(w http.ResponseWriter, r *http.Request) {
	// Extract name from URL path manually or use query parameter
	name := r.URL.Query().Get("name")
	if name == "" {
		// Extract from path - simple extraction for now
		// This would be better handled by a proper router
		http.Error(w, "Provider name is required", http.StatusBadRequest)
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	ctx := r.Context()

	provider, err := s.llmProviderService.GetProvider(ctx, namespace, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("LLM provider not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

// HealthCheck handles GET /health
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.clientManager.HealthCheck(ctx); err != nil {
		http.Error(w, fmt.Sprintf("Health check failed: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// GetLLMProviderByName handles GET /api/v1/llm/providers/:name with Gin
func (s *Server) GetLLMProviderByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider name is required"})
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	provider, err := s.llmProviderService.GetProvider(c.Request.Context(), namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("LLM provider not found: %v", err)})
		return
	}

	c.JSON(http.StatusOK, provider)
}

// CreateLLMProvider handles POST /api/v1/llm/providers with Gin
func (s *Server) CreateLLMProvider(c *gin.Context) {
	var provider llm.LLMProvider

	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	// Validate required fields
	if provider.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider name is required"})
		return
	}

	if provider.Namespace == "" {
		provider.Namespace = "default"
	}

	// Create the provider
	err := s.llmProviderService.CreateProvider(c.Request.Context(), &provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create LLM provider: %v", err)})
		return
	}

	// Return the created provider (masked)
	maskedProvider := provider.MaskSecret()
	c.JSON(http.StatusCreated, maskedProvider)
}

// DeleteLLMProvider handles DELETE /api/v1/llm/providers/{name} with Gin
func (s *Server) DeleteLLMProvider(c *gin.Context) {
	name := c.Param("name")
	namespace := c.DefaultQuery("namespace", "default")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider name is required"})
		return
	}

	// Delete the provider
	err := s.llmProviderService.DeleteProvider(c.Request.Context(), namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete LLM provider: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Provider '%s' deleted successfully", name)})
}
