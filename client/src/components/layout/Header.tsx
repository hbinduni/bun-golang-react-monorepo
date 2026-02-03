import reactLogo from '@client/assets/react.svg'
import viteLogo from '/vite.svg'

export function Header() {
  return (
    <header className="w-full px-4 py-6 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto text-center">
        <div className="flex items-center justify-center gap-4 sm:gap-6 mb-4 sm:mb-6">
          <a href="https://vite.dev" target="_blank" rel="noopener noreferrer" aria-label="Visit Vite website">
            <img src={viteLogo} className="logo" alt="Vite logo" />
          </a>
          <a href="https://react.dev" target="_blank" rel="noopener noreferrer" aria-label="Visit React website">
            <img src={reactLogo} className="logo react" alt="React logo" />
          </a>
        </div>
        <h1 className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-gray-900 dark:text-gray-100">
          Bun + Hono + React Monorepo
        </h1>
      </div>
    </header>
  )
}
