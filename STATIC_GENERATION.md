# Static Site Generation (SSG) for Vite + React

This document explains how to add static generation capabilities to this Vite + React project, similar to Next.js's SSG/SSR features.

## Current Stack Limitations

Our current setup:
- **Client**: Bun + Vite + React + TypeScript + Tailwind CSS 4 + Biome
- **Server**: Bun + Hono + TypeScript

**Issue**: Vite + React doesn't support static generation for dynamic routes out of the box like Next.js does.

## Available Solutions

### 1. Vike (Recommended) ⭐

**Best for**: Full-featured SSR/SSG with maximum flexibility

- Full SSR/SSG capabilities for Vite
- Works seamlessly with custom servers (Hono)
- Pre-render dynamic routes at build time
- Most Next.js-like developer experience
- Active maintenance and community

**Installation**:
```bash
bun add vike
```

**Resources**:
- Documentation: https://vike.dev
- GitHub: https://github.com/vikejs/vike

**Use Cases**:
- Complex applications requiring both SSR and SSG
- Projects needing custom server integration (already have Hono)
- Dynamic routing with pre-rendering requirements
- SEO-critical pages with dynamic content

---

### 2. vite-ssg

**Best for**: Simple static site generation

- Focused on static site generation only
- Perfect for blogs, documentation, marketing sites
- Less complex setup than Vike
- Good TypeScript support

**Installation**:
```bash
bun add vite-ssg
bun add -d @types/node
```

**Resources**:
- GitHub: https://github.com/antfu/vite-ssg

**Use Cases**:
- Blogs, portfolios, documentation sites
- Marketing landing pages
- Content-heavy sites without complex server logic
- Projects that don't need SSR, only SSG

---

### 3. vite-plugin-prerender

**Best for**: Minimal pre-rendering needs

- Lightweight solution
- Pre-render specific routes only
- Good for hybrid SPA + static pages
- Simple configuration

**Installation**:
```bash
bun add -d vite-plugin-prerender
```

**Resources**:
- npm: https://www.npmjs.com/package/vite-plugin-prerender

**Use Cases**:
- Mostly SPA with a few static pages
- Simple landing pages
- Projects needing minimal pre-rendering
- Quick wins for SEO on specific routes

---

### 4. Custom Solution with Hono

**Best for**: Full control and custom requirements

- Use existing Hono server to render React at build time
- Store generated HTML files
- Integrate with current build pipeline
- Complete customization

**Approach**:
1. Create build script that uses React Server Components or `renderToString`
2. Generate HTML for each route
3. Save static files to `dist/`
4. Serve with Hono or static hosting

**Use Cases**:
- Unique requirements not covered by existing solutions
- Full control over rendering logic
- Integration with custom build tools
- Learning SSG internals

---

## Comparison Matrix

| Feature | Vike | vite-ssg | prerender | Custom |
|---------|------|----------|-----------|--------|
| SSG Support | ✅ | ✅ | ✅ | ✅ |
| SSR Support | ✅ | ❌ | ❌ | ✅ (DIY) |
| Hono Integration | ✅ Native | ⚠️ Manual | ⚠️ Manual | ✅ Native |
| File Routing | ✅ Auto | ⚠️ Manual | ❌ | ❌ |
| Complexity | Medium | Low | Very Low | High |
| Maintenance | Active | Active | Minimal | Self |
| Learning Curve | Moderate | Easy | Easy | Steep |
| Bundle Size | Medium | Small | Tiny | Variable |

## Next.js Feature Parity

Comparison with Next.js features:

| Feature | Next.js | Vike | vite-ssg | prerender |
|---------|---------|------|----------|-----------|
| SSG | ✅ Built-in | ✅ Plugin | ✅ Plugin | ✅ Plugin |
| SSR | ✅ Built-in | ✅ Plugin | ❌ | ❌ |
| ISR | ✅ Built-in | ⚠️ Manual | ❌ | ❌ |
| File Routing | ✅ Automatic | ✅ Configurable | ⚠️ Manual | ❌ |
| API Routes | ✅ Built-in | ➡️ Use Hono | ➡️ Use Hono | ➡️ Use Hono |
| Image Optimization | ✅ Built-in | ⚠️ Plugin/Manual | ⚠️ Manual | ⚠️ Manual |
| Middleware | ✅ Built-in | ✅ Custom | ❌ | ❌ |

## Recommendation

### Choose **Vike** if you need:
- Both SSR and SSG capabilities
- Next.js-like developer experience
- Dynamic routing with pre-rendering
- Full integration with existing Hono server
- SEO-critical application

### Choose **vite-ssg** if you need:
- Simple static generation only
- Blog, docs, or marketing site
- Minimal complexity
- No server-side rendering requirements

### Choose **vite-plugin-prerender** if you need:
- Quick pre-rendering for a few pages
- Mostly SPA with some static pages
- Minimal setup and maintenance

### Build **Custom Solution** if you need:
- Very specific requirements
- Full control over every aspect
- Learning opportunity
- Unique build pipeline integration

## Implementation Example: Vike (Quick Start)

```bash
# Install Vike
bun add vike

# Install React integration
bun add vike-react

# Install server adapter for Hono
bun add @hattip/adapter-hono
```

**Basic Configuration** (`vite.config.ts`):
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import vike from 'vike/plugin'

export default defineConfig({
  plugins: [
    react(),
    vike({
      prerender: true // Enable SSG
    })
  ]
})
```

**Server Integration** (`server/index.ts`):
```typescript
import { serve } from '@hono/node-server'
import { Hono } from 'hono'
import { renderPage } from 'vike/server'

const app = new Hono()

// Your existing API routes
app.get('/api/*', (c) => {
  // ... existing routes
})

// Vike middleware for SSR/SSG
app.use('*', async (c) => {
  const pageContext = await renderPage({ urlOriginal: c.req.url })
  const { httpResponse } = pageContext

  if (!httpResponse) return c.notFound()

  return c.html(httpResponse.body, httpResponse.statusCode)
})

serve(app)
```

## Additional Resources

- [Vite SSR Guide](https://vitejs.dev/guide/ssr.html)
- [React Server Components](https://react.dev/reference/react-dom/server)
- [Hono SSR Integration](https://hono.dev/guides/ssr)

## Notes

- All solutions maintain compatibility with existing Bun, TypeScript, Tailwind CSS 4, and Biome setup
- Hono server can continue serving API routes independently
- Consider SEO, performance, and hosting requirements when choosing
- Test thoroughly in development before production deployment

---

**Last Updated**: November 2024
**Project Stack**: Bun + Vite + React + Hono + TypeScript + Tailwind CSS 4 + Biome
