# Quick Start Guide

## 🎯 Your New Monorepo Structure

```
monorepo/
├── server/              # Hono backend (http://localhost:3000)
├── client/              # React Vite (http://localhost:5173)
├── packages/
│   └── shared/          # Shared types & utils
└── package.json         # Bun workspace config
```

## 🚀 Running the Project

### Development Mode

Run both server and client together:
```bash
# Terminal 1 - Server
bun dev:server

# Terminal 2 - Client
bun dev:client
```

Or run them separately as needed.

### Type Checking

```bash
bun run typecheck
```

### Building

```bash
# Build everything
bun build

# Build individually
bun build:server
bun build:client
```

## 💡 Using Shared Types

### Server (`server/src/index.ts`)
```typescript
import type { ApiResponse, Feedback } from '@shared/types'

app.get('/api/feedback', (c) => {
  const response: ApiResponse<Feedback[]> = {
    success: true,
    data: [/* ... */]
  }
  return c.json(response)
})
```

### Client (`client/src/api/feedback.ts`)
```typescript
import type { ApiResponse, Feedback } from '@shared/types'

export async function getFeedback(): Promise<Feedback[]> {
  const response = await fetch(`${API_BASE_URL}/api/feedback`)
  const data: ApiResponse<Feedback[]> = await response.json()
  return data.data!
}
```

### Component (`client/src/App.tsx`)
```typescript
import type { Feedback } from '@shared/types'
import { formatDate } from '@shared/utils'

const [feedbacks, setFeedbacks] = useState<Feedback[]>([])
```

## 🔧 Adding Dependencies

```bash
# Server dependency
bun add <package> --filter @monorepo/server

# Client dependency
bun add <package> --filter @monorepo/client

# Shared package dependency
bun add <package> --filter @monorepo/shared

# Dev dependency (root)
bun add -d <package>
```

## 📦 Key Benefits of This Setup

✅ **Type Safety**: Shared types between frontend and backend
✅ **Single Source of Truth**: API contracts defined once
✅ **Fast**: Bun workspaces for instant installs
✅ **Hot Reload**: Both server and client support HMR
✅ **Clean Architecture**: Clear separation of concerns

## 🎨 Path Aliases

### Server
- `@server/*` → `server/src/*` (access server modules from within server)
- `@shared/*` → `packages/shared/src/*`

### Client
- `@/*` → `client/src/*` (access client modules from within client)
- `@shared/*` → `packages/shared/src/*`

## 🧪 Example Workflow

1. Define types in `packages/shared/src/types/`
2. Create API endpoint in `server/src/`
3. Use types in both server response and client request
4. TypeScript ensures everything matches!

## 📝 Next Steps

1. Add more shared types to `packages/shared/src/types/`
2. Create API routes in `server/src/routes/`
3. Build React components in `client/src/components/`
4. Enjoy type-safe full-stack development!
