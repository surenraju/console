package tests

import (
	"testing"

	"github.com/envoyproxy/ai-gateway/console/backend/pkg/llm"
)

func TestLLMProviderMaskSecret(t *testing.T) {
	provider := &llm.LLMProvider{
		Name:      "test-provider",
		Namespace: "default",
		Schema:    "OpenAI",
		Version:   "v1",
		Auth: llm.AuthConfig{
			Type:   "apiKey",
			APIKey: "sk-super-secret-key",
		},
		Backend: llm.Backend{
			Host: "api.openai.com",
			Port: 443,
		},
		TLS: llm.TLSValidation{
			Hostname:                "api.openai.com",
			WellKnownCACertificates: "System",
		},
	}

	masked := provider.MaskSecret()

	// Verify non-sensitive fields are preserved
	if masked.Name != provider.Name {
		t.Errorf("Expected name %s, got %s", provider.Name, masked.Name)
	}
	if masked.Namespace != provider.Namespace {
		t.Errorf("Expected namespace %s, got %s", provider.Namespace, masked.Namespace)
	}
	if masked.Schema != provider.Schema {
		t.Errorf("Expected schema %s, got %s", provider.Schema, masked.Schema)
	}
	if masked.Backend.Host != provider.Backend.Host {
		t.Errorf("Expected backend host %s, got %s", provider.Backend.Host, masked.Backend.Host)
	}

	// Verify sensitive fields are masked
	if masked.Auth.APIKey != llm.MaskedSecretValue {
		t.Errorf("Expected API key to be masked, got %s", masked.Auth.APIKey)
	}

	// Verify original is unchanged
	if provider.Auth.APIKey == llm.MaskedSecretValue {
		t.Error("Original provider was modified during masking")
	}
}

func TestAWSAuthMaskSecret(t *testing.T) {
	aws := &llm.AWSAuth{
		Region:          "us-east-1",
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	masked := aws.MaskSecret()

	// Verify non-sensitive field is preserved
	if masked.Region != aws.Region {
		t.Errorf("Expected region %s, got %s", aws.Region, masked.Region)
	}

	// Verify sensitive fields are masked
	if masked.AccessKeyID != llm.MaskedSecretValue {
		t.Errorf("Expected access key ID to be masked, got %s", masked.AccessKeyID)
	}
	if masked.SecretAccessKey != llm.MaskedSecretValue {
		t.Errorf("Expected secret access key to be masked, got %s", masked.SecretAccessKey)
	}

	// Verify original is unchanged
	if aws.AccessKeyID == llm.MaskedSecretValue || aws.SecretAccessKey == llm.MaskedSecretValue {
		t.Error("Original AWS auth was modified during masking")
	}
}

func TestGCPAuthMaskSecret(t *testing.T) {
	gcp := &llm.GCPAuth{
		ProjectID:                    "my-project",
		Location:                     "us-central1",
		WorkloadIdentityPoolName:     "my-pool",
		WorkloadIdentityProviderName: "my-provider",
		ServiceAccountName:           "my-service-account",
		OIDCIssuer:                   "https://token.actions.githubusercontent.com",
		OIDCClientID:                 "my-client-id",
		OIDCClientSecret:             "super-secret-client-secret",
		PrivateKey:                   "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
	}

	masked := gcp.MaskSecret()

	// Verify non-sensitive fields are preserved
	if masked.ProjectID != gcp.ProjectID {
		t.Errorf("Expected project ID %s, got %s", gcp.ProjectID, masked.ProjectID)
	}
	if masked.Location != gcp.Location {
		t.Errorf("Expected location %s, got %s", gcp.Location, masked.Location)
	}
	if masked.OIDCClientID != gcp.OIDCClientID {
		t.Errorf("Expected OIDC client ID %s, got %s", gcp.OIDCClientID, masked.OIDCClientID)
	}

	// Verify sensitive fields are masked
	if masked.OIDCClientSecret != llm.MaskedSecretValue {
		t.Errorf("Expected OIDC client secret to be masked, got %s", masked.OIDCClientSecret)
	}
	if masked.PrivateKey != llm.MaskedSecretValue {
		t.Errorf("Expected private key to be masked, got %s", masked.PrivateKey)
	}

	// Verify original is unchanged
	if gcp.OIDCClientSecret == llm.MaskedSecretValue || gcp.PrivateKey == llm.MaskedSecretValue {
		t.Error("Original GCP auth was modified during masking")
	}
}

func TestAzureAuthMaskSecret(t *testing.T) {
	azure := &llm.AzureAuth{
		ClientID: "client-id-123",
		TenantID: "tenant-id-456",
		APIKey:   "super-secret-api-key",
	}

	masked := azure.MaskSecret()

	// Verify non-sensitive fields are preserved
	if masked.ClientID != azure.ClientID {
		t.Errorf("Expected client ID %s, got %s", azure.ClientID, masked.ClientID)
	}
	if masked.TenantID != azure.TenantID {
		t.Errorf("Expected tenant ID %s, got %s", azure.TenantID, masked.TenantID)
	}

	// Verify sensitive field is masked
	if masked.APIKey != llm.MaskedSecretValue {
		t.Errorf("Expected API key to be masked, got %s", masked.APIKey)
	}

	// Verify original is unchanged
	if azure.APIKey == llm.MaskedSecretValue {
		t.Error("Original Azure auth was modified during masking")
	}
}

func TestNilMaskSecret(t *testing.T) {
	// Test that nil pointers are handled gracefully
	var provider *llm.LLMProvider
	if masked := provider.MaskSecret(); masked != nil {
		t.Error("Expected nil when masking nil provider")
	}

	var aws *llm.AWSAuth
	if masked := aws.MaskSecret(); masked != nil {
		t.Error("Expected nil when masking nil AWS auth")
	}

	var gcp *llm.GCPAuth
	if masked := gcp.MaskSecret(); masked != nil {
		t.Error("Expected nil when masking nil GCP auth")
	}

	var azure *llm.AzureAuth
	if masked := azure.MaskSecret(); masked != nil {
		t.Error("Expected nil when masking nil Azure auth")
	}
}
