// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretClient handles operations for Secret resources
type SecretClient struct {
	client client.Client
	logger logr.Logger
}

// NewSecretClient creates a new SecretClient
func NewSecretClient(client client.Client, logger logr.Logger) *SecretClient {
	return &SecretClient{
		client: client,
		logger: logger,
	}
}

// Create creates a new Secret
func (c *SecretClient) Create(ctx context.Context, secret *corev1.Secret) error {
	if err := c.client.Create(ctx, secret); err != nil {
		return fmt.Errorf("failed to create Secret: %w", err)
	}
	return nil
}

// Get retrieves a specific Secret by name in a namespace
func (c *SecretClient) Get(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	var secret corev1.Secret
	if err := c.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &secret); err != nil {
		return nil, fmt.Errorf("failed to get Secret: %w", err)
	}
	return &secret, nil
}

// List retrieves all Secret resources in a namespace
func (c *SecretClient) List(ctx context.Context, namespace string) (*corev1.SecretList, error) {
	var list corev1.SecretList
	if err := c.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list Secrets: %w", err)
	}
	return &list, nil
}

// Update updates an existing Secret
func (c *SecretClient) Update(ctx context.Context, secret *corev1.Secret) error {
	if err := c.client.Update(ctx, secret); err != nil {
		return fmt.Errorf("failed to update Secret: %w", err)
	}
	return nil
}

// Delete deletes a Secret by name in a namespace
func (c *SecretClient) Delete(ctx context.Context, namespace, name string) error {
	secret := &corev1.Secret{}
	secret.Namespace = namespace
	secret.Name = name
	if err := c.client.Delete(ctx, secret); err != nil {
		return fmt.Errorf("failed to delete Secret: %w", err)
	}
	return nil
}
