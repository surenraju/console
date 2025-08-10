// Types matching the backend Go structs

export interface LLMProvider {
  name: string;
  schema: string; // OpenAI, AWSBedrock, AzureOpenAI, GCPVertexAI
  version?: string;
  auth: AuthConfig;
  backend: Backend;
  tls: TLSValidation;
}

export interface AuthConfig {
  type: string; // APIKey, AWS, Azure, GCP
  secretRef?: SecretRef;
  apiKey?: string;
  aws?: AWSAuth;
  gcp?: GCPAuth;
  azure?: AzureAuth;
}

export interface SecretRef {
  name: string;
  namespace: string;
}

export interface AWSAuth {
  region: string;
  accessKeyId?: string;
  secretAccessKey?: string;
}

export interface GCPAuth {
  projectId: string;
  location: string;
  workloadIdentityPoolName: string;
  workloadIdentityProviderName: string;
  serviceAccountName: string;
  oidcIssuer: string;
  oidcClientId: string;
  oidcClientSecret: string;
  // Legacy fields for backward compatibility
  privateKey?: string;
  clientEmail?: string;
  serviceAccountProjectId?: string;
  clientId?: string;
  authUri?: string;
  tokenUri?: string;
}

export interface AzureAuth {
  clientId?: string;
  tenantId?: string;
  apiKey?: string;
}

export interface Backend {
  host: string;
  port: number;
}

export interface TLSValidation {
  hostname: string;
  wellKnownCACertificates: string; // e.g., "System"
}

// Simplified interface for creating new providers (frontend form)
export interface CreateLLMProviderRequest {
  name: string;
  schema: string; // OpenAI, AWSBedrock, AzureOpenAI, GCPVertexAI
  version?: string;
  
  // Simplified auth - will be transformed to full AuthConfig
  authType: string; // APIKey, AWS, Azure, GCP
  
  // Direct credentials for APIKey (OpenAI, Azure)
  apiKey?: string;
  
  // AWS specific
  awsRegion?: string;
  awsAccessKeyId?: string;
  awsSecretAccessKey?: string;
  
  // GCP specific - Workload Identity Federation
  gcpProjectId?: string;
  gcpLocation?: string;
  gcpWorkloadIdentityPoolName?: string;
  gcpWorkloadIdentityProviderName?: string;
  gcpServiceAccountName?: string;
  gcpOidcIssuer?: string;
  gcpOidcClientId?: string;
  gcpOidcClientSecret?: string;
  
  // GCP Legacy fields for backward compatibility
  gcpPrivateKey?: string;
  gcpClientEmail?: string;
  gcpServiceAccountProjectId?: string;
  gcpClientId?: string;
  gcpAuthUri?: string;
  gcpTokenUri?: string;
  
  // Azure specific
  azureClientId?: string;
  azureTenantId?: string;
  azureApiKey?: string;
  
  // Backend
  host: string;
  port: number;
  
  // TLS
  tlsHostname?: string;
  tlsWellKnownCACertificates?: string;
}

// Helper functions to transform between frontend and backend formats
export function createLLMProviderFromForm(form: CreateLLMProviderRequest): LLMProvider {
  const authConfig: AuthConfig = { type: form.authType };
  
  // Handle direct credentials
  if (form.authType === 'APIKey' && form.apiKey) {
    authConfig.apiKey = form.apiKey;
  } else if (form.authType === 'AWS') {
    authConfig.aws = {
      region: form.awsRegion || '',
      accessKeyId: form.awsAccessKeyId || '',
      secretAccessKey: form.awsSecretAccessKey || '',
    };
  } else if (form.authType === 'GCP') {
    authConfig.gcp = {
      projectId: form.gcpProjectId || '',
      location: form.gcpLocation || '',
      workloadIdentityPoolName: form.gcpWorkloadIdentityPoolName || '',
      workloadIdentityProviderName: form.gcpWorkloadIdentityProviderName || '',
      serviceAccountName: form.gcpServiceAccountName || '',
      oidcIssuer: form.gcpOidcIssuer || '',
      oidcClientId: form.gcpOidcClientId || '',
      oidcClientSecret: form.gcpOidcClientSecret || '',
      // Legacy fields
      privateKey: form.gcpPrivateKey,
      clientEmail: form.gcpClientEmail,
      serviceAccountProjectId: form.gcpServiceAccountProjectId,
      clientId: form.gcpClientId,
      authUri: form.gcpAuthUri,
      tokenUri: form.gcpTokenUri,
    };
  } else if (form.authType === 'Azure') {
    authConfig.azure = {
      clientId: form.azureClientId || '',
      tenantId: form.azureTenantId || '',
      apiKey: form.azureApiKey || '',
    };
  }

  return {
    name: form.name,
    schema: form.schema,
    version: form.version,
    auth: authConfig,
    backend: {
      host: form.host,
      port: form.port,
    },
    tls: {
      hostname: form.tlsHostname || form.host,
      wellKnownCACertificates: form.tlsWellKnownCACertificates || 'System',
    },
  };
}

// Display helpers for the UI
export interface LLMProviderDisplay {
  name: string;
  type: string; // schema
  model: string; // derived from schema/version
  endpoint: string; // derived from backend
  status: 'active' | 'inactive' | 'error'; // will need to be determined by the frontend
  authType: string;
  createdAt?: string; // if available from metadata
  lastUsed?: string; // if available from metrics
}

export function toLLMProviderDisplay(provider: LLMProvider): LLMProviderDisplay {
  const endpoint = `${provider.backend.host}:${provider.backend.port}`;
  const model = provider.version ? `${provider.schema}-${provider.version}` : provider.schema;
  
  return {
    name: provider.name,
    type: provider.schema,
    model,
    endpoint,
    status: 'active', // TODO: determine actual status
    authType: provider.auth.type,
    // createdAt and lastUsed would come from Kubernetes metadata or metrics
  };
}
