package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	aigatewayv1alpha1 "github.com/envoyproxy/ai-gateway/api/v1alpha1"
	"github.com/envoyproxy/ai-gateway/console/backend/pkg/llm"
	gatewayv1alpha1 "github.com/envoyproxy/gateway/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	gwapiv1a3 "sigs.k8s.io/gateway-api/apis/v1alpha3"
)

func TestLLMProviderToEnvoyGatewayResources(t *testing.T) {
	testCases := []struct {
		Name       string
		Input      string
		GoldenFile string
	}{
		{Name: "OpenAI", Input: "testdata/llm_provider/openai.json", GoldenFile: "testdata/gateway/openai.json"},
		{Name: "AWS", Input: "testdata/llm_provider/aws.json", GoldenFile: "testdata/gateway/aws.json"},
		{Name: "Azure", Input: "testdata/llm_provider/azure_openai.json", GoldenFile: "testdata/gateway/azure_openai.json"},
		{Name: "GCP", Input: "testdata/llm_provider/gcp_vertex.json", GoldenFile: "testdata/gateway/gcp_vertex.json"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			jsonPath, err := filepath.Abs(tc.Input)
			if err != nil {
				t.Fatalf("failed to get abs path: %v", err)
			}
			jsonData, err := os.ReadFile(jsonPath)
			if err != nil {
				t.Fatalf("failed to read json: %v", err)
			}
			var provider llm.LLMProvider
			if err := json.Unmarshal(jsonData, &provider); err != nil {
				t.Fatalf("failed to unmarshal json: %v", err)
			}
			resources, err := provider.ToEnvoyGatewayResources()

			// Verify consistent resource naming for all resources
			// This should apply to all providers - each resource should match the provider name
			for _, res := range resources {
				if meta, ok := res.(interface {
					GetName() string
					GetNamespace() string
				}); ok {
					if meta.GetName() != provider.Name {
						t.Fatalf("Resource name %s does not match provider name %s",
							meta.GetName(), provider.Name)
					}
					if meta.GetNamespace() != provider.Namespace {
						t.Fatalf("Resource namespace %s does not match provider namespace %s",
							meta.GetNamespace(), provider.Namespace)
					}
				}
			}

			if err != nil {
				// We shouldn't get errors for any provider now that GCP is fully implemented
				t.Fatalf("ToEnvoyGatewayResources failed: %v", err)
			}
			// Marshal resources to JSON array
			marshaled, err := json.MarshalIndent(resources, "", "  ")
			if err != nil {
				t.Fatalf("failed to marshal resources to json: %v", err)
			}
			// Load golden JSON
			goldenPath, err := filepath.Abs(tc.GoldenFile)
			if err != nil {
				t.Fatalf("failed to get abs path: %v", err)
			}
			goldenData, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden json: %v", err)
			}
			// Normalize JSON by removing null fields
			var actualObj interface{}
			if err := json.Unmarshal(marshaled, &actualObj); err != nil {
				t.Fatalf("failed to unmarshal actual json: %v", err)
			}
			var expectedObj interface{}
			if err := json.Unmarshal(goldenData, &expectedObj); err != nil {
				t.Fatalf("failed to unmarshal golden json: %v", err)
			}
			var removeNulls func(interface{}) interface{}
			removeNulls = func(v interface{}) interface{} {
				switch vv := v.(type) {
				case map[string]interface{}:
					m := map[string]interface{}{}
					for k, val := range vv {
						cleaned := removeNulls(val)
						if cleaned != nil {
							m[k] = cleaned
						}
					}
					if len(m) == 0 {
						return nil
					}
					return m
				case []interface{}:
					arr := make([]interface{}, 0, len(vv))
					for _, item := range vv {
						cleaned := removeNulls(item)
						if cleaned != nil {
							arr = append(arr, cleaned)
						}
					}
					return arr
				default:
					if vv == nil {
						return nil
					}
					return vv
				}
			}
			actualClean := removeNulls(actualObj)
			expectedClean := removeNulls(expectedObj)
			actualNorm, _ := json.MarshalIndent(actualClean, "", "  ")
			expectedNorm, _ := json.MarshalIndent(expectedClean, "", "  ")
			if string(actualNorm) != string(expectedNorm) {
				t.Errorf("JSON output does not match golden file for %s\nExpected:\n%s\nActual:\n%s", tc.Name, string(expectedNorm), string(actualNorm))
			}
		})
	}
}

// TestEnvoyGatewayResourcesToLLMProvider tests the conversion from Envoy Gateway resources back to LLMProvider
func TestEnvoyGatewayResourcesToLLMProvider(t *testing.T) {
	testCases := []struct {
		Name       string
		GoldenFile string // Expected LLMProvider JSON
		InputFile  string // Gateway resources JSON
	}{
		{Name: "OpenAI", InputFile: "testdata/gateway/openai.json", GoldenFile: "testdata/llm_provider/openai.json"},
		{Name: "AWS", InputFile: "testdata/gateway/aws.json", GoldenFile: "testdata/llm_provider/aws.json"},
		{Name: "Azure", InputFile: "testdata/gateway/azure_openai.json", GoldenFile: "testdata/llm_provider/azure_openai.json"},
		{Name: "GCP", InputFile: "testdata/gateway/gcp_vertex.json", GoldenFile: "testdata/llm_provider/gcp_vertex.json"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Load gateway resources JSON
			gatewayPath, err := filepath.Abs(tc.InputFile)
			if err != nil {
				t.Fatalf("failed to get abs path: %v", err)
			}
			gatewayData, err := os.ReadFile(gatewayPath)
			if err != nil {
				t.Fatalf("failed to read gateway json: %v", err)
			}

			// Load expected LLMProvider JSON
			expectedPath, err := filepath.Abs(tc.GoldenFile)
			if err != nil {
				t.Fatalf("failed to get abs path: %v", err)
			}
			expectedData, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("failed to read expected json: %v", err)
			}

			var expectedProvider llm.LLMProvider
			if err := json.Unmarshal(expectedData, &expectedProvider); err != nil {
				t.Fatalf("failed to unmarshal expected provider: %v", err)
			}

			// Parse gateway resources into typed objects
			var resources []interface{}
			var gatewayResources []json.RawMessage
			if err := json.Unmarshal(gatewayData, &gatewayResources); err != nil {
				t.Fatalf("failed to unmarshal gateway resources: %v", err)
			}

			for _, res := range gatewayResources {
				var typeMeta struct {
					Kind       string `json:"kind"`
					APIVersion string `json:"apiVersion"`
				}
				if err := json.Unmarshal(res, &typeMeta); err != nil {
					t.Fatalf("failed to unmarshal typeMeta: %v", err)
				}

				switch typeMeta.Kind {
				case "Backend":
					var backend gatewayv1alpha1.Backend
					if err := json.Unmarshal(res, &backend); err != nil {
						t.Fatalf("failed to unmarshal Backend: %v", err)
					}
					resources = append(resources, &backend)
				case "BackendTLSPolicy":
					var tlsPolicy gwapiv1a3.BackendTLSPolicy
					if err := json.Unmarshal(res, &tlsPolicy); err != nil {
						t.Fatalf("failed to unmarshal BackendTLSPolicy: %v", err)
					}
					resources = append(resources, &tlsPolicy)
				case "Secret":
					var secret corev1.Secret
					if err := json.Unmarshal(res, &secret); err != nil {
						t.Fatalf("failed to unmarshal Secret: %v", err)
					}
					resources = append(resources, &secret)
				case "BackendSecurityPolicy":
					var bsp aigatewayv1alpha1.BackendSecurityPolicy
					if err := json.Unmarshal(res, &bsp); err != nil {
						t.Fatalf("failed to unmarshal BackendSecurityPolicy: %v", err)
					}
					resources = append(resources, &bsp)
				case "AIServiceBackend":
					var aisb aigatewayv1alpha1.AIServiceBackend
					if err := json.Unmarshal(res, &aisb); err != nil {
						t.Fatalf("failed to unmarshal AIServiceBackend: %v", err)
					}
					resources = append(resources, &aisb)
				default:
					t.Fatalf("unknown resource kind: %s", typeMeta.Kind)
				}
			}

			// Convert gateway resources to LLMProvider
			actualProvider, err := llm.ToLLMProvider(resources)
			if err != nil {
				t.Fatalf("ToLLMProvider failed: %v", err)
			}

			// Compare providers
			actualJSON, _ := json.MarshalIndent(actualProvider, "", "  ")
			expectedJSON, _ := json.MarshalIndent(expectedProvider, "", "  ")

			if string(actualJSON) != string(expectedJSON) {
				t.Errorf("Converted LLMProvider does not match expected for %s\nExpected:\n%s\nActual:\n%s",
					tc.Name, string(expectedJSON), string(actualJSON))
			}
		})
	}
}

type LLMProviderTestCase struct {
	Name     string
	YamlFile string
	JsonFile string
}
