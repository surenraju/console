// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"fmt"

	aigv1a1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	gwapiv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

// Manager provides centralized access to all Kubernetes resource clients
type Manager struct {
	client client.Client
	logger logr.Logger

	// Individual typed clients for each resource type
	Backend               *BackendClient
	BackendTLSPolicy      *BackendTLSPolicyClient
	Secret                *SecretClient
	BackendSecurityPolicy *BackendSecurityPolicyClient
	AIServiceBackend      *AIServiceBackendClient
}

// Config holds configuration for the Kubernetes client manager
type Config struct {
	// Kubeconfig path, if empty uses in-cluster config
	Kubeconfig string
	// Logger for the client manager
	Logger logr.Logger
}

// NewManager creates a new Kubernetes client manager with all resource services
func NewManager(cfg Config) (*Manager, error) {
	// Get Kubernetes config
	var restConfig *rest.Config
	var err error

	if cfg.Kubeconfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	} else {
		restConfig, err = config.GetConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config: %w", err)
	}

	// Create a new scheme and add all required types
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add core/v1 to scheme: %w", err)
	}
	if err := gwapiv1.Install(scheme); err != nil {
		return nil, fmt.Errorf("failed to add gateway-api/v1 to scheme: %w", err)
	}
	if err := gwapiv1a3.Install(scheme); err != nil {
		return nil, fmt.Errorf("failed to add gateway-api/v1alpha3 to scheme: %w", err)
	}
	if err := gwapiv1a1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add envoy-gateway/v1alpha1 to scheme: %w", err)
	}
	if err := aigv1a1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add ai-gateway/v1alpha1 to scheme: %w", err)
	}

	// Create the controller-runtime client
	k8sClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	logger := cfg.Logger
	if logger.GetSink() == nil {
		// Use a no-op logger if none provided
		logger = logr.Discard()
	}

	// Initialize all clients
	manager := &Manager{
		client:                k8sClient,
		logger:                logger,
		Backend:               NewBackendClient(k8sClient, logger),
		BackendTLSPolicy:      NewBackendTLSPolicyClient(k8sClient, logger),
		Secret:                NewSecretClient(k8sClient, logger),
		BackendSecurityPolicy: NewBackendSecurityPolicyClient(k8sClient, logger),
		AIServiceBackend:      NewAIServiceBackendClient(k8sClient, logger),
	}

	return manager, nil
}

// Client returns the underlying Kubernetes client
func (m *Manager) Client() client.Client {
	return m.client
}

// HealthCheck performs a basic health check against the Kubernetes API
func (m *Manager) HealthCheck(ctx context.Context) error {
	// Try to list namespaces as a basic connectivity test
	var namespaces corev1.NamespaceList
	if err := m.client.List(ctx, &namespaces, client.Limit(1)); err != nil {
		return fmt.Errorf("kubernetes health check failed: %w", err)
	}
	return nil
}

// LoadEnvoyGatewayResources loads all Envoy Gateway resources for a given provider
// This is a convenience method that matches the test case requirements
func (m *Manager) LoadEnvoyGatewayResources(ctx context.Context, namespace, name string) ([]interface{}, error) {
	var resources []interface{}

	// Load Backend
	if backend, err := m.Backend.Get(ctx, namespace, name); err == nil {
		resources = append(resources, backend)
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to load Backend: %w", err)
	}

	// Load BackendTLSPolicy
	if tlsPolicy, err := m.BackendTLSPolicy.Get(ctx, namespace, name); err == nil {
		resources = append(resources, tlsPolicy)
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to load BackendTLSPolicy: %w", err)
	}

	// Load Secret
	if secret, err := m.Secret.Get(ctx, namespace, name); err == nil {
		resources = append(resources, secret)
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to load Secret: %w", err)
	}

	// Load BackendSecurityPolicy
	if bsp, err := m.BackendSecurityPolicy.Get(ctx, namespace, name); err == nil {
		resources = append(resources, bsp)
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to load BackendSecurityPolicy: %w", err)
	}

	// Load AIServiceBackend
	if aisb, err := m.AIServiceBackend.Get(ctx, namespace, name); err == nil {
		resources = append(resources, aisb)
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to load AIServiceBackend: %w", err)
	}

	return resources, nil
}

// GetBackendClient returns the Backend client
func (m *Manager) GetBackendClient() BackendClientInterface {
	return m.Backend
}

// GetSecretClient returns the Secret client
func (m *Manager) GetSecretClient() SecretClientInterface {
	return m.Secret
}

// GetAIServiceBackendClient returns the AIServiceBackend client
func (m *Manager) GetAIServiceBackendClient() AIServiceBackendClientInterface {
	return m.AIServiceBackend
}

// GetBackendSecurityPolicyClient returns the BackendSecurityPolicy client
func (m *Manager) GetBackendSecurityPolicyClient() BackendSecurityPolicyClientInterface {
	return m.BackendSecurityPolicy
}

// GetBackendTLSPolicyClient returns the BackendTLSPolicy client
func (m *Manager) GetBackendTLSPolicyClient() BackendTLSPolicyClientInterface {
	return m.BackendTLSPolicy
}
