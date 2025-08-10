// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package services

import (
	"context"
	"fmt"

	gwapiv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackendService handles operations for Backend resources
type BackendService struct {
	client client.Client
	logger logr.Logger
}

// NewBackendService creates a new BackendService
func NewBackendService(client client.Client, logger logr.Logger) *BackendService {
	return &BackendService{
		client: client,
		logger: logger,
	}
}

// Create creates a new Backend
func (s *BackendService) Create(ctx context.Context, backend *gwapiv1a1.Backend) error {
	if err := s.client.Create(ctx, backend); err != nil {
		return fmt.Errorf("failed to create Backend: %w", err)
	}
	return nil
}

// Get retrieves a specific Backend by name in a namespace
func (s *BackendService) Get(ctx context.Context, namespace, name string) (*gwapiv1a1.Backend, error) {
	var backend gwapiv1a1.Backend
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &backend); err != nil {
		return nil, fmt.Errorf("failed to get Backend: %w", err)
	}
	return &backend, nil
}

// List retrieves all Backend resources in a namespace
func (s *BackendService) List(ctx context.Context, namespace string) (*gwapiv1a1.BackendList, error) {
	var list gwapiv1a1.BackendList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list Backends: %w", err)
	}
	return &list, nil
}

// Update updates an existing Backend
func (s *BackendService) Update(ctx context.Context, backend *gwapiv1a1.Backend) error {
	if err := s.client.Update(ctx, backend); err != nil {
		return fmt.Errorf("failed to update Backend: %w", err)
	}
	return nil
}

// Delete deletes a Backend by name in a namespace
func (s *BackendService) Delete(ctx context.Context, namespace, name string) error {
	backend := &gwapiv1a1.Backend{}
	backend.Namespace = namespace
	backend.Name = name
	if err := s.client.Delete(ctx, backend); err != nil {
		return fmt.Errorf("failed to delete Backend: %w", err)
	}
	return nil
}
