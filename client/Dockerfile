# Client Dockerfile
FROM oven/bun:1.3.0-slim AS base
WORKDIR /app

# Install dependencies
FROM base AS deps
COPY package.json bun.lock* ./
COPY client/package.json ./client/
RUN bun install --frozen-lockfile

# Build client
FROM base AS builder
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/client/node_modules ./client/node_modules
COPY package.json ./
COPY client ./client
WORKDIR /app/client
RUN bun run build

# Production image with nginx
FROM nginx:alpine AS runner

# Copy built files
COPY --from=builder /app/client/dist /usr/share/nginx/html
COPY client/nginx.conf /etc/nginx/conf.d/default.conf

# Copy and set up entrypoint script for runtime config injection
COPY client/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80

# Use entrypoint script to generate runtime config before starting nginx
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
