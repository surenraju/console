// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package client

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	aigv1a1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	gwapiv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

func TestManager_NewManager(t *testing.T) {
	// Create a test scheme with all required types
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, gwapiv1a1.AddToScheme(scheme))
	require.NoError(t, aigv1a1.AddToScheme(scheme))
	require.NoError(t, gwapiv1a3.AddToScheme(scheme))

	// Create a fake client for testing
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()

	// Create manager with fake client
	manager := &Manager{
		client:                fakeClient,
		logger:                logr.Discard(),
		Backend:               NewBackendClient(fakeClient, logr.Discard()),
		BackendTLSPolicy:      NewBackendTLSPolicyClient(fakeClient, logr.Discard()),
		Secret:                NewSecretClient(fakeClient, logr.Discard()),
		BackendSecurityPolicy: NewBackendSecurityPolicyClient(fakeClient, logr.Discard()),
		AIServiceBackend:      NewAIServiceBackendClient(fakeClient, logr.Discard()),
	}

	// Test that all clients are properly initialized
	assert.NotNil(t, manager.Backend)
	assert.NotNil(t, manager.BackendTLSPolicy)
	assert.NotNil(t, manager.Secret)
	assert.NotNil(t, manager.BackendSecurityPolicy)
	assert.NotNil(t, manager.AIServiceBackend)

	// Test interface methods
	assert.NotNil(t, manager.GetBackendClient())
	assert.NotNil(t, manager.GetSecretClient())
	assert.NotNil(t, manager.GetAIServiceBackendClient())
	assert.NotNil(t, manager.GetBackendSecurityPolicyClient())
	assert.NotNil(t, manager.GetBackendTLSPolicyClient())

	// Test client access
	assert.Equal(t, fakeClient, manager.Client())
}

func TestManager_LoadEnvoyGatewayResources(t *testing.T) {
	// Create a test scheme with all required types
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, gwapiv1a1.AddToScheme(scheme))
	require.NoError(t, aigv1a1.AddToScheme(scheme))
	require.NoError(t, gwapiv1a3.AddToScheme(scheme))

	// Create a fake client for testing
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()

	// Create manager with fake client
	manager := &Manager{
		client:                fakeClient,
		logger:                logr.Discard(),
		Backend:               NewBackendClient(fakeClient, logr.Discard()),
		BackendTLSPolicy:      NewBackendTLSPolicyClient(fakeClient, logr.Discard()),
		Secret:                NewSecretClient(fakeClient, logr.Discard()),
		BackendSecurityPolicy: NewBackendSecurityPolicyClient(fakeClient, logr.Discard()),
		AIServiceBackend:      NewAIServiceBackendClient(fakeClient, logr.Discard()),
	}

	ctx := context.Background()

	// Test loading resources when none exist (should return empty slice, no error)
	resources, err := manager.LoadEnvoyGatewayResources(ctx, "test-namespace", "test-provider")
	require.NoError(t, err)
	assert.Empty(t, resources)
}

func TestManager_HealthCheck(t *testing.T) {
	// Create a test scheme with all required types
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	// Create a fake client for testing
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()

	// Create manager with fake client
	manager := &Manager{
		client: fakeClient,
		logger: logr.Discard(),
	}

	ctx := context.Background()

	// Test health check (should pass with fake client)
	err := manager.HealthCheck(ctx)
	assert.NoError(t, err)
}
