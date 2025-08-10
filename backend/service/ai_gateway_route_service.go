// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package services

import (
	"context"
	"fmt"

	aigv1a1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AIGatewayRouteService handles operations for AIGatewayRoute resources
type AIGatewayRouteService struct {
	client client.Client
	logger logr.Logger
}

// NewAIGatewayRouteService creates a new AIGatewayRouteService
func NewAIGatewayRouteService(client client.Client, logger logr.Logger) *AIGatewayRouteService {
	return &AIGatewayRouteService{
		client: client,
		logger: logger,
	}
}

// Create creates a new AIGatewayRoute
func (s *AIGatewayRouteService) Create(ctx context.Context, route *aigv1a1.AIGatewayRoute) error {
	if err := s.client.Create(ctx, route); err != nil {
		return fmt.Errorf("failed to create AIGatewayRoute: %w", err)
	}
	return nil
}

// Get retrieves a specific AIGatewayRoute by name in a namespace
func (s *AIGatewayRouteService) Get(ctx context.Context, namespace, name string) (*aigv1a1.AIGatewayRoute, error) {
	var route aigv1a1.AIGatewayRoute
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &route); err != nil {
		return nil, fmt.Errorf("failed to get AIGatewayRoute: %w", err)
	}
	return &route, nil
}

// List retrieves all AIGatewayRoute resources in a namespace
func (s *AIGatewayRouteService) List(ctx context.Context, namespace string) (*aigv1a1.AIGatewayRouteList, error) {
	var list aigv1a1.AIGatewayRouteList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list AIGatewayRoutes: %w", err)
	}
	return &list, nil
}

// Update updates an existing AIGatewayRoute
func (s *AIGatewayRouteService) Update(ctx context.Context, route *aigv1a1.AIGatewayRoute) error {
	if err := s.client.Update(ctx, route); err != nil {
		return fmt.Errorf("failed to update AIGatewayRoute: %w", err)
	}
	return nil
}

// Delete deletes an AIGatewayRoute by name in a namespace
func (s *AIGatewayRouteService) Delete(ctx context.Context, namespace, name string) error {
	route := &aigv1a1.AIGatewayRoute{}
	route.Namespace = namespace
	route.Name = name
	if err := s.client.Delete(ctx, route); err != nil {
		return fmt.Errorf("failed to delete AIGatewayRoute: %w", err)
	}
	return nil
}
