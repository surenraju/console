// Package llm contains internal data model for LLMProvider
package llm

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"

	aigwv1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	egv1 "github.com/envoyproxy/gateway/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// LLMProvider is the internal data model used by the application.
type LLMProvider struct {
	Name      string
	Namespace string
	Schema    string // OpenAI, AWS, AzureOpenAI, GCP
	Version   string

	APIKeyRef     *corev1.SecretReference
	AWSConfig     *AWSConfig
	GCPConfig     *GCPConfig
	AzureClientID string
	AzureTenantID string

	Backend   Backend
	TLS       TLSValidation
	AuthType  string // APIKey, AWS, Azure, GCP
}

type AWSConfig struct {
	Region         string
	CredentialsRef *corev1.SecretReference
}

type GCPConfig struct {
	ProjectID      string
	Location       string
	CredentialsRef *corev1.SecretReference
}

type Backend struct {
	Host string
	Port int32
}

type TLSValidation struct {
	Hostname                string
	WellKnownCACertificates string
}

func (p *LLMProvider) ToEnvoyGatewayResources() ([]client.Object, error) {
	resources := []client.Object{}

	// Backend
	backend := &egv1.Backend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: egv1.BackendSpec{
			Endpoints: []egv1.BackendEndpoint{{
				FQDN: &egv1.BackendFQDN{
					Hostname: p.Backend.Host,
					Port:     p.Backend.Port,
				},
			}},
		},
	}
	resources = append(resources, backend)

	// BackendTLSPolicy
	backendTLS := &egv1.BackendTLSPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-tls", p.Name),
			Namespace: p.Namespace,
		},
		Spec: egv1.BackendTLSPolicySpec{
			TargetRefs: []egv1.LocalPolicyTargetReference{{
				Group: "gateway.envoyproxy.io",
				Kind:  "Backend",
				Name:  p.Name,
			}},
			Validation: &egv1.BackendTLSValidation{
				Hostname:                p.TLS.Hostname,
				WellKnownCACertificates: p.TLS.WellKnownCACertificates,
			},
		},
	}
	resources = append(resources, backendTLS)

	// Secret & BackendSecurityPolicy
	var bsp *aigwv1.BackendSecurityPolicy
	if p.AuthType == "APIKey" && p.APIKeyRef != nil {
		bsp = &aigwv1.BackendSecurityPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-auth", p.Name),
				Namespace: p.Namespace,
			},
			Spec: aigwv1.BackendSecurityPolicySpec{
				Type: aigwv1.APIKey,
				APIKey: &aigwv1.APIKeyAuth{
					SecretRef: &aigwv1.SecretObjectReference{
						Name:      p.APIKeyRef.Name,
						Namespace: p.APIKeyRef.Namespace,
					},
				},
			},
		}
	} else if p.AuthType == "AWS" && p.AWSConfig != nil {
		bsp = &aigwv1.BackendSecurityPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-auth", p.Name),
				Namespace: p.Namespace,
			},
			Spec: aigwv1.BackendSecurityPolicySpec{
				Type: aigwv1.AWSCredentials,
				AWSCredentials: &aigwv1.AWSCredentialsAuth{
					Region: p.AWSConfig.Region,
					Credentials: &aigwv1.SecretObjectReference{
						Name:      p.AWSConfig.CredentialsRef.Name,
						Namespace: p.AWSConfig.CredentialsRef.Namespace,
					},
				},
			},
		}
	} else if p.AuthType == "GCP" && p.GCPConfig != nil {
		bsp = &aigwv1.BackendSecurityPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-auth", p.Name),
				Namespace: p.Namespace,
			},
			Spec: aigwv1.BackendSecurityPolicySpec{
				Type: aigwv1.GCPCredentials,
				GCPCredentials: &aigwv1.GCPCredentialsAuth{
					ProjectName: p.GCPConfig.ProjectID,
					Region:      p.GCPConfig.Location,
					WorkloadIdentityFederationConfig: &aigwv1.WorkloadIdentityFederationConfig{
						// Fill in if required
					},
				},
			},
		}
	} else if p.AuthType == "Azure" && p.APIKeyRef != nil {
		bsp = &aigwv1.BackendSecurityPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-auth", p.Name),
				Namespace: p.Namespace,
			},
			Spec: aigwv1.BackendSecurityPolicySpec{
				Type: aigwv1.APIKey,
				APIKey: &aigwv1.APIKeyAuth{
					SecretRef: &aigwv1.SecretObjectReference{
						Name:      p.APIKeyRef.Name,
						Namespace: p.APIKeyRef.Namespace,
					},
				},
			},
		}
	}
	if bsp != nil {
		resources = append(resources, bsp)
	}

	// AIServiceBackend
	aisb := &aigwv1.AIServiceBackend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Spec: aigwv1.AIServiceBackendSpec{
			Schema: aigwv1.APISchema{Name: p.Schema, Version: p.Version},
			BackendRef: aigwv1.BackendReference{
				Group: "gateway.envoyproxy.io",
				Kind:  "Backend",
				Name:  p.Name,
			},
			BackendSecurityPolicyRef: &aigwv1.NamespacedObjectReference{
				Group: "aigateway.envoyproxy.io",
				Kind:  "BackendSecurityPolicy",
				Name:  fmt.Sprintf("%s-auth", p.Name),
			},
		},
	}
	resources = append(resources, isb)

	return resources, nil
}

func FromEnvoyGatewayResources(resources []client.Object) (*LLMProvider, error) {
	// TODO: Implement parsing logic to map CRDs back to LLMProvider struct.
	return nil, fmt.Errorf("not yet implemented")
}
