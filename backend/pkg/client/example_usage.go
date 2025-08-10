// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"log"

	"github.com/go-logr/logr"
)

// ExampleUsage demonstrates how to use the client manager
func ExampleUsage() {
	// Example usage of the client manager

	// Create a new client manager
	config := Config{
		// Kubeconfig: "/path/to/kubeconfig", // Optional, uses in-cluster config if empty
		Logger: logr.Discard(), // Or use a real logger
	}

	manager, err := NewManager(config)
	if err != nil {
		log.Fatalf("Failed to create client manager: %v", err)
	}

	ctx := context.Background()

	// Perform a health check
	if err := manager.HealthCheck(ctx); err != nil {
		log.Fatalf("Kubernetes health check failed: %v", err)
	}

	// Example: List all backends in a namespace
	backends, err := manager.Backend.List(ctx, "default")
	if err != nil {
		log.Printf("Failed to list backends: %v", err)
	} else {
		log.Printf("Found %d backends", len(backends.Items))
	}

	// Example: Load all Envoy Gateway resources for a specific provider
	resources, err := manager.LoadEnvoyGatewayResources(ctx, "default", "my-llm-provider")
	if err != nil {
		log.Printf("Failed to load Envoy Gateway resources: %v", err)
	} else {
		log.Printf("Found %d resources for provider", len(resources))
	}

	// Example: Get a specific secret
	secret, err := manager.Secret.Get(ctx, "default", "my-api-key")
	if err != nil {
		log.Printf("Failed to get secret: %v", err)
	} else {
		log.Printf("Found secret with %d keys", len(secret.Data))
	}

	log.Println("Client manager example completed successfully")
}
