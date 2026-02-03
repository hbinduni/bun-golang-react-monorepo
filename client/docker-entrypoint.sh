#!/bin/sh
# Docker entrypoint script for runtime configuration injection
# This script generates a config.js file with runtime environment variables

set -e

# Generate runtime configuration file
cat > /usr/share/nginx/html/config.js <<EOF
// Runtime configuration - generated at container startup
window.__RUNTIME_CONFIG__ = {
  VITE_API_URL: "${VITE_API_URL:-http://localhost:3000}"
};
EOF

echo "Runtime configuration generated:"
cat /usr/share/nginx/html/config.js

# Start nginx
exec "$@"
