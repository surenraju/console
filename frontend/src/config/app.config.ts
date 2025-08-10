interface AppConfig {
  apiBaseUrl: string;
  apiVersion: string;
  environment: 'development' | 'production' | 'staging';
}

class ConfigService {
  private static instance: ConfigService;
  private config: AppConfig;

  private constructor() {
    this.config = this.loadConfig();
  }

  public static getInstance(): ConfigService {
    if (!ConfigService.instance) {
      ConfigService.instance = new ConfigService();
    }
    return ConfigService.instance;
  }

  private loadConfig(): AppConfig {
    // Get environment variables with fallbacks
    const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
    const apiVersion = import.meta.env.VITE_API_VERSION || 'v1';
    const environment = (import.meta.env.VITE_ENVIRONMENT || 'development') as AppConfig['environment'];

    return {
      apiBaseUrl,
      apiVersion,
      environment,
    };
  }

  public get apiBaseUrl(): string {
    return this.config.apiBaseUrl;
  }

  public get apiVersion(): string {
    return this.config.apiVersion;
  }

  public get environment(): AppConfig['environment'] {
    return this.config.environment;
  }

  public getFullApiUrl(): string {
    return `${this.config.apiBaseUrl}/api/${this.config.apiVersion}`;
  }

  public isDevelopment(): boolean {
    return this.config.environment === 'development';
  }

  public isProduction(): boolean {
    return this.config.environment === 'production';
  }
}

// Export singleton instance
export const config = ConfigService.getInstance();
export default config;
