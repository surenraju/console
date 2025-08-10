// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

// BackendTLSPolicyClient handles operations for BackendTLSPolicy resources
type BackendTLSPolicyClient struct {
	client client.Client
	logger logr.Logger
}

// NewBackendTLSPolicyClient creates a new BackendTLSPolicyClient
func NewBackendTLSPolicyClient(client client.Client, logger logr.Logger) *BackendTLSPolicyClient {
	return &BackendTLSPolicyClient{
		client: client,
		logger: logger,
	}
}

// Create creates a new BackendTLSPolicy
func (c *BackendTLSPolicyClient) Create(ctx context.Context, policy *gwapiv1a3.BackendTLSPolicy) error {
	if err := c.client.Create(ctx, policy); err != nil {
		return fmt.Errorf("failed to create BackendTLSPolicy: %w", err)
	}
	return nil
}

// Get retrieves a specific BackendTLSPolicy by name in a namespace
func (c *BackendTLSPolicyClient) Get(ctx context.Context, namespace, name string) (*gwapiv1a3.BackendTLSPolicy, error) {
	var policy gwapiv1a3.BackendTLSPolicy
	if err := c.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &policy); err != nil {
		return nil, fmt.Errorf("failed to get BackendTLSPolicy: %w", err)
	}
	return &policy, nil
}

// List retrieves all BackendTLSPolicy resources in a namespace
func (c *BackendTLSPolicyClient) List(ctx context.Context, namespace string) (*gwapiv1a3.BackendTLSPolicyList, error) {
	var list gwapiv1a3.BackendTLSPolicyList
	if err := c.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list BackendTLSPolicies: %w", err)
	}
	return &list, nil
}

// Update updates an existing BackendTLSPolicy
func (c *BackendTLSPolicyClient) Update(ctx context.Context, policy *gwapiv1a3.BackendTLSPolicy) error {
	if err := c.client.Update(ctx, policy); err != nil {
		return fmt.Errorf("failed to update BackendTLSPolicy: %w", err)
	}
	return nil
}

// Delete deletes a BackendTLSPolicy by name in a namespace
func (c *BackendTLSPolicyClient) Delete(ctx context.Context, namespace, name string) error {
	policy := &gwapiv1a3.BackendTLSPolicy{}
	policy.Namespace = namespace
	policy.Name = name
	if err := c.client.Delete(ctx, policy); err != nil {
		return fmt.Errorf("failed to delete BackendTLSPolicy: %w", err)
	}
	return nil
}
