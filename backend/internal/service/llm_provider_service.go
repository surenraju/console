// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package service

import (
	"context"
	"fmt"

	aigatewayv1alpha1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	"github.com/envoyproxy/ai-gateway/console/backend/pkg/client"
	"github.com/envoyproxy/ai-gateway/console/backend/pkg/llm"
	gatewayv1alpha1 "github.com/envoyproxy/gateway/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

// LLMProviderService handles business logic for LLM providers
// It orchestrates loading Kubernetes resources and translating them to LLMProvider objects
type LLMProviderService struct {
	clientManager *client.Manager
}

// NewLLMProviderService creates a new LLMProviderService
func NewLLMProviderService(clientManager *client.Manager) *LLMProviderService {
	return &LLMProviderService{
		clientManager: clientManager,
	}
}

// ListProviders returns all available LLM providers
func (s *LLMProviderService) ListProviders(ctx context.Context, namespace string) ([]llm.LLMProvider, error) {
	backends, err := s.clientManager.AIServiceBackend.List(ctx, namespace)
	if err != nil {
		return []llm.LLMProvider{}, fmt.Errorf("failed to list AIServiceBackends: %w", err)
	}

	// Initialize with empty slice to ensure we never return nil
	providers := make([]llm.LLMProvider, 0)

	for _, backend := range backends.Items {
		resources, err := s.loadProviderResources(ctx, backend.Namespace, backend.Name)
		if err != nil {
			// Log error but continue with other providers
			continue
		}

		provider, err := llm.ToLLMProvider(resources)
		if err != nil {
			// Log error but continue with other providers
			continue
		}

		// Mask sensitive information before adding to the list
		maskedProvider := provider.MaskSecret()
		providers = append(providers, *maskedProvider)
	}

	return providers, nil
}

// GetProvider returns a specific LLM provider by namespace and name
func (s *LLMProviderService) GetProvider(ctx context.Context, namespace, name string) (*llm.LLMProvider, error) {
	resources, err := s.loadProviderResources(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	provider, err := llm.ToLLMProvider(resources)
	if err != nil {
		return nil, err
	}

	// Mask sensitive information before returning
	return provider.MaskSecret(), nil
}

// CreateProvider creates a new LLM provider by converting it to Kubernetes resources
func (s *LLMProviderService) CreateProvider(ctx context.Context, provider *llm.LLMProvider) error {
	// Convert LLMProvider to Kubernetes resources
	resources, err := provider.ToEnvoyGatewayResources()
	if err != nil {
		return fmt.Errorf("failed to convert provider to Kubernetes resources: %w", err)
	}

	// Create each resource in the cluster
	for _, resource := range resources {
		switch r := resource.(type) {
		case *gatewayv1alpha1.Backend:
			err = s.clientManager.Backend.Create(ctx, r)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return fmt.Errorf("backend '%s' already exists. Please choose a different name or delete the existing backend first", r.Name)
				}
				return fmt.Errorf("failed to create Backend: %w", err)
			}

		case *gwapiv1a3.BackendTLSPolicy:
			err = s.clientManager.BackendTLSPolicy.Create(ctx, r)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return fmt.Errorf("backend TLS policy '%s' already exists. Please choose a different name or delete the existing policy first", r.Name)
				}
				return fmt.Errorf("failed to create BackendTLSPolicy: %w", err)
			}

		case *aigatewayv1alpha1.BackendSecurityPolicy:
			err = s.clientManager.BackendSecurityPolicy.Create(ctx, r)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return fmt.Errorf("backend security policy '%s' already exists. Please choose a different name or delete the existing policy first", r.Name)
				}
				return fmt.Errorf("failed to create BackendSecurityPolicy: %w", err)
			}

		case *aigatewayv1alpha1.AIServiceBackend:
			err = s.clientManager.AIServiceBackend.Create(ctx, r)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return fmt.Errorf("AI service backend '%s' already exists. Please choose a different name or delete the existing provider first", r.Name)
				}
				return fmt.Errorf("failed to create AIServiceBackend: %w", err)
			}

		case *corev1.Secret:
			err = s.clientManager.Secret.Create(ctx, r)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					return fmt.Errorf("secret '%s' already exists. Please choose a different name or delete the existing secret first", r.Name)
				}
				return fmt.Errorf("failed to create Secret: %w", err)
			}

		default:
			return fmt.Errorf("unknown resource type: %T", r)
		}
	}

	return nil
}

// DeleteProvider deletes an LLM provider by removing all its Kubernetes resources
func (s *LLMProviderService) DeleteProvider(ctx context.Context, namespace, name string) error {
	// Load all resources for this provider first
	resources, err := s.loadProviderResources(ctx, namespace, name)
	if err != nil {
		return fmt.Errorf("failed to load provider resources for deletion: %w", err)
	}

	// Delete resources in reverse order to avoid dependency issues
	// Delete AIServiceBackend first (it references other resources)
	for _, resource := range resources {
		switch r := resource.(type) {
		case *aigatewayv1alpha1.AIServiceBackend:
			err = s.clientManager.AIServiceBackend.Delete(ctx, r.Namespace, r.Name)
			if err != nil {
				return fmt.Errorf("failed to delete AIServiceBackend: %w", err)
			}
		}
	}

	// Then delete BackendSecurityPolicy
	for _, resource := range resources {
		switch r := resource.(type) {
		case *aigatewayv1alpha1.BackendSecurityPolicy:
			err = s.clientManager.BackendSecurityPolicy.Delete(ctx, r.Namespace, r.Name)
			if err != nil {
				return fmt.Errorf("failed to delete BackendSecurityPolicy: %w", err)
			}
		}
	}

	// Then delete BackendTLSPolicy
	for _, resource := range resources {
		switch r := resource.(type) {
		case *gwapiv1a3.BackendTLSPolicy:
			err = s.clientManager.BackendTLSPolicy.Delete(ctx, r.Namespace, r.Name)
			if err != nil {
				return fmt.Errorf("failed to delete BackendTLSPolicy: %w", err)
			}
		}
	}

	// Then delete Backend
	for _, resource := range resources {
		switch r := resource.(type) {
		case *gatewayv1alpha1.Backend:
			err = s.clientManager.Backend.Delete(ctx, r.Namespace, r.Name)
			if err != nil {
				return fmt.Errorf("failed to delete Backend: %w", err)
			}
		}
	}

	// Finally delete Secrets
	for _, resource := range resources {
		switch r := resource.(type) {
		case *corev1.Secret:
			err = s.clientManager.Secret.Delete(ctx, r.Namespace, r.Name)
			if err != nil {
				return fmt.Errorf("failed to delete Secret: %w", err)
			}
		}
	}

	return nil
}

// loadProviderResources loads all resources for a specific provider hierarchically
func (s *LLMProviderService) loadProviderResources(ctx context.Context, namespace, name string) ([]interface{}, error) {
	var resources []interface{}

	// 1. First load AIServiceBackend
	aisb, err := s.clientManager.AIServiceBackend.Get(ctx, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get AIServiceBackend %s/%s: %w", namespace, name, err)
	}
	resources = append(resources, aisb)

	// 2. Based on backendRef.name load Backend
	backendName := string(aisb.Spec.BackendRef.Name)
	backendNamespace := namespace
	if aisb.Spec.BackendRef.Namespace != nil {
		backendNamespace = string(*aisb.Spec.BackendRef.Namespace)
	}

	backend, err := s.clientManager.Backend.Get(ctx, backendNamespace, backendName)
	if err != nil {
		return nil, fmt.Errorf("failed to get Backend %s/%s: %w", backendNamespace, backendName, err)
	}
	resources = append(resources, backend)

	// 5. Based on Backend find BackendTLSPolicy with matching targetRefs.name
	tlsPolicies, err := s.clientManager.BackendTLSPolicy.List(ctx, backendNamespace)
	if err == nil {
		for _, policy := range tlsPolicies.Items {
			// Check if this TLS policy targets our backend
			for _, targetRef := range policy.Spec.TargetRefs {
				if string(targetRef.Name) == backendName &&
					string(targetRef.Kind) == "Backend" &&
					(string(targetRef.Group) == "" || string(targetRef.Group) == "gateway.envoyproxy.io") {
					resources = append(resources, &policy)
					break
				}
			}
		}
	}

	// 3. Find BackendSecurityPolicy that targets this AIServiceBackend via targetRefs
	// This replaces the deprecated BackendSecurityPolicyRef field
	securityPolicies, err := s.clientManager.BackendSecurityPolicy.List(ctx, namespace)
	if err == nil {
		for _, policy := range securityPolicies.Items {
			// Check if this security policy targets our AIServiceBackend
			for _, targetRef := range policy.Spec.TargetRefs {
				if string(targetRef.Name) == aisb.Name &&
					string(targetRef.Kind) == "AIServiceBackend" &&
					(string(targetRef.Group) == "" || string(targetRef.Group) == "aigateway.envoyproxy.io") {
					resources = append(resources, &policy)

					// 4. Based on BackendSecurityPolicy.secretRef find secret
					err = s.loadSecretsForAuthType(ctx, &policy, namespace, &resources)
					if err != nil {
						return nil, fmt.Errorf("failed to load secrets for auth type: %w", err)
					}
					break
				}
			}
		}
	}

	return resources, nil
}

// loadSecretsForAuthType loads secrets based on the authentication type
func (s *LLMProviderService) loadSecretsForAuthType(ctx context.Context, securityPolicy interface{}, namespace string, resources *[]interface{}) error {
	// Type assertion to get the actual BackendSecurityPolicy
	bsp, ok := securityPolicy.(*aigatewayv1alpha1.BackendSecurityPolicy)
	if !ok {
		return fmt.Errorf("invalid security policy type")
	}

	switch bsp.Spec.Type {
	case aigatewayv1alpha1.BackendSecurityPolicyTypeAPIKey:
		// Load API key secret
		if bsp.Spec.APIKey != nil && bsp.Spec.APIKey.SecretRef != nil {
			secretName := string(bsp.Spec.APIKey.SecretRef.Name)
			secretNamespace := namespace
			if bsp.Spec.APIKey.SecretRef.Namespace != nil {
				secretNamespace = string(*bsp.Spec.APIKey.SecretRef.Namespace)
			}

			secret, err := s.clientManager.Secret.Get(ctx, secretNamespace, secretName)
			if err == nil {
				*resources = append(*resources, secret)
			}
		}

	case aigatewayv1alpha1.BackendSecurityPolicyTypeGCPCredentials:
		// Load GCP client secret
		if bsp.Spec.GCPCredentials != nil {
			clientSecret := bsp.Spec.GCPCredentials.WorkloadIdentityFederationConfig.OIDCExchangeToken.OIDC.ClientSecret
			secretName := string(clientSecret.Name)
			secretNamespace := namespace
			if clientSecret.Namespace != nil {
				secretNamespace = string(*clientSecret.Namespace)
			}

			secret, err := s.clientManager.Secret.Get(ctx, secretNamespace, secretName)
			if err == nil {
				*resources = append(*resources, secret)
			}
		}

	case aigatewayv1alpha1.BackendSecurityPolicyTypeAWSCredentials:
		// Load AWS credentials secret
		if bsp.Spec.AWSCredentials != nil &&
			bsp.Spec.AWSCredentials.CredentialsFile != nil &&
			bsp.Spec.AWSCredentials.CredentialsFile.SecretRef != nil {
			secretName := string(bsp.Spec.AWSCredentials.CredentialsFile.SecretRef.Name)
			secretNamespace := namespace
			if bsp.Spec.AWSCredentials.CredentialsFile.SecretRef.Namespace != nil {
				secretNamespace = string(*bsp.Spec.AWSCredentials.CredentialsFile.SecretRef.Namespace)
			}

			secret, err := s.clientManager.Secret.Get(ctx, secretNamespace, secretName)
			if err == nil {
				*resources = append(*resources, secret)
			}
		}

	case aigatewayv1alpha1.BackendSecurityPolicyTypeAzureCredentials:
		// Load Azure client secret
		if bsp.Spec.AzureCredentials != nil && bsp.Spec.AzureCredentials.ClientSecretRef != nil {
			secretName := string(bsp.Spec.AzureCredentials.ClientSecretRef.Name)
			secretNamespace := namespace
			if bsp.Spec.AzureCredentials.ClientSecretRef.Namespace != nil {
				secretNamespace = string(*bsp.Spec.AzureCredentials.ClientSecretRef.Namespace)
			}

			secret, err := s.clientManager.Secret.Get(ctx, secretNamespace, secretName)
			if err == nil {
				*resources = append(*resources, secret)
			}
		}
	}

	return nil
}
