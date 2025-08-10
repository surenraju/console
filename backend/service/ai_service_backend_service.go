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

// AIServiceBackendService handles operations for AIServiceBackend resources
type AIServiceBackendService struct {
	client client.Client
	logger logr.Logger
}

// NewAIServiceBackendService creates a new AIServiceBackendService
func NewAIServiceBackendService(client client.Client, logger logr.Logger) *AIServiceBackendService {
	return &AIServiceBackendService{
		client: client,
		logger: logger,
	}
}

// Create creates a new AIServiceBackend
func (s *AIServiceBackendService) Create(ctx context.Context, backend *aigv1a1.AIServiceBackend) error {
	if err := s.client.Create(ctx, backend); err != nil {
		return fmt.Errorf("failed to create AIServiceBackend: %w", err)
	}
	return nil
}

// Get retrieves a specific AIServiceBackend by name in a namespace
func (s *AIServiceBackendService) Get(ctx context.Context, namespace, name string) (*aigv1a1.AIServiceBackend, error) {
	var backend aigv1a1.AIServiceBackend
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &backend); err != nil {
		return nil, fmt.Errorf("failed to get AIServiceBackend: %w", err)
	}
	return &backend, nil
}

// List retrieves all AIServiceBackend resources in a namespace
func (s *AIServiceBackendService) List(ctx context.Context, namespace string) (*aigv1a1.AIServiceBackendList, error) {
	var list aigv1a1.AIServiceBackendList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list AIServiceBackends: %w", err)
	}
	return &list, nil
}

// Update updates an existing AIServiceBackend
func (s *AIServiceBackendService) Update(ctx context.Context, backend *aigv1a1.AIServiceBackend) error {
	if err := s.client.Update(ctx, backend); err != nil {
		return fmt.Errorf("failed to update AIServiceBackend: %w", err)
	}
	return nil
}

// Delete deletes an AIServiceBackend by name in a namespace
func (s *AIServiceBackendService) Delete(ctx context.Context, namespace, name string) error {
	backend := &aigv1a1.AIServiceBackend{}
	backend.Namespace = namespace
	backend.Name = name
	if err := s.client.Delete(ctx, backend); err != nil {
		return fmt.Errorf("failed to delete AIServiceBackend: %w", err)
	}
	return nil
}