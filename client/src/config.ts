// Runtime configuration loader
// Supports both runtime config (production) and build-time env (development)

interface RuntimeConfig {
  VITE_API_URL: string
}

// Extend Window interface to include runtime config
declare global {
  interface Window {
    __RUNTIME_CONFIG__?: RuntimeConfig
  }
}

/**
 * Get configuration value with fallback priority:
 * 1. Runtime config (window.__RUNTIME_CONFIG__) - for production containers
 * 2. Build-time env (import.meta.env) - for local development
 * 3. Default value
 */
function getConfig(): RuntimeConfig {
  // Production: Runtime config injected by docker-entrypoint.sh
  if (window.__RUNTIME_CONFIG__) {
    return window.__RUNTIME_CONFIG__
  }

  // Development: Build-time environment variables
  return {
    VITE_API_URL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
  }
}

export const config = getConfig()
