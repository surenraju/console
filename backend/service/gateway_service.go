// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package services

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GatewayService handles operations for Gateway resources
type GatewayService struct {
	client client.Client
	logger logr.Logger
}

// NewGatewayService creates a new GatewayService
func NewGatewayService(client client.Client, logger logr.Logger) *GatewayService {
	return &GatewayService{
		client: client,
		logger: logger,
	}
}

// Create creates a new Gateway
func (s *GatewayService) Create(ctx context.Context, gateway *gwapiv1.Gateway) error {
	if err := s.client.Create(ctx, gateway); err != nil {
		return fmt.Errorf("failed to create Gateway: %w", err)
	}
	return nil
}

// Get retrieves a specific Gateway by name in a namespace
func (s *GatewayService) Get(ctx context.Context, namespace, name string) (*gwapiv1.Gateway, error) {
	var gateway gwapiv1.Gateway
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &gateway); err != nil {
		return nil, fmt.Errorf("failed to get Gateway: %w", err)
	}
	return &gateway, nil
}

// List retrieves all Gateway resources in a namespace
func (s *GatewayService) List(ctx context.Context, namespace string) (*gwapiv1.GatewayList, error) {
	var list gwapiv1.GatewayList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list Gateways: %w", err)
	}
	return &list, nil
}

// Update updates an existing Gateway
func (s *GatewayService) Update(ctx context.Context, gateway *gwapiv1.Gateway) error {
	if err := s.client.Update(ctx, gateway); err != nil {
		return fmt.Errorf("failed to update Gateway: %w", err)
	}
	return nil
}

// Delete deletes a Gateway by name in a namespace
func (s *GatewayService) Delete(ctx context.Context, namespace, name string) error {
	gateway := &gwapiv1.Gateway{}
	gateway.Namespace = namespace
	gateway.Name = name
	if err := s.client.Delete(ctx, gateway); err != nil {
		return fmt.Errorf("failed to delete Gateway: %w", err)
	}
	return nil
}