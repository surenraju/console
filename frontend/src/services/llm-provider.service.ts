import type { 
  LLMProvider, 
  CreateLLMProviderRequest, 
  LLMProviderDisplay
} from '@/types/llm-provider';
import { createLLMProviderFromForm, toLLMProviderDisplay } from '@/types/llm-provider';
import { ApiService } from './api.service';

export class LLMProviderService {
  private static readonly BASE_ENDPOINT = '/llm/providers';

  static async getProviders(): Promise<LLMProviderDisplay[]> {
    const providers = await ApiService.get<LLMProvider[]>(`${this.BASE_ENDPOINT}`);
    return providers.map(toLLMProviderDisplay);
  }

  static async getProvidersRaw(): Promise<LLMProvider[]> {
    return ApiService.get<LLMProvider[]>(`${this.BASE_ENDPOINT}`);
  }

  static async getProviderByName(name: string): Promise<LLMProviderDisplay> {
    const provider = await ApiService.get<LLMProvider>(`${this.BASE_ENDPOINT}/${encodeURIComponent(name)}`);
    return toLLMProviderDisplay(provider);
  }

  static async getProviderRaw(name: string): Promise<LLMProvider> {
    return ApiService.get<LLMProvider>(`${this.BASE_ENDPOINT}/${encodeURIComponent(name)}`);
  }

  static async createProvider(providerForm: CreateLLMProviderRequest): Promise<LLMProviderDisplay> {
    const provider = createLLMProviderFromForm(providerForm);
    const created = await ApiService.post<LLMProvider>(this.BASE_ENDPOINT, provider);
    return toLLMProviderDisplay(created);
  }

  static async deleteProvider(name: string): Promise<void> {
    return ApiService.delete<void>(`${this.BASE_ENDPOINT}/${encodeURIComponent(name)}`);
  }

  static async updateProvider(name: string, providerForm: CreateLLMProviderRequest): Promise<LLMProviderDisplay> {
    const provider = createLLMProviderFromForm(providerForm);
    const updated = await ApiService.put<LLMProvider>(`${this.BASE_ENDPOINT}/${encodeURIComponent(name)}`, provider);
    return toLLMProviderDisplay(updated);
  }
}
