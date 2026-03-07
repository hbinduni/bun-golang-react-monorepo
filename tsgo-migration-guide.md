# Migrating from tsc to tsgo (@typescript/native-preview)

> Instruction document for Claude Code. Paste this into your conversation when migrating a TypeScript project from `tsc` to `tsgo` (TypeScript 7 native compiler).

## What is tsgo?

`tsgo` is the Go-based native TypeScript compiler (TypeScript 7 / Project Corsa). It provides ~10x faster type-checking through native parallelization. Install via `@typescript/native-preview` on npm.

## Migration Strategy

**Use tsgo for type-checking. Keep tsc for emit.** tsgo's `.d.ts` and `.js` generation is still incomplete. The safe path is:

| Use Case | Tool | Why |
|----------|------|-----|
| Type-checking (`--noEmit`) | **tsgo** | ~10x faster, production-ready |
| Build mode (`-b`, no emit) | **tsgo** | Supported, parallel project refs |
| `.d.ts` + `.js` emit | **tsc** | tsgo emit still maturing |
| Watch mode (dev) | **tsc** | tsgo watch less battle-tested |

## Step-by-Step Migration

### 1. Install

```bash
bun add -d @typescript/native-preview
# or: npm install -D @typescript/native-preview
```

This provides the `tsgo` binary.

### 2. Audit tsconfig for Breaking Changes

tsgo removes/changes several tsconfig options. Check each tsconfig in the project:

#### Removed Options

| Option | Fix |
|--------|-----|
| `baseUrl` | Remove it. Change paths values to use `./` prefix (e.g., `"src/*"` → `"./src/*"`). Or run: `npx @andrewbranch/ts5to6 --fixBaseUrl tsconfig.json` |
| `moduleResolution: "node"` / `"node10"` | Change to `"bundler"` (Vite/esbuild) or `"nodenext"` (Node.js) |
| `target: "es5"` | Minimum is `"es2015"`, downlevel emit only works for `es2021+` |
| `experimentalDecorators` | Removed. Check if code actually uses decorators. If not, delete the flag |
| `emitDecoratorMetadata` | Removed. Same — delete if unused |

#### Changed Defaults

| Option | Old Default | New Default |
|--------|------------|-------------|
| `strict` | `false` | `true` |
| `types` | auto-enumerate `@types/*` | `[]` (must be explicit) |
| `rootDir` | inferred from sources | directory of `tsconfig.json` |

### 3. Handle the `baseUrl` + Runtime Conflict

**This is the biggest gotcha.** Runtimes like Bun read `tsconfig.json` paths at runtime and need `baseUrl` for non-relative path mappings. tsgo rejects `baseUrl`.

**If your runtime reads tsconfig paths (Bun, ts-node, tsx):**

Create a separate tsconfig for type-checking:

```
tsconfig.json              ← Runtime uses this (keeps baseUrl)
tsconfig.typecheck.json    ← tsgo uses this (no baseUrl, ./ paths)
```

`tsconfig.typecheck.json`:
```jsonc
{
  "extends": "./tsconfig.json",  // or "../tsconfig.json" for monorepo
  "compilerOptions": {
    // Override: remove baseUrl, use ./ prefixed paths
    "noEmit": true,
    "paths": {
      "@/*": ["./src/*"],
      "@lib/*": ["./src/lib/*"]
      // ... mirror all paths from tsconfig.json with ./ prefix
    }
  }
}
```

> **Note:** You cannot "unset" `baseUrl` via extends — the child config must NOT extend a tsconfig that has `baseUrl`. Either extend a root tsconfig without `baseUrl`, or make the typecheck config standalone.

**If your runtime does NOT read tsconfig paths (Vite, Next.js, esbuild):**

Just remove `baseUrl` and add `./` prefix to paths. No separate config needed.

### 4. Update Scripts

#### Type-checking scripts (swap to tsgo)

```jsonc
// Before
"typecheck": "tsc --noEmit -p tsconfig.json"

// After (no baseUrl conflict)
"typecheck": "bunx tsgo --noEmit -p tsconfig.json"

// After (with baseUrl conflict — use separate config)
"typecheck": "bunx tsgo --noEmit -p tsconfig.typecheck.json"
```

#### Build scripts

```jsonc
// Client (Vite) — tsc -b is just type-checking here, swap to tsgo
"build": "bunx tsgo -b && vite build"

// Shared lib (needs .d.ts emit) — keep tsc
"build": "tsc"

// Server (Bun runtime) — optional: add tsgo gate before bun build
"build": "bunx tsgo --noEmit -p tsconfig.typecheck.json && bun build src/index.ts --outdir dist"
```

#### Add safe fallbacks

```jsonc
"typecheck:safe": "tsc --noEmit -p tsconfig.json"
```

### 5. Verify

Run these checks after migration:

```bash
# 1. tsgo type-checking passes
bunx tsgo --noEmit -p tsconfig.json  # (or tsconfig.typecheck.json)

# 2. tsc fallback still works
bun x tsc --noEmit -p tsconfig.json

# 3. Runtime still starts (if runtime reads tsconfig paths)
bun run src/index.ts

# 4. Build still works
bun run build

# 5. Compare speed
time bunx tsgo --noEmit -p tsconfig.json
time bun x tsc --noEmit -p tsconfig.json
```

### 6. Monorepo-Specific Notes

- Build shared/library packages with `tsc` (they need `.d.ts` emit)
- Type-check everything with `tsgo` (including shared, via `--noEmit`)
- tsgo runs project references (`-b`) in **parallel** — big win for monorepos
- Each package with `baseUrl` in its tsconfig may need its own `tsconfig.typecheck.json`

## What NOT to Migrate

| Keep on tsc | Reason |
|-------------|--------|
| Packages that emit `.d.ts` | tsgo declaration emit incomplete |
| `tsc --watch` in dev scripts | tsgo watch mode less stable |
| Any tooling using TypeScript's JS API | tsgo uses new "Corsa" API, old "Strada" API not available |
| CI that needs 100% reliability | Keep `:safe` fallback scripts |

## Dockerfile Considerations

- If Docker build runs `bun install` and your `package.json` has `@typescript/native-preview` in devDependencies, `tsgo` binary will be available in the build stage
- No changes needed to Dockerfile commands themselves — the package.json scripts handle the tool selection
- For multi-stage builds: tsgo is only needed in the builder stage, not the runtime stage

## Troubleshooting

| Error | Cause | Fix |
|-------|-------|-----|
| `TS5102: Option 'baseUrl' has been removed` | tsgo rejects `baseUrl` | Remove `baseUrl`, use `./` paths. Or use separate `tsconfig.typecheck.json` |
| `TS5090: Non-relative paths are not allowed` | Paths values need `./` prefix | Change `"src/*"` to `"./src/*"` |
| `Cannot find module '@lib/...'` at runtime | Removed `baseUrl` broke Bun/ts-node path resolution | Restore `baseUrl` in main tsconfig, use separate typecheck config for tsgo |
| `TS5102: Option 'moduleResolution' value 'node' removed` | `"node"` / `"node10"` removed | Use `"bundler"` or `"nodenext"` |

## Expected Performance Gains

| Project Size | tsc | tsgo | Speedup |
|-------------|-----|------|---------|
| Small (~50 files) | 3-5s | 0.3-0.5s | ~10x |
| Medium (~200 files) | 7-15s | 0.8-1.5s | ~9x |
| Large (~1000+ files) | 30-60s | 3-6s | ~10x |

tsgo uses all CPU cores (you'll see 400-1000%+ CPU utilization).

## References

- [Progress on TypeScript 7 - December 2025](https://devblogs.microsoft.com/typescript/progress-on-typescript-7-december-2025/)
- [Announcing TypeScript Native Previews](https://devblogs.microsoft.com/typescript/announcing-typescript-native-previews/)
- [microsoft/typescript-go](https://github.com/microsoft/typescript-go)
- [ts5to6 migration tool](https://github.com/andrewbranch/ts5to6)
- [TS6 Deprecation Candidates](https://github.com/microsoft/TypeScript/issues/54500)
- [@typescript/native-preview on npm](https://www.npmjs.com/package/@typescript/native-preview)
