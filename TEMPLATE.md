# 🚀 Using This as a Project Template

This repository is a full-stack monorepo template. Follow these steps to create your own project from this template.

## Quick Start

### 1. Clone or Fork Repository

```bash
# Option A: Use GitHub's "Use this template" button (recommended)
# Click the "Use this template" button on GitHub

# Option B: Clone and remove git history
git clone https://github.com/YOUR_USERNAME/bun-hono-react-monorepo.git my-project
cd my-project
rm -rf .git
git init
git add .
git commit -m "Initial commit from template"
```

### 2. Customize Project Names

Update these files with your project-specific names:

#### Root Package (`package.json`)
```json
{
  "name": "my-project-monorepo",  // Change this
  "description": "My awesome project"  // Change this
}
```

#### Workspace Packages
Replace `@monorepo` with your organization/project name:

**`server/package.json`:**
```json
{
  "name": "@my-project/server"  // Change this
}
```

**`client/package.json`:**
```json
{
  "name": "@my-project/client"  // Change this
}
```

**`packages/shared/package.json`:**
```json
{
  "name": "@my-project/shared"  // Change this
}
```

#### Update Script References
In root `package.json`, update filter references:
```json
{
  "scripts": {
    "dev:server": "bun run --filter='@my-project/server' dev",
    "dev:client": "bun run --filter='@my-project/client' dev",
    // ... update all @monorepo references
  }
}
```

#### Docker Configuration (`docker-compose.yml`)
```yaml
services:
  postgres:
    container_name: my-project-db  # Change this
  server:
    container_name: my-project-server  # Change this
  client:
    container_name: my-project-client  # Change this

networks:
  my-project-network:  # Change this
```

#### Environment Files
Update `.env.production` and `.env.production.example`:
```env
POSTGRES_DB=my_project_db  # Change this
COMPOSE_PROJECT_NAME=my-project  # Change this
```

### 3. Customize Database Schema

Edit `db/schema.sql` for your data model:
```sql
-- Replace the example users and items tables with your tables
CREATE TABLE IF NOT EXISTS your_table (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  -- your fields here
);
```

Edit `db/seed.sql` with your sample data:
```sql
-- Add your seed data
INSERT INTO your_table (...) VALUES (...);
```

### 4. Update Application Code

#### Server (`server/src/index.ts`)
The template includes a basic Hono server. Customize it:
```typescript
import {Hono} from 'hono'
import {cors} from 'hono/cors'

const app = new Hono()

app.use('/*', cors())

// Add your routes here
app.get('/api/your-endpoint', (c) => {
  return c.json({message: 'Hello from your API'})
})

export default {
  port: 3000,
  fetch: app.fetch,
}
```

#### Client (`client/src/App.tsx`)
Replace the template UI with your React components.

#### Shared Types (`packages/shared/src/types/index.ts`)
Define your application-specific types:
```typescript
export interface YourType {
  id: string
  // your fields
}
```

### 5. Update Documentation

- **README.md**: Replace template info with your project description
- **DOCKER.md**: Update container names in examples
- **QUICK_START.md**: Customize for your project

### 6. Initialize Your Git Repository

```bash
# If you removed .git earlier
git init
git add .
git commit -m "Initial commit: My Project"

# Create GitHub repository and push
git remote add origin https://github.com/YOUR_USERNAME/my-project.git
git branch -M main
git push -u origin main
```

## What's Included

### Tech Stack
- ✅ **Runtime**: Bun 1.3.0
- ✅ **Server**: Hono (TypeScript)
- ✅ **Client**: React 19 + Vite 7 + TypeScript
- ✅ **Database**: PostgreSQL 16
- ✅ **Monorepo**: Turborepo + Bun Workspaces
- ✅ **Linting/Formatting**: Biome.js
- ✅ **Testing**: Bun Test
- ✅ **Docker**: Multi-stage builds + docker-compose

### Path Aliases
Pre-configured TypeScript path aliases for clean imports:

**Server (`server/`):**
- `@server/*` → `server/src/*` - Import server modules
- `@shared/*` → `packages/shared/src/*` - Import shared code

**Client (`client/`):**
- `@/*` → `client/src/*` - Import client modules
- `@shared/*` → `packages/shared/src/*` - Import shared code

**Example Usage:**
```typescript
// Instead of: import { User } from '../../../packages/shared/src/types'
import { User } from '@shared/types'

// Instead of: import { api } from '../../utils/api'
import { api } from '@/utils/api'  // client
import { db } from '@server/utils/db'  // server
```

### Project Structure
```
my-project/
├── server/              # Hono backend
├── client/              # React frontend
├── packages/shared/     # Shared types & utils
├── db/                  # Database schema & seeds
├── tests/               # Test files
├── docker-compose.yml   # Docker orchestration
└── Makefile            # Convenience commands
```

### Available Scripts
```bash
# Development
bun dev              # Run all services
bun dev:server       # Server only
bun dev:client       # Client only

# Building
bun build            # Build all packages
bun typecheck        # Type check
bun lint             # Lint code
bun check            # Check & format
bun test             # Run tests

# Database
bun db:create        # Create schema
bun db:seed          # Seed data
bun db:fresh         # Fresh database

# Docker
make build           # Build images
make up              # Start containers
make logs            # View logs
```

## Customization Checklist

- [ ] Update project name in all `package.json` files
- [ ] Update `@monorepo` namespace to your project name
- [ ] Update Docker container names
- [ ] Update database name in environment files
- [ ] Customize database schema (`db/schema.sql`)
- [ ] Customize seed data (`db/seed.sql`)
- [ ] Update server routes and logic
- [ ] Update client UI components
- [ ] Define shared types for your domain
- [ ] Update README.md with your project info
- [ ] Update environment variable examples
- [ ] Remove or customize sample code
- [ ] Add your project-specific features

## Environment Setup

### Development
```bash
# Server environment
cp server/.env.example server/.env
# Edit server/.env with your values

# Client environment
cp .env.example .env
# Edit .env with your values
```

### Production
```bash
# Docker environment
cp .env.production.example .env.production
# Edit .env.production with secure values
```

## Best Practices

1. **Environment Variables**: Never commit `.env` or `.env.production.local` files
2. **Database Migrations**: Create migration files in `db/migrations/` for schema changes
3. **Type Safety**: Define all types in `packages/shared` for consistency
4. **Testing**: Add tests in `tests/` directory
5. **Docker**: Use multi-stage builds to keep images small
6. **Monorepo**: Keep packages loosely coupled and reusable

## Getting Help

- Check `README.md` for project overview
- Check `DOCKER.md` for Docker deployment guide
- Check `QUICK_START.md` for quick reference

## Contributing to the Template

If you improve this template, consider contributing back:
1. Fork the template repository
2. Make your improvements
3. Submit a pull request

## License

Update the LICENSE file with your preferred license.

---

**Happy Building!** 🚀

If you found this template useful, please give it a ⭐ on GitHub!
