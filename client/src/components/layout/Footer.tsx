import {config} from '@client/config'

export function Footer() {
  return (
    <footer className="w-full border-t border-white/5 mt-auto">
      <div className="max-w-6xl mx-auto px-6 py-6 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-white/30 font-mono">
        <span>Go + Fiber + React Monorepo Template</span>
        <span className="text-white/20">API: {config.VITE_API_URL}</span>
      </div>
    </footer>
  )
}
