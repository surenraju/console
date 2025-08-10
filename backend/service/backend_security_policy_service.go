// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package services

import (
	"context"
	"fmt"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackendSecurityPolicyService handles operations for BackendSecurityPolicy resources
type BackendSecurityPolicyService struct {
	client client.Client
	logger logr.Logger
}

// NewBackendSecurityPolicyService creates a new BackendSecurityPolicyService
func NewBackendSecurityPolicyService(client client.Client, logger logr.Logger) *BackendSecurityPolicyService {
	return &BackendSecurityPolicyService{
		client: client,
		logger: logger,
	}
}

// Create creates a new SecurityPolicy
func (s *BackendSecurityPolicyService) Create(ctx context.Context, policy *egv1a1.SecurityPolicy) error {
	if err := s.client.Create(ctx, policy); err != nil {
		return fmt.Errorf("failed to create SecurityPolicy: %w", err)
	}
	return nil
}

// Get retrieves a specific SecurityPolicy by name in a namespace
func (s *BackendSecurityPolicyService) Get(ctx context.Context, namespace, name string) (*egv1a1.SecurityPolicy, error) {
	var policy egv1a1.SecurityPolicy
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &policy); err != nil {
		return nil, fmt.Errorf("failed to get SecurityPolicy: %w", err)
	}
	return &policy, nil
}

// List retrieves all SecurityPolicy resources in a namespace
func (s *BackendSecurityPolicyService) List(ctx context.Context, namespace string) (*egv1a1.SecurityPolicyList, error) {
	var list egv1a1.SecurityPolicyList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list SecurityPolicies: %w", err)
	}
	return &list, nil
}

// Update updates an existing SecurityPolicy
func (s *BackendSecurityPolicyService) Update(ctx context.Context, policy *egv1a1.SecurityPolicy) error {
	if err := s.client.Update(ctx, policy); err != nil {
		return fmt.Errorf("failed to update SecurityPolicy: %w", err)
	}
	return nil
}

// Delete deletes a SecurityPolicy by name in a namespace
func (s *BackendSecurityPolicyService) Delete(ctx context.Context, namespace, name string) error {
	policy := &egv1a1.SecurityPolicy{}
	policy.Namespace = namespace
	policy.Name = name
	if err := s.client.Delete(ctx, policy); err != nil {
		return fmt.Errorf("failed to delete SecurityPolicy: %w", err)
	}
	return nil
}