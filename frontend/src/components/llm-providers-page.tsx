"use client"

import { useState, useEffect } from "react"
import { 
  IconPlus, 
  IconSearch, 
  IconTrash, 
  IconEdit,
  IconRobot,
  IconAlertCircle
} from "@tabler/icons-react"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { LLMProviderService } from "@/services/llm-provider.service"
import type { LLMProvider, CreateLLMProviderRequest } from "@/types/llm-provider"

export function LLMProvidersPage() {
  const [providers, setProviders] = useState<LLMProvider[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState("")
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [editingProvider, setEditingProvider] = useState<LLMProvider | null>(null)
  
  // Form state for creating new provider
  const [newProvider, setNewProvider] = useState<CreateLLMProviderRequest>({
    name: "",
    schema: "",
    authType: "APIKey",
    host: "",
    port: 443
  })

  // Form state for editing provider
  const [editProvider, setEditProvider] = useState<CreateLLMProviderRequest>({
    name: "",
    schema: "",
    authType: "APIKey",
    host: "",
    port: 443
  })

  // Helper function to get default auth type based on schema
  const getDefaultAuthType = (schema: string): string => {
    switch (schema) {
      case "OpenAI":
        return "APIKey"
      case "AWSBedrock":
        return "AWS"
      case "AzureOpenAI":
        return "Azure"
      case "GCPVertexAI":
        return "GCP"
      default:
        return "APIKey"
    }
  }

  // Helper function to get default host based on schema
  const getDefaultHost = (schema: string): string => {
    switch (schema) {
      case "OpenAI":
        return "api.openai.com"
      case "AWSBedrock":
        return "bedrock-runtime.us-east-1.amazonaws.com"
      case "AzureOpenAI":
        return "your-resource.openai.azure.com"
      case "GCPVertexAI":
        return "us-central1-aiplatform.googleapis.com"
      default:
        return ""
    }
  }

  // Handle schema change to auto-populate defaults
  const handleSchemaChange = (value: string) => {
    setNewProvider(prev => ({
      ...prev,
      schema: value,
      authType: getDefaultAuthType(value),
      host: getDefaultHost(value)
    }))
  }

  // Load providers on component mount
  useEffect(() => {
    loadProviders()
  }, [])

  const loadProviders = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await LLMProviderService.getProvidersRaw()
      setProviders(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load providers')
      setProviders([])
    } finally {
      setLoading(false)
    }
  }

  const handleCreateProvider = async () => {
    try {
      const created = await LLMProviderService.createProvider(newProvider)
      // Convert display back to raw for consistency
      const rawProvider = await LLMProviderService.getProviderRaw(created.name)
      setProviders(prev => [...prev, rawProvider])
      setIsCreateDialogOpen(false)
      setNewProvider({
        name: "",
        schema: "",
        authType: "APIKey",
        host: "",
        port: 443
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create provider')
    }
  }

  const handleDeleteProvider = async (name: string) => {
    try {
      await LLMProviderService.deleteProvider(name)
      setProviders(prev => prev.filter(p => p.name !== name))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete provider')
    }
  }

  const handleEditProvider = (provider: LLMProvider) => {
    setEditingProvider(provider)
    
    // Convert provider to edit form format
    const editForm: CreateLLMProviderRequest = {
      name: provider.name,
      schema: provider.schema,
      authType: provider.auth.type,
      host: provider.backend.host,
      port: provider.backend.port,
      tlsHostname: provider.tls?.hostname,
      tlsWellKnownCACertificates: provider.tls?.wellKnownCACertificates,
    }

    // Add auth-specific fields
    if (provider.auth.apiKey) {
      editForm.apiKey = provider.auth.apiKey
    }
    if (provider.auth.aws) {
      editForm.awsRegion = provider.auth.aws.region
      editForm.awsAccessKeyId = provider.auth.aws.accessKeyId
      editForm.awsSecretAccessKey = provider.auth.aws.secretAccessKey
    }
    if (provider.auth.gcp) {
      editForm.gcpProjectId = provider.auth.gcp.projectId
      editForm.gcpLocation = provider.auth.gcp.location
      editForm.gcpWorkloadIdentityPoolName = provider.auth.gcp.workloadIdentityPoolName
      editForm.gcpWorkloadIdentityProviderName = provider.auth.gcp.workloadIdentityProviderName
      editForm.gcpServiceAccountName = provider.auth.gcp.serviceAccountName
      editForm.gcpOidcIssuer = provider.auth.gcp.oidcIssuer
      editForm.gcpOidcClientId = provider.auth.gcp.oidcClientId
      editForm.gcpOidcClientSecret = provider.auth.gcp.oidcClientSecret
      editForm.gcpPrivateKey = provider.auth.gcp.privateKey
      editForm.gcpClientEmail = provider.auth.gcp.clientEmail
      editForm.gcpServiceAccountProjectId = provider.auth.gcp.serviceAccountProjectId
      editForm.gcpClientId = provider.auth.gcp.clientId
      editForm.gcpAuthUri = provider.auth.gcp.authUri
      editForm.gcpTokenUri = provider.auth.gcp.tokenUri
    }
    if (provider.auth.azure) {
      editForm.azureClientId = provider.auth.azure.clientId
      editForm.azureTenantId = provider.auth.azure.tenantId
      editForm.azureApiKey = provider.auth.azure.apiKey
    }

    setEditProvider(editForm)
    setIsEditDialogOpen(true)
  }

  const handleUpdateProvider = async () => {
    if (!editingProvider) return
    
    try {
      // For now, delete and recreate the provider (since there's no PUT endpoint)
      await LLMProviderService.deleteProvider(editingProvider.name)
      const updated = await LLMProviderService.createProvider(editProvider)
      const rawProvider = await LLMProviderService.getProviderRaw(updated.name)
      
      setProviders(prev => prev.filter(p => p.name !== editingProvider.name).concat(rawProvider))
      setIsEditDialogOpen(false)
      setEditingProvider(null)
      setEditProvider({
        name: "",
        schema: "",
        authType: "APIKey",
        host: "",
        port: 443
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update provider')
    }
  }

  const handleEditSchemaChange = (schema: string) => {
    const authType = getDefaultAuthType(schema)
    const host = getDefaultHost(schema)
    
    setEditProvider(prev => ({
      ...prev,
      schema,
      authType,
      host,
      // Clear auth fields when schema changes
      apiKey: "",
      awsRegion: "",
      awsAccessKeyId: "",
      awsSecretAccessKey: "",
      gcpProjectId: "",
      gcpLocation: "",
      gcpWorkloadIdentityPoolName: "",
      gcpWorkloadIdentityProviderName: "",
      gcpServiceAccountName: "",
      gcpOidcIssuer: "",
      gcpOidcClientId: "",
      gcpOidcClientSecret: "",
      gcpPrivateKey: "",
      gcpClientEmail: "",
      gcpServiceAccountProjectId: "",
      gcpClientId: "",
      gcpAuthUri: "",
      gcpTokenUri: "",
      azureClientId: "",
      azureTenantId: "",
      azureApiKey: "",
    }))
  }

  const filteredProviders = providers.filter(provider =>
    provider.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    provider.schema.toLowerCase().includes(searchTerm.toLowerCase()) ||
    provider.backend.host.toLowerCase().includes(searchTerm.toLowerCase())
  )

  if (loading) {
    return (
      <div className="@container/main flex flex-1 flex-col gap-2">
        <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
          <div className="px-4 lg:px-6">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-center h-32">
                  <div className="text-muted-foreground">Loading providers...</div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="@container/main flex flex-1 flex-col gap-2">
      <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        <div className="px-4 lg:px-6">
          {error && (
            <Card className="mb-4 border-red-200 bg-red-50">
              <CardContent className="p-4">
                <div className="flex items-center gap-2 text-red-800">
                  <IconAlertCircle className="w-4 h-4" />
                  <span>{error}</span>
                </div>
              </CardContent>
            </Card>
          )}
          
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <IconRobot className="w-5 h-5" />
                    LLM Providers
                  </CardTitle>
                  <CardDescription>
                    Manage your LLM providers and configurations
                  </CardDescription>
                </div>
                <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                  <DialogTrigger asChild>
                    <Button>
                      <IconPlus className="w-4 h-4 mr-2" />
                      Add Provider
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-md max-h-[80vh] overflow-y-auto">
                    <DialogHeader>
                      <DialogTitle>Create LLM Provider</DialogTitle>
                      <DialogDescription>
                        Add a new LLM provider configuration
                      </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                      <div className="grid gap-2">
                        <Label htmlFor="name">Name</Label>
                        <Input
                          id="name"
                          value={newProvider.name}
                          onChange={(e) => setNewProvider(prev => ({...prev, name: e.target.value}))}
                          placeholder="e.g., openai-gpt4"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="schema">Schema</Label>
                        <Select
                          value={newProvider.schema}
                          onValueChange={handleSchemaChange}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select provider schema" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="OpenAI">OpenAI</SelectItem>
                            <SelectItem value="AzureOpenAI">Azure OpenAI</SelectItem>
                            <SelectItem value="AWSBedrock">AWS Bedrock</SelectItem>
                            <SelectItem value="GCPVertexAI">Google Vertex AI</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="host">Host</Label>
                        <Input
                          id="host"
                          value={newProvider.host}
                          onChange={(e) => setNewProvider(prev => ({...prev, host: e.target.value}))}
                          placeholder="api.openai.com"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="port">Port</Label>
                        <Input
                          id="port"
                          type="number"
                          value={newProvider.port}
                          onChange={(e) => setNewProvider(prev => ({...prev, port: parseInt(e.target.value) || 443}))}
                          placeholder="443"
                        />
                      </div>
                      
                      {/* Conditional Auth Fields based on Schema */}
                      {newProvider.schema === "OpenAI" && (
                        <div className="grid gap-2">
                          <Label htmlFor="apiKey">API Key</Label>
                          <Input
                            id="apiKey"
                            type="password"
                            value={newProvider.apiKey || ""}
                            onChange={(e) => setNewProvider(prev => ({...prev, apiKey: e.target.value}))}
                            placeholder="sk-..."
                            required
                          />
                        </div>
                      )}
                      
                      {newProvider.schema === "AWSBedrock" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="awsRegion">AWS Region</Label>
                            <Input
                              id="awsRegion"
                              value={newProvider.awsRegion || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, awsRegion: e.target.value}))}
                              placeholder="us-east-1"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="awsAccessKeyId">AWS Access Key ID (Optional)</Label>
                            <Input
                              id="awsAccessKeyId"
                              value={newProvider.awsAccessKeyId || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, awsAccessKeyId: e.target.value}))}
                              placeholder="Use IAM roles when possible"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="awsSecretAccessKey">AWS Secret Access Key (Optional)</Label>
                            <Input
                              id="awsSecretAccessKey"
                              type="password"
                              value={newProvider.awsSecretAccessKey || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, awsSecretAccessKey: e.target.value}))}
                              placeholder="Use IAM roles when possible"
                            />
                          </div>
                        </>
                      )}
                      
                      {newProvider.schema === "GCPVertexAI" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpProjectId">GCP Project ID</Label>
                            <Input
                              id="gcpProjectId"
                              value={newProvider.gcpProjectId || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpProjectId: e.target.value}))}
                              placeholder="your-project-id"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpLocation">GCP Location</Label>
                            <Input
                              id="gcpLocation"
                              value={newProvider.gcpLocation || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpLocation: e.target.value}))}
                              placeholder="us-central1"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpWorkloadIdentityPoolName">Workload Identity Pool Name</Label>
                            <Input
                              id="gcpWorkloadIdentityPoolName"
                              value={newProvider.gcpWorkloadIdentityPoolName || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpWorkloadIdentityPoolName: e.target.value}))}
                              placeholder="my-pool"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpWorkloadIdentityProviderName">Workload Identity Provider Name</Label>
                            <Input
                              id="gcpWorkloadIdentityProviderName"
                              value={newProvider.gcpWorkloadIdentityProviderName || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpWorkloadIdentityProviderName: e.target.value}))}
                              placeholder="my-provider"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpServiceAccountName">Service Account Name</Label>
                            <Input
                              id="gcpServiceAccountName"
                              value={newProvider.gcpServiceAccountName || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpServiceAccountName: e.target.value}))}
                              placeholder="my-service-account@project.iam.gserviceaccount.com"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpOidcIssuer">OIDC Issuer</Label>
                            <Input
                              id="gcpOidcIssuer"
                              value={newProvider.gcpOidcIssuer || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpOidcIssuer: e.target.value}))}
                              placeholder="https://token.actions.githubusercontent.com"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpOidcClientId">OIDC Client ID</Label>
                            <Input
                              id="gcpOidcClientId"
                              value={newProvider.gcpOidcClientId || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpOidcClientId: e.target.value}))}
                              placeholder="your-client-id"
                              required
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="gcpOidcClientSecret">OIDC Client Secret</Label>
                            <Input
                              id="gcpOidcClientSecret"
                              type="password"
                              value={newProvider.gcpOidcClientSecret || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, gcpOidcClientSecret: e.target.value}))}
                              placeholder="your-client-secret"
                              required
                            />
                          </div>
                        </>
                      )}
                      
                      {newProvider.schema === "AzureOpenAI" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="azureClientId">Azure Client ID (Optional)</Label>
                            <Input
                              id="azureClientId"
                              value={newProvider.azureClientId || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, azureClientId: e.target.value}))}
                              placeholder="your-client-id"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="azureTenantId">Azure Tenant ID (Optional)</Label>
                            <Input
                              id="azureTenantId"
                              value={newProvider.azureTenantId || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, azureTenantId: e.target.value}))}
                              placeholder="your-tenant-id"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="azureApiKey">Azure API Key</Label>
                            <Input
                              id="azureApiKey"
                              type="password"
                              value={newProvider.azureApiKey || ""}
                              onChange={(e) => setNewProvider(prev => ({...prev, azureApiKey: e.target.value}))}
                              placeholder="Enter Azure OpenAI API key"
                              required
                            />
                          </div>
                        </>
                      )}
                    </div>
                    <DialogFooter>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => {
                          setIsCreateDialogOpen(false)
                          setNewProvider({
                            name: "",
                            schema: "",
                            authType: "APIKey",
                            host: "",
                            port: 443
                          })
                        }}
                      >
                        Cancel
                      </Button>
                      <Button onClick={handleCreateProvider}>
                        Create Provider
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>

                {/* Edit Dialog */}
                <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
                  <DialogContent className="sm:max-w-md max-h-[80vh] overflow-y-auto">
                    <DialogHeader>
                      <DialogTitle>Edit LLM Provider</DialogTitle>
                      <DialogDescription>
                        Update the LLM provider configuration
                      </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                      <div className="grid gap-2">
                        <Label htmlFor="edit-name">Name</Label>
                        <Input
                          id="edit-name"
                          value={editProvider.name}
                          readOnly
                          disabled
                          className="bg-muted text-muted-foreground cursor-not-allowed"
                          placeholder="e.g., openai-gpt4"
                        />
                        <p className="text-xs text-muted-foreground">Name cannot be changed after creation</p>
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="edit-schema">Schema</Label>
                        <Select
                          value={editProvider.schema}
                          onValueChange={handleEditSchemaChange}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select provider schema" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="OpenAI">OpenAI</SelectItem>
                            <SelectItem value="AzureOpenAI">Azure OpenAI</SelectItem>
                            <SelectItem value="AWSBedrock">AWS Bedrock</SelectItem>
                            <SelectItem value="GCPVertexAI">Google Vertex AI</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="edit-host">Host</Label>
                        <Input
                          id="edit-host"
                          value={editProvider.host}
                          onChange={(e) => setEditProvider(prev => ({...prev, host: e.target.value}))}
                          placeholder="api.openai.com"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="edit-port">Port</Label>
                        <Input
                          id="edit-port"
                          type="number"
                          value={editProvider.port}
                          onChange={(e) => setEditProvider(prev => ({...prev, port: parseInt(e.target.value) || 443}))}
                          placeholder="443"
                        />
                      </div>

                      {/* Authentication Fields for Edit */}
                      {editProvider.schema === "OpenAI" && (
                        <div className="grid gap-2">
                          <Label htmlFor="edit-api-key">API Key</Label>
                          <Input
                            id="edit-api-key"
                            type="password"
                            value={editProvider.apiKey || ""}
                            onChange={(e) => setEditProvider(prev => ({...prev, apiKey: e.target.value}))}
                            placeholder="sk-..."
                          />
                        </div>
                      )}

                      {editProvider.schema === "AWSBedrock" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-aws-region">AWS Region</Label>
                            <Input
                              id="edit-aws-region"
                              value={editProvider.awsRegion || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, awsRegion: e.target.value}))}
                              placeholder="us-east-1"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-aws-access-key">Access Key ID</Label>
                            <Input
                              id="edit-aws-access-key"
                              value={editProvider.awsAccessKeyId || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, awsAccessKeyId: e.target.value}))}
                              placeholder="AKIA..."
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-aws-secret">Secret Access Key</Label>
                            <Input
                              id="edit-aws-secret"
                              type="password"
                              value={editProvider.awsSecretAccessKey || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, awsSecretAccessKey: e.target.value}))}
                              placeholder="..."
                            />
                          </div>
                        </>
                      )}

                      {editProvider.schema === "AzureOpenAI" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-azure-client-id">Client ID</Label>
                            <Input
                              id="edit-azure-client-id"
                              value={editProvider.azureClientId || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, azureClientId: e.target.value}))}
                              placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-azure-tenant-id">Tenant ID</Label>
                            <Input
                              id="edit-azure-tenant-id"
                              value={editProvider.azureTenantId || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, azureTenantId: e.target.value}))}
                              placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-azure-api-key">API Key</Label>
                            <Input
                              id="edit-azure-api-key"
                              type="password"
                              value={editProvider.azureApiKey || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, azureApiKey: e.target.value}))}
                              placeholder="..."
                            />
                          </div>
                        </>
                      )}

                      {editProvider.schema === "GCPVertexAI" && (
                        <>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-project-id">Project ID</Label>
                            <Input
                              id="edit-gcp-project-id"
                              value={editProvider.gcpProjectId || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpProjectId: e.target.value}))}
                              placeholder="my-gcp-project"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-location">Location</Label>
                            <Input
                              id="edit-gcp-location"
                              value={editProvider.gcpLocation || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpLocation: e.target.value}))}
                              placeholder="us-central1"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-pool-name">Workload Identity Pool Name</Label>
                            <Input
                              id="edit-gcp-pool-name"
                              value={editProvider.gcpWorkloadIdentityPoolName || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpWorkloadIdentityPoolName: e.target.value}))}
                              placeholder="my-pool"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-provider-name">Workload Identity Provider Name</Label>
                            <Input
                              id="edit-gcp-provider-name"
                              value={editProvider.gcpWorkloadIdentityProviderName || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpWorkloadIdentityProviderName: e.target.value}))}
                              placeholder="my-provider"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-service-account">Service Account Name</Label>
                            <Input
                              id="edit-gcp-service-account"
                              value={editProvider.gcpServiceAccountName || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpServiceAccountName: e.target.value}))}
                              placeholder="my-service-account"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-oidc-issuer">OIDC Issuer</Label>
                            <Input
                              id="edit-gcp-oidc-issuer"
                              value={editProvider.gcpOidcIssuer || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpOidcIssuer: e.target.value}))}
                              placeholder="https://token.actions.githubusercontent.com"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-oidc-client-id">OIDC Client ID</Label>
                            <Input
                              id="edit-gcp-oidc-client-id"
                              value={editProvider.gcpOidcClientId || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpOidcClientId: e.target.value}))}
                              placeholder="repo:owner/repo:environment:production"
                            />
                          </div>
                          <div className="grid gap-2">
                            <Label htmlFor="edit-gcp-oidc-client-secret">OIDC Client Secret</Label>
                            <Input
                              id="edit-gcp-oidc-client-secret"
                              type="password"
                              value={editProvider.gcpOidcClientSecret || ""}
                              onChange={(e) => setEditProvider(prev => ({...prev, gcpOidcClientSecret: e.target.value}))}
                              placeholder="..."
                            />
                          </div>
                        </>
                      )}

                      {/* TLS Configuration for Edit */}
                      <div className="grid gap-2">
                        <Label htmlFor="edit-tls-hostname">TLS Hostname (optional)</Label>
                        <Input
                          id="edit-tls-hostname"
                          value={editProvider.tlsHostname || ""}
                          onChange={(e) => setEditProvider(prev => ({...prev, tlsHostname: e.target.value}))}
                          placeholder="Leave empty to use host"
                        />
                      </div>
                      <div className="grid gap-2">
                        <Label htmlFor="edit-tls-ca">TLS CA Certificates</Label>
                        <Select
                          value={editProvider.tlsWellKnownCACertificates || "System"}
                          onValueChange={(value) => setEditProvider(prev => ({...prev, tlsWellKnownCACertificates: value}))}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select CA certificates" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="System">System</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        variant="outline"
                        onClick={() => {
                          setIsEditDialogOpen(false)
                          setEditingProvider(null)
                          setEditProvider({
                            name: "",
                            schema: "",
                            authType: "APIKey",
                            host: "",
                            port: 443
                          })
                        }}
                      >
                        Cancel
                      </Button>
                      <Button onClick={handleUpdateProvider}>
                        Update Provider
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardHeader>
            <CardContent>
              <div className="flex items-center gap-2 mb-4">
                <div className="relative flex-1 max-w-sm">
                  <IconSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
                  <Input
                    placeholder="Search providers..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>
              
              <div className="border rounded-lg">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Name</TableHead>
                      <TableHead>Host</TableHead>
                      <TableHead>Port</TableHead>
                      <TableHead>Schema</TableHead>
                      <TableHead className="w-[100px]">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredProviders.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                          No providers found
                        </TableCell>
                      </TableRow>
                    ) : (
                      filteredProviders.map((provider) => (
                        <TableRow key={provider.name}>
                          <TableCell className="font-medium">{provider.name}</TableCell>
                          <TableCell>{provider.backend.host}</TableCell>
                          <TableCell>{provider.backend.port}</TableCell>
                          <TableCell>
                            <code className="bg-muted px-2 py-1 rounded text-sm">
                              {provider.schema}
                            </code>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-1">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleEditProvider(provider)}
                              >
                                <IconEdit className="w-4 h-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDeleteProvider(provider.name)}
                                className="text-red-600 hover:text-red-700"
                              >
                                <IconTrash className="w-4 h-4" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
