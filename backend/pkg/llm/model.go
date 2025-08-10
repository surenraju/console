package llm

// LLMProvider represents a logical configuration to interact with LLM APIs via Envoy AI Gateway.
type LLMProvider struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	Schema  string `json:"schema"` // OpenAI, AWS, AzureOpenAI, GCP
	Version string `json:"version,omitempty"`

	Auth AuthConfig `json:"auth"`

	Backend Backend       `json:"backend"`
	TLS     TLSValidation `json:"tls"`
}

// AuthType represents the type of authentication used.
type AuthType string

// AuthConfig supports multiple authentication strategies.
type AuthConfig struct {
	Type string `json:"type"` // APIKey, AWS, Azure, GCP (use string for flexible matching)

	// Generic secret reference (for APIKey, or JSON creds)
	SecretRef *SecretRef `json:"secretRef,omitempty"`

	// Direct credentials if not using secretRef
	APIKey string     `json:"apiKey,omitempty"` // for OpenAI / Azure (not pointer)
	AWS    *AWSAuth   `json:"aws,omitempty"`
	GCP    *GCPAuth   `json:"gcp,omitempty"`
	Azure  *AzureAuth `json:"azure,omitempty"`
}

// SecretRef is a simplified form of corev1.SecretReference.
type SecretRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// AWSAuth defines raw credentials for AWS Bedrock.
type AWSAuth struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// GCPAuth defines configuration for GCP Vertex AI using workload identity federation.
// Fields are named to match the structure needed for BackendSecurityPolicy with GCPCredentials
type GCPAuth struct {
	// Main project configuration
	ProjectID string `json:"projectId"` // Used as both projectName in BackendSecurityPolicy and projectID in workloadIdentityFederationConfig
	Location  string `json:"location"`  // Used as region in BackendSecurityPolicy

	// Workload identity federation configuration
	WorkloadIdentityPoolName     string `json:"workloadIdentityPoolName"`     // Pool name for workload identity federation
	WorkloadIdentityProviderName string `json:"workloadIdentityProviderName"` // Provider name registered in GCP
	ServiceAccountName           string `json:"serviceAccountName"`           // Service account to impersonate in GCP

	// OIDC configuration for token exchange
	OIDCIssuer       string `json:"oidcIssuer"`       // OIDC provider issuer URL (e.g., "https://token.actions.githubusercontent.com")
	OIDCClientID     string `json:"oidcClientId"`     // OIDC client ID for authentication with the identity provider
	OIDCClientSecret string `json:"oidcClientSecret"` // OIDC client secret - will be stored in a separate Secret resource

	// Legacy fields for backward compatibility
	PrivateKey              string `json:"privateKey,omitempty"`              // Legacy field - service account key JSON
	ClientEmail             string `json:"clientEmail,omitempty"`             // Legacy field - fallback for WorkloadIdentityPoolName
	ServiceAccountProjectID string `json:"serviceAccountProjectId,omitempty"` // Legacy field - fallback for WorkloadIdentityProviderName
	ClientID                string `json:"clientId,omitempty"`                // Legacy field - fallback for ServiceAccountName
	AuthURI                 string `json:"authUri,omitempty"`                 // Legacy field - fallback for OIDCIssuer
	TokenURI                string `json:"tokenUri,omitempty"`                // Legacy field - fallback for OIDCClientID
}

// AzureAuth defines direct input for Azure OpenAI (if not using secret).
type AzureAuth struct {
	ClientID string `json:"clientId,omitempty"`
	TenantID string `json:"tenantId,omitempty"`
	APIKey   string `json:"apiKey,omitempty"` // some use cases use apikey over header
}

// Backend represents the network address of the provider's API.
type Backend struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

// TLSValidation specifies the upstream TLS validation settings.
type TLSValidation struct {
	Hostname                string `json:"hostname"`
	WellKnownCACertificates string `json:"wellKnownCACertificates"` // e.g., "System"
}
