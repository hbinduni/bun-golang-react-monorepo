// Development placeholder - runtime config is injected in production by docker-entrypoint.sh
// This file prevents 404 errors in development while maintaining production runtime config capability
//
// In production, this file is replaced with actual runtime configuration:
// window.__RUNTIME_CONFIG__ = {
//   VITE_API_URL: '<injected-at-runtime>'
// }
//
// In development, the app uses build-time environment variables from .env file

console.log('Development mode: Using build-time configuration from .env')
