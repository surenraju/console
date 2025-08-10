package llm

// LLMProvider represents a logical configuration to interact with LLM APIs via Envoy AI Gateway.
type LLMProvider struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	Schema  string `json:"schema"` // supported schema OpenAI, AWSBedrock, AzureOpenAI, GCPVertexAI
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

const MaskedSecretValue = "***MASKED***"

// MaskSecret returns a copy of the LLMProvider with all sensitive information masked.
func (l *LLMProvider) MaskSecret() *LLMProvider {
	if l == nil {
		return nil
	}

	masked := &LLMProvider{
		Name:      l.Name,
		Namespace: l.Namespace,
		Schema:    l.Schema,
		Version:   l.Version,
		Auth:      l.Auth.MaskSecret(),
		Backend:   l.Backend,
		TLS:       l.TLS,
	}

	return masked
}

// MaskSecret returns a copy of the AuthConfig with all sensitive information masked.
func (a AuthConfig) MaskSecret() AuthConfig {
	masked := AuthConfig{
		Type:      a.Type,
		SecretRef: a.SecretRef, // SecretRef contains only references, not actual secrets
	}

	// Mask direct credentials
	if a.APIKey != "" {
		masked.APIKey = MaskedSecretValue
	}

	if a.AWS != nil {
		masked.AWS = a.AWS.MaskSecret()
	}

	if a.GCP != nil {
		masked.GCP = a.GCP.MaskSecret()
	}

	if a.Azure != nil {
		masked.Azure = a.Azure.MaskSecret()
	}

	return masked
}

// MaskSecret returns a copy of the AWSAuth with all sensitive information masked.
func (a *AWSAuth) MaskSecret() *AWSAuth {
	if a == nil {
		return nil
	}

	masked := &AWSAuth{
		Region: a.Region, // Region is not sensitive
	}

	if a.AccessKeyID != "" {
		masked.AccessKeyID = MaskedSecretValue
	}

	if a.SecretAccessKey != "" {
		masked.SecretAccessKey = MaskedSecretValue
	}

	return masked
}

// MaskSecret returns a copy of the GCPAuth with all sensitive information masked.
func (g *GCPAuth) MaskSecret() *GCPAuth {
	if g == nil {
		return nil
	}

	masked := &GCPAuth{
		// Non-sensitive configuration fields
		ProjectID:                    g.ProjectID,
		Location:                     g.Location,
		WorkloadIdentityPoolName:     g.WorkloadIdentityPoolName,
		WorkloadIdentityProviderName: g.WorkloadIdentityProviderName,
		ServiceAccountName:           g.ServiceAccountName,
		OIDCIssuer:                   g.OIDCIssuer,
		OIDCClientID:                 g.OIDCClientID,
		ClientEmail:                  g.ClientEmail,
		ServiceAccountProjectID:      g.ServiceAccountProjectID,
		ClientID:                     g.ClientID,
		AuthURI:                      g.AuthURI,
		TokenURI:                     g.TokenURI,
	}

	// Mask sensitive fields
	if g.OIDCClientSecret != "" {
		masked.OIDCClientSecret = MaskedSecretValue
	}

	if g.PrivateKey != "" {
		masked.PrivateKey = MaskedSecretValue
	}

	return masked
}

// MaskSecret returns a copy of the AzureAuth with all sensitive information masked.
func (a *AzureAuth) MaskSecret() *AzureAuth {
	if a == nil {
		return nil
	}

	masked := &AzureAuth{
		ClientID: a.ClientID, // ClientID is not sensitive
		TenantID: a.TenantID, // TenantID is not sensitive
	}

	if a.APIKey != "" {
		masked.APIKey = MaskedSecretValue
	}

	return masked
}
