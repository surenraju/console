// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"

	aigv1a1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	gwapiv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

// BackendClientInterface defines the interface for Backend operations
type BackendClientInterface interface {
	Create(ctx context.Context, backend *gwapiv1a1.Backend) error
	Get(ctx context.Context, namespace, name string) (*gwapiv1a1.Backend, error)
	List(ctx context.Context, namespace string) (*gwapiv1a1.BackendList, error)
	Update(ctx context.Context, backend *gwapiv1a1.Backend) error
	Delete(ctx context.Context, namespace, name string) error
}

// SecretClientInterface defines the interface for Secret operations
type SecretClientInterface interface {
	Create(ctx context.Context, secret *corev1.Secret) error
	Get(ctx context.Context, namespace, name string) (*corev1.Secret, error)
	List(ctx context.Context, namespace string) (*corev1.SecretList, error)
	Update(ctx context.Context, secret *corev1.Secret) error
	Delete(ctx context.Context, namespace, name string) error
}

// AIServiceBackendClientInterface defines the interface for AIServiceBackend operations
type AIServiceBackendClientInterface interface {
	Create(ctx context.Context, backend *aigv1a1.AIServiceBackend) error
	Get(ctx context.Context, namespace, name string) (*aigv1a1.AIServiceBackend, error)
	List(ctx context.Context, namespace string) (*aigv1a1.AIServiceBackendList, error)
	Update(ctx context.Context, backend *aigv1a1.AIServiceBackend) error
	Delete(ctx context.Context, namespace, name string) error
}

// BackendSecurityPolicyClientInterface defines the interface for BackendSecurityPolicy operations
type BackendSecurityPolicyClientInterface interface {
	Create(ctx context.Context, policy *aigv1a1.BackendSecurityPolicy) error
	Get(ctx context.Context, namespace, name string) (*aigv1a1.BackendSecurityPolicy, error)
	List(ctx context.Context, namespace string) (*aigv1a1.BackendSecurityPolicyList, error)
	Update(ctx context.Context, policy *aigv1a1.BackendSecurityPolicy) error
	Delete(ctx context.Context, namespace, name string) error
}

// BackendTLSPolicyClientInterface defines the interface for BackendTLSPolicy operations
type BackendTLSPolicyClientInterface interface {
	Create(ctx context.Context, policy *gwapiv1a3.BackendTLSPolicy) error
	Get(ctx context.Context, namespace, name string) (*gwapiv1a3.BackendTLSPolicy, error)
	List(ctx context.Context, namespace string) (*gwapiv1a3.BackendTLSPolicyList, error)
	Update(ctx context.Context, policy *gwapiv1a3.BackendTLSPolicy) error
	Delete(ctx context.Context, namespace, name string) error
}

// ManagerInterface defines the interface for the client manager
type ManagerInterface interface {
	// Client returns the underlying Kubernetes client
	Client() client.Client

	// HealthCheck performs a basic health check against the Kubernetes API
	HealthCheck(ctx context.Context) error

	// LoadEnvoyGatewayResources loads all Envoy Gateway resources for a given provider
	LoadEnvoyGatewayResources(ctx context.Context, namespace, name string) ([]interface{}, error)

	// Resource clients
	GetBackendClient() BackendClientInterface
	GetSecretClient() SecretClientInterface
	GetAIServiceBackendClient() AIServiceBackendClientInterface
	GetBackendSecurityPolicyClient() BackendSecurityPolicyClientInterface
	GetBackendTLSPolicyClient() BackendTLSPolicyClientInterface
}

// Ensure our implementations satisfy the interfaces
var _ BackendClientInterface = &BackendClient{}
var _ SecretClientInterface = &SecretClient{}
var _ AIServiceBackendClientInterface = &AIServiceBackendClient{}
var _ BackendSecurityPolicyClientInterface = &BackendSecurityPolicyClient{}
var _ BackendTLSPolicyClientInterface = &BackendTLSPolicyClient{}
var _ ManagerInterface = &Manager{}
