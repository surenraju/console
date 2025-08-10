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
		return nil, fmt.Errorf("failed to list AIServiceBackends: %w", err)
	}

	var providers []llm.LLMProvider
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

		providers = append(providers, *provider)
	}

	return providers, nil
}

// GetProvider returns a specific LLM provider by namespace and name
func (s *LLMProviderService) GetProvider(ctx context.Context, namespace, name string) (*llm.LLMProvider, error) {
	resources, err := s.loadProviderResources(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	return llm.ToLLMProvider(resources)
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
