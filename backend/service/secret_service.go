// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package services

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretService handles operations for Secret resources
type SecretService struct {
	client client.Client
	logger logr.Logger
}

// NewSecretService creates a new SecretService
func NewSecretService(client client.Client, logger logr.Logger) *SecretService {
	return &SecretService{
		client: client,
		logger: logger,
	}
}

// Create creates a new Secret
func (s *SecretService) Create(ctx context.Context, secret *corev1.Secret) error {
	if err := s.client.Create(ctx, secret); err != nil {
		return fmt.Errorf("failed to create Secret: %w", err)
	}
	return nil
}

// Get retrieves a specific Secret by name in a namespace
func (s *SecretService) Get(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	var secret corev1.Secret
	if err := s.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &secret); err != nil {
		return nil, fmt.Errorf("failed to get Secret: %w", err)
	}
	return &secret, nil
}

// List retrieves all Secret resources in a namespace
func (s *SecretService) List(ctx context.Context, namespace string) (*corev1.SecretList, error) {
	var list corev1.SecretList
	if err := s.client.List(ctx, &list, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list Secrets: %w", err)
	}
	return &list, nil
}

// Update updates an existing Secret
func (s *SecretService) Update(ctx context.Context, secret *corev1.Secret) error {
	if err := s.client.Update(ctx, secret); err != nil {
		return fmt.Errorf("failed to update Secret: %w", err)
	}
	return nil
}

// Delete deletes a Secret by name in a namespace
func (s *SecretService) Delete(ctx context.Context, namespace, name string) error {
	secret := &corev1.Secret{}
	secret.Namespace = namespace
	secret.Name = name
	if err := s.client.Delete(ctx, secret); err != nil {
		return fmt.Errorf("failed to delete Secret: %w", err)
	}
	return nil
}