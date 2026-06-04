import { Config, ConfigSchema } from '@doota/pb/doota/portal/v1/portal_pb'
import { portalClient } from './grpc'
import { log } from './logger'
import { create } from '@bufbuild/protobuf'

// this is present on build (i.e. http://api.freightstream.ai)
export const CONFIG_API_URI = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8787'

// this is present on build (i.e. http://app.freightstream.ai)
export const CONFIG_PORTAL_URI = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'

export class ConfigProvider {
  config: Config

  constructor() {
    this.config = create(ConfigSchema, {
      auth0Domain: 'domain.auth0.com',
      auth0ClientId: 'xxxxxxxxxxxxxxxx',
      auth0Scope: 'openid email',
      msoftAuth0CallbackUrl: 'http://msoftcallback',
      googleAuth0CallbackUrl: 'http://googlecallback'
    })
  }

  async bootstrap(): Promise<Config> {
    this.config = await this.buildConfig()

    return this.config
  }

  async fetchFromBackend(): Promise<Config> {
    return portalClient.getConfig({})
  }

  async buildConfig(): Promise<Config> {
    try {
      const backendConfig = await this.fetchFromBackend()

      if (backendConfig === null || backendConfig === undefined) {
        log.warn('No backend configuration found, using defaults')
        return this.config
      }

      log.info('retrieve config', { config: backendConfig })
      return backendConfig
    } catch (error) {
      log.error('Failed to fetch backend config, using defaults', { error })
      // Return default config instead of crashing
      return this.config
    }
  }
}

export const configProvider = new ConfigProvider()
