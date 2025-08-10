// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package router

import (
	"github.com/envoyproxy/ai-gateway/console/backend/internal/server"
	"github.com/gin-gonic/gin"
)

// NewRouter creates a new Gin router with all routes configured
func NewRouter(srv *server.Server) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())   // Recovery middleware
	router.Use(corsMiddleware()) // CORS middleware

	// Health check endpoint
	router.GET("/health", gin.WrapF(srv.HealthCheck))

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// LLM provider routes
		llm := apiV1.Group("/llm")
		{
			llm.GET("/providers", gin.WrapF(srv.GetLLMProviders))
			llm.GET("/providers/:name", srv.GetLLMProviderByName)
		}
	}

	return router
}

// corsMiddleware returns a Gin middleware for CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
