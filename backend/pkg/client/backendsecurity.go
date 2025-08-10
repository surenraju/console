// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"fmt"

	aigv1a1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackendSecurityPolicyClient handles operations for BackendSecurityPolicy resources
type BackendSecurityPolicyClient struct {
	client client.Client
	logger logr.Logger
}

// NewBackendSecurityPolicyClient creates a new BackendSecurityPolicyClient
func NewBackendSecurityPolicyClient(client client.Client, logger logr.Logger) *BackendSecurityPolicyClient {
	return &BackendSecurityPolicyClient{
		client: client,
		logger: logger,
	}
}

// Create creates a new BackendSecurityPolicy
func (c *BackendSecurityPolicyClient) Create(ctx context.Context, policy *aigv1a1.BackendSecurityPolicy) error {
	if err := c.client.Create(ctx, policy); err != nil {
		return fmt.Errorf("failed to create BackendSecurityPolicy: %w", err)
	}
	return nil
}

// Get retrieves a specific BackendSecurityPolicy by name in a namespace
func (c *BackendSecurityPolicyClient) Get(ctx context.Context, namespace, name string) (*aigv1a1.BackendSecurityPolicy, error) {
	var policy aigv1a1.BackendSecurityPolicy
	if err := c.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &policy); err != nil {
		return nil, fmt.Errorf("failed to get BackendSecurityPolicy: %w", err)
	}
	return &policy, nil
}

// List retrieves all BackendSecurityPolicy resources in a namespace
func (c *BackendSecurityPolicyClient) List(ctx context.Context, namespace string) (*aigv1a1.BackendSecurityPolicyList, error) {
	var list aigv1a1.BackendSecurityPolicyList
	if err := c.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list BackendSecurityPolicies: %w", err)
	}
	return &list, nil
}

// Update updates an existing BackendSecurityPolicy
func (c *BackendSecurityPolicyClient) Update(ctx context.Context, policy *aigv1a1.BackendSecurityPolicy) error {
	if err := c.client.Update(ctx, policy); err != nil {
		return fmt.Errorf("failed to update BackendSecurityPolicy: %w", err)
	}
	return nil
}

// Delete deletes a BackendSecurityPolicy by name in a namespace
func (c *BackendSecurityPolicyClient) Delete(ctx context.Context, namespace, name string) error {
	policy := &aigv1a1.BackendSecurityPolicy{}
	policy.Namespace = namespace
	policy.Name = name
	if err := c.client.Delete(ctx, policy); err != nil {
		return fmt.Errorf("failed to delete BackendSecurityPolicy: %w", err)
	}
	return nil
}
