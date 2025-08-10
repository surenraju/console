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
  
  // Form state for creating new provider
  const [newProvider, setNewProvider] = useState<CreateLLMProviderRequest>({
    name: "",
    schema: "",
    authType: "apiKey",
    host: "",
    port: 443
  })

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
        authType: "apiKey",
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
                  <DialogContent className="sm:max-w-md">
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
                          onValueChange={(value) => setNewProvider(prev => ({...prev, schema: value}))}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select provider schema" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="OpenAI">OpenAI</SelectItem>
                            <SelectItem value="AzureOpenAI">Azure OpenAI</SelectItem>
                            <SelectItem value="AWS">AWS Bedrock</SelectItem>
                            <SelectItem value="GCP">Google Vertex AI</SelectItem>
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
                      <div className="grid gap-2">
                        <Label htmlFor="apiKey">API Key (Optional)</Label>
                        <Input
                          id="apiKey"
                          type="password"
                          value={newProvider.apiKey}
                          onChange={(e) => setNewProvider(prev => ({...prev, apiKey: e.target.value}))}
                          placeholder="sk-..."
                        />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setIsCreateDialogOpen(false)}
                      >
                        Cancel
                      </Button>
                      <Button onClick={handleCreateProvider}>
                        Create Provider
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
                                onClick={() => {/* TODO: Edit functionality */}}
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
