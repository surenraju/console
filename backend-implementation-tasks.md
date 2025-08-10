# Envoy AI Gateway Console - Backend Implementation Tasks

## Overview
This document outlines the step-by-step tasks for implementing the backend service for the Envoy AI Gateway Console. The backend will provide REST APIs and WebSocket connectivity for managing AI Gateway CRDs in Kubernetes.

## Architecture Summary

Based on our analysis of the AI Gateway CRDs, the backend needs to handle:

### Core CRDs
- **AIGatewayRoute** (`aigatewayroutes.aigateway.envoyproxy.io`)
- **AIServiceBackend** (`aiservicebackends.aigateway.envoyproxy.io`) 
- **BackendSecurityPolicy** (`backendsecuritypolicies.aigateway.envoyproxy.io`)
- **Backend** (`backends.gateway.envoyproxy.io`)
- **BackendTLSPolicy** (`backendtlspolicies.gateway.envoyproxy.io`)

### Provider Templates
- OpenAI (`api.openai.com`)
- AWS Bedrock (`bedrock-runtime.{region}.amazonaws.com`)
- Google Vertex AI (`{region}-aiplatform.googleapis.com`)
- Azure OpenAI (`{account}.openai.azure.com`)

## Implementation Tasks

### 🏗️ Task 1: Setup Backend Project Structure
**Status:** Pending
**Dependencies:** None

#### Actions:
1. Create Go module structure:
   ```
   backend/
   ├── cmd/server/main.go
   ├── internal/
   │   ├── api/         # REST API handlers
   │   ├── handlers/    # Business logic handlers  
   │   ├── server/      # HTTP/WebSocket server
   │   ├── k8s/         # Kubernetes client & watchers
   │   ├── models/      # CRD Go structs
   │   └── config/      # Configuration management
   ├── pkg/
   │   └── types/       # Shared types
   └── go.mod
   ```

2. Initialize Go module with dependencies:
   - `k8s.io/client-go` for Kubernetes API
   - `github.com/gorilla/websocket` for WebSocket
   - `github.com/gin-gonic/gin` for REST API
   - `sigs.k8s.io/controller-runtime` for CRD handling

### 🔧 Task 2: Implement CRD Models
**Status:** Pending
**Dependencies:** Task 1

#### Actions:
1. Create Go structs matching AI Gateway CRDs:
   ```go
   // AIGatewayRoute
   type AIGatewayRouteSpec struct {
       Rules []AIGatewayRouteRule `json:"rules"`
   }
   
   type AIGatewayRouteRule struct {
       Matches     []AIGatewayRouteMatch `json:"matches"`
       BackendRefs []AIGatewayBackendRef `json:"backendRefs"`
   }
   
   // AIServiceBackend  
   type AIServiceBackendSpec struct {
       BackendRef         BackendReference         `json:"backendRef"`
       SecurityPolicyRef  SecurityPolicyReference  `json:"securityPolicyRef,omitempty"`
       TLSPolicyRef      TLSPolicyReference       `json:"tlsPolicyRef,omitempty"`
       Schema            VersionedApiSchema       `json:"schema,omitempty"`
   }
   
   // BackendSecurityPolicy
   type BackendSecurityPolicySpec struct {
       APIKey      *APIKeyAuth      `json:"apiKey,omitempty"`
       OIDC        *OIDCAuth        `json:"oidc,omitempty"`
       AWSCredentials *AWSAuth      `json:"awsCredentials,omitempty"`
   }
   ```

2. Add JSON/YAML tags and validation annotations
3. Create factory functions for each provider template

### 🔌 Task 3: Create Kubernetes Client
**Status:** Pending  
**Dependencies:** Task 2

#### Actions:
1. Set up Kubernetes client configuration:
   ```go
   type K8sClient struct {
       clientset    kubernetes.Interface
       dynamicClient dynamic.Interface
       config       *rest.Config
   }
   ```

2. Configure RBAC permissions for:
   - Reading/Writing AI Gateway CRDs
   - Watching for resource changes
   - Access to secrets for API keys

3. Implement connection handling and retry logic

### 👀 Task 4: Implement CRD Watchers
**Status:** Pending
**Dependencies:** Task 3

#### Actions:
1. Create watchers for each CRD type:
   ```go
   func (k *K8sClient) WatchAIGatewayRoutes(ctx context.Context, eventCh chan<- WatchEvent) error
   func (k *K8sClient) WatchAIServiceBackends(ctx context.Context, eventCh chan<- WatchEvent) error
   func (k *K8sClient) WatchBackendSecurityPolicies(ctx context.Context, eventCh chan<- WatchEvent) error
   ```

2. Handle watch events: ADDED, MODIFIED, DELETED
3. Implement event aggregation and filtering
4. Add error handling and reconnection logic

### 🌐 Task 5: Create API Handlers
**Status:** Pending
**Dependencies:** Task 4

#### Actions:
1. Implement REST API endpoints:
   ```
   GET    /api/v1/providers                    # List all LLM providers
   POST   /api/v1/providers                    # Create new provider
   GET    /api/v1/providers/{name}             # Get specific provider
   PUT    /api/v1/providers/{name}             # Update provider
   DELETE /api/v1/providers/{name}             # Delete provider
   
   GET    /api/v1/routes                       # List AI Gateway routes
   POST   /api/v1/routes                       # Create route
   GET    /api/v1/routes/{name}                # Get specific route
   PUT    /api/v1/routes/{name}                # Update route
   DELETE /api/v1/routes/{name}                # Delete route
   ```

2. Implement provider templates API:
   ```
   GET    /api/v1/templates                    # List available templates
   GET    /api/v1/templates/{provider}         # Get provider template
   ```

3. Add input validation and error handling
4. Implement proper HTTP status codes and responses

### 🔗 Task 6: Implement WebSocket Server
**Status:** Pending
**Dependencies:** Task 5

#### Actions:
1. Create WebSocket endpoint for real-time updates:
   ```go
   func (s *Server) HandleWebSocket(c *gin.Context) {
       conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
       // Handle connection lifecycle
   }
   ```

2. Implement message types:
   ```go
   type WSMessage struct {
       Type string      `json:"type"`
       Data interface{} `json:"data"`
   }
   
   // Message types: "provider_created", "provider_updated", "provider_deleted"
   ```

3. Broadcast CRD changes to connected clients
4. Handle client connection management
5. Implement heartbeat and reconnection logic

### ✅ Task 7: Add Validation Logic
**Status:** Pending
**Dependencies:** Task 6

#### Actions:
1. Implement schema validation for each provider:
   ```go
   func ValidateOpenAIConfig(config *OpenAIProviderConfig) error
   func ValidateAWSConfig(config *AWSProviderConfig) error
   func ValidateGCPConfig(config *GCPProviderConfig) error
   func ValidateAzureConfig(config *AzureProviderConfig) error
   ```

2. Add business logic validation:
   - Unique provider names
   - Valid API endpoints
   - Required authentication fields
   - Schema compatibility

3. Implement dry-run validation mode

### 🧪 Task 8: Create Backend Tests
**Status:** Pending
**Dependencies:** Task 7

#### Actions:
1. Unit tests for all components:
   ```
   backend/tests/
   ├── api/           # API handler tests
   ├── handlers/      # Business logic tests
   ├── k8s/          # Kubernetes client tests
   └── integration/   # End-to-end tests
   ```

2. Mock Kubernetes API responses
3. Test WebSocket functionality
4. Integration tests with real CRDs
5. Load testing for WebSocket connections

## Provider Configuration Examples

### OpenAI Template
```yaml
apiVersion: aigateway.envoyproxy.io/v1alpha1
kind: AIServiceBackend
metadata:
  name: openai-gpt4
spec:
  backendRef:
    name: openai-backend
  securityPolicyRef:
    name: openai-apikey
  schema:
    openAPIV3Schema:
      openapi: 3.0.1
      servers:
        - url: https://api.openai.com/v1
```

### AWS Bedrock Template  
```yaml
apiVersion: aigateway.envoyproxy.io/v1alpha1
kind: AIServiceBackend
metadata:
  name: aws-claude
spec:
  backendRef:
    name: aws-bedrock-backend
  securityPolicyRef:
    name: aws-credentials
  schema:
    openAPIV3Schema:
      openapi: 3.0.1
      servers:
        - url: https://bedrock-runtime.us-east-1.amazonaws.com
```

## Success Criteria
- [ ] All AI Gateway CRDs can be managed via REST API
- [ ] Real-time updates work via WebSocket
- [ ] Provider templates are available for major LLM services
- [ ] Configuration validation prevents invalid deployments
- [ ] Comprehensive test coverage (>80%)
- [ ] Docker image builds successfully
- [ ] Integration with frontend dashboard works

## Next Phase
After completing these tasks, the backend will be ready for:
1. Frontend integration testing
2. Kubernetes deployment
3. Production configuration
4. Monitoring and observability setup