# Docker Deployment Guide

Complete guide for deploying the monorepo application using Docker.

## 📋 Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Make (optional, for convenience commands)

## 🚀 Quick Start

### 1. Setup Environment Variables

```bash
# Copy the production environment template
cp .env.production.example .env.production

# Edit .env.production and update with your values
nano .env.production
```

**Important**: Update these values in `.env.production`:
- `POSTGRES_PASSWORD` - Use a strong password
- `VITE_API_URL` - Your server URL (for production: `http://your-domain.com:3000`)

### 2. Build and Start Services

```bash
# Build Docker images
docker-compose --env-file .env.production build

# Start all services
docker-compose --env-file .env.production up -d

# Or use Make commands
make build
make up
```

### 3. Check Services

```bash
# View logs
docker-compose --env-file .env.production logs -f

# Check container status
docker-compose --env-file .env.production ps

# Or use Make
make logs
make ps
```

### 4. Access Application

- **Client**: http://localhost (port 80)
- **Server API**: http://localhost:3000
- **Database**: localhost:5432

## 🏗️ Architecture

```
┌─────────────────┐
│  React Client   │  Port 80 (Nginx)
│   (Container)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Hono Server    │  Port 3000
│   (Container)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │  Port 5432
│   (Container)   │
└─────────────────┘
```

## 📦 Services

### 1. PostgreSQL Database
- **Image**: postgres:16-alpine
- **Port**: 5432
- **Volume**: postgres_data (persistent)
- **Auto-initialization**: Runs schema.sql and seed.sql on first start

### 2. Hono Server
- **Base**: oven/bun:1.3.0-slim
- **Port**: 3000
- **Health Check**: GET /
- **Dependencies**: PostgreSQL

### 3. React Client
- **Base**: nginx:alpine
- **Port**: 80
- **Build**: Vite production build
- **Health Check**: GET /health

## 🛠️ Make Commands

```bash
make help          # Show all available commands
make dev           # Start development mode
make build         # Build Docker images
make up            # Start all services
make down          # Stop all services
make restart       # Restart all services
make logs          # Show logs from all services
make logs-server   # Show server logs only
make logs-client   # Show client logs only
make logs-db       # Show database logs only
make ps            # Show running containers
make clean         # Remove containers, volumes, and images
make shell-server  # Open shell in server container
make shell-db      # Open PostgreSQL shell
make test          # Run tests
```

## 🔧 Docker Compose Commands

### Start Services
```bash
# Start in detached mode
docker-compose --env-file .env.production up -d

# Start with logs
docker-compose --env-file .env.production up

# Start specific service
docker-compose --env-file .env.production up -d server
```

### Stop Services
```bash
# Stop all services
docker-compose --env-file .env.production down

# Stop and remove volumes
docker-compose --env-file .env.production down -v

# Stop and remove images
docker-compose --env-file .env.production down --rmi all
```

### View Logs
```bash
# All services
docker-compose --env-file .env.production logs -f

# Specific service
docker-compose --env-file .env.production logs -f server

# Last 100 lines
docker-compose --env-file .env.production logs --tail=100 server
```

### Execute Commands
```bash
# Server shell
docker-compose --env-file .env.production exec server sh

# Database shell
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db

# Run migrations
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db -f /docker-entrypoint-initdb.d/01-schema.sql
```

## 🔍 Debugging

### Check Container Health
```bash
docker-compose --env-file .env.production ps
```

### View Container Logs
```bash
# Server logs
docker logs monorepo-server -f

# Client logs
docker logs monorepo-client -f

# Database logs
docker logs monorepo-db -f
```

### Inspect Container
```bash
docker inspect monorepo-server
docker inspect monorepo-client
docker inspect monorepo-db
```

### Test Database Connection
```bash
# From host
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db -c "SELECT version();"

# From server container
docker-compose --env-file .env.production exec server sh -c 'psql "$DATABASE_URL" -c "SELECT version();"'
```

## 🔄 Database Operations

### Backup Database
```bash
docker-compose --env-file .env.production exec postgres pg_dump -U postgres app_db > backup.sql
```

### Restore Database
```bash
docker-compose --env-file .env.production exec -T postgres psql -U postgres app_db < backup.sql
```

### Reset Database
```bash
# Drop and recreate
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Re-run migrations
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db -f /docker-entrypoint-initdb.d/01-schema.sql
docker-compose --env-file .env.production exec postgres psql -U postgres -d app_db -f /docker-entrypoint-initdb.d/02-seed.sql
```

## 📊 Monitoring

### Container Stats
```bash
docker stats monorepo-server monorepo-client monorepo-db
```

### Health Checks
```bash
# Server health
curl http://localhost:3000/

# Client health
curl http://localhost/health

# Database health
docker-compose --env-file .env.production exec postgres pg_isready -U postgres
```

## 🚀 Production Deployment

### 1. Environment Variables
```bash
# Use production values
cp .env.production.example .env.production

# Update these values:
POSTGRES_PASSWORD=strong_random_password
VITE_API_URL=https://api.yourdomain.com
SERVER_PORT=3000
CLIENT_PORT=80
```

### 2. Build for Production
```bash
# Build with no cache for fresh production build
docker-compose --env-file .env.production build --no-cache
```

### 3. Deploy
```bash
# Start services
docker-compose --env-file .env.production up -d

# Check all services are healthy
docker-compose --env-file .env.production ps
```

### 4. SSL/TLS (Optional)
Add reverse proxy like Nginx or Traefik for HTTPS support.

## 🧹 Cleanup

### Remove Everything
```bash
# Stop and remove containers, networks, volumes, and images
make clean

# Or manually
docker-compose --env-file .env.production down -v --rmi all
```

### Remove Only Volumes
```bash
make clean-volumes

# Or manually
docker-compose --env-file .env.production down -v
```

## 📝 Environment Variables

### Required Variables (.env.production)

| Variable | Description | Example |
|----------|-------------|---------|
| `NODE_ENV` | Environment mode | `production` |
| `POSTGRES_USER` | Database user | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `your_secure_password` |
| `POSTGRES_DB` | Database name | `app_db` |
| `POSTGRES_PORT` | Database port | `5432` |
| `SERVER_PORT` | Server port | `3000` |
| `CLIENT_PORT` | Client port | `80` |
| `VITE_API_URL` | API URL for client | `http://localhost:3000` |

## 🔐 Security Best Practices

1. **Never commit** `.env.production` with real credentials
2. Use **strong passwords** for `POSTGRES_PASSWORD`
3. Run containers as **non-root users** (already configured)
4. Keep Docker images **up to date**
5. Use **secrets management** for production (Docker Secrets, Vault)
6. Enable **SSL/TLS** for production deployments
7. Implement **rate limiting** and **CORS** properly
8. Regular **security updates** for base images

## 🐛 Troubleshooting

### Server Won't Start
```bash
# Check logs
docker logs monorepo-server

# Common issues:
# - Database not ready: Wait for health check
# - Port conflict: Change SERVER_PORT in .env.production
```

### Client Build Fails
```bash
# Check build args
docker-compose --env-file .env.production config

# Rebuild without cache
docker-compose --env-file .env.production build --no-cache client
```

### Database Connection Failed
```bash
# Check database is running
docker-compose --env-file .env.production ps postgres

# Test connection
docker-compose --env-file .env.production exec postgres pg_isready

# Check environment variables
docker-compose --env-file .env.production exec server env | grep DATABASE_URL
```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/compose-file/)
- [Bun Documentation](https://bun.sh/docs)
- [PostgreSQL Docker](https://hub.docker.com/_/postgres)
- [Nginx Docker](https://hub.docker.com/_/nginx)
