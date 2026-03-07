import {getItems} from '@client/api/items'
import {config} from '@client/config'
import type {Item} from '@client/types'
import {formatDate} from '@client/utils'
import {useCallback, useEffect, useState} from 'react'

const techStack = [
  {name: 'Go', color: 'var(--color-go)', icon: 'Go'},
  {name: 'Fiber', color: 'var(--color-fiber)', icon: 'F'},
  {name: 'React 19', color: 'var(--color-react)', icon: 'R'},
  {name: 'Vite 7', color: 'var(--color-vite)', icon: 'V'},
  {name: 'Tailwind 4', color: 'var(--color-tailwind)', icon: 'T'},
  {name: 'tsgo', color: 'var(--color-tsgo)', icon: 'TS'},
  {name: 'PostgreSQL', color: 'var(--color-postgres)', icon: 'PG'},
  {name: 'Docker', color: 'var(--color-docker)', icon: 'D'},
]

function StatusDot({alive}: {alive: boolean}) {
  return (
    <span className="relative flex size-2">
      {alive && (
        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
      )}
      <span className={`relative inline-flex size-2 rounded-full ${alive ? 'bg-emerald-400' : 'bg-white/20'}`} />
    </span>
  )
}

export function Home() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>()
  const [apiAlive, setApiAlive] = useState(false)

  const loadItems = useCallback(() => {
    setLoading(true)
    setError(undefined)

    getItems()
      .then((data) => {
        setItems(data)
        setApiAlive(true)
      })
      .catch((err) => {
        setError(err.message)
        setApiAlive(false)
      })
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    loadItems()
  }, [loadItems])

  return (
    <div className="w-full">
      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-accent-500/5 via-transparent to-transparent" />
        <div className="relative max-w-6xl mx-auto px-6 pt-20 pb-16 text-center">
          <p
            className="animate-fade-up font-mono text-xs tracking-widest uppercase text-white/40 mb-6"
            style={{animationDelay: '0.05s'}}
          >
            Production-Ready Monorepo Template
          </p>
          <h1
            className="animate-fade-up text-4xl sm:text-5xl md:text-6xl font-display font-bold tracking-tight leading-[1.1] mb-6"
            style={{animationDelay: '0.15s'}}
          >
            <span className="text-gradient">Go + Fiber</span>
            <br />
            <span className="text-white/90">meets React</span>
          </h1>
          <p
            className="animate-fade-up max-w-xl mx-auto text-sm sm:text-base text-white/40 leading-relaxed font-light mb-10"
            style={{animationDelay: '0.25s'}}
          >
            A full-stack monorepo with a Go/Fiber API, React 19 SPA, Tailwind CSS 4, native TypeScript via tsgo, and
            one-command Docker/K8s deployment.
          </p>

          {/* API status pill */}
          <div
            className="animate-fade-up inline-flex items-center gap-2 rounded-full border border-white/10 bg-surface-1/80 backdrop-blur px-4 py-2 text-xs font-mono text-white/50"
            style={{animationDelay: '0.35s'}}
          >
            <StatusDot alive={apiAlive} />
            <span className={apiAlive ? 'text-emerald-400' : 'text-white/30'}>
              {loading ? 'connecting...' : apiAlive ? 'API connected' : 'API offline'}
            </span>
            <span className="text-white/15">|</span>
            <span className="text-white/30 truncate max-w-[200px]">{config.VITE_API_URL}</span>
          </div>
        </div>
      </section>

      {/* Tech Stack */}
      <section className="max-w-6xl mx-auto px-6 pb-16">
        <div className="animate-fade-up grid grid-cols-2 sm:grid-cols-4 gap-3" style={{animationDelay: '0.4s'}}>
          {techStack.map((tech) => (
            <div
              key={tech.name}
              className="tech-badge group flex items-center gap-3 rounded-xl border border-white/5 bg-surface-1/60 px-4 py-3"
            >
              <span
                className="flex size-8 shrink-0 items-center justify-center rounded-lg text-[11px] font-mono font-semibold"
                style={{
                  backgroundColor: `color-mix(in srgb, ${tech.color} 15%, transparent)`,
                  color: tech.color,
                }}
              >
                {tech.icon}
              </span>
              <span className="text-sm text-white/60 group-hover:text-white/80 transition-colors">{tech.name}</span>
            </div>
          ))}
        </div>
      </section>

      {/* API Demo */}
      <section className="max-w-6xl mx-auto px-6 pb-20">
        <div
          className="animate-fade-up rounded-2xl border border-white/5 bg-surface-1/40 overflow-hidden"
          style={{animationDelay: '0.5s'}}
        >
          {/* Section header */}
          <div className="flex items-center justify-between border-b border-white/5 px-6 py-4">
            <div className="flex items-center gap-3">
              <span className="font-mono text-xs text-accent-400">GET</span>
              <span className="font-mono text-xs text-white/40">/api/items</span>
            </div>
            {!loading && error && (
              <button
                type="button"
                onClick={loadItems}
                className="text-xs font-mono text-accent-400 hover:text-accent-300 transition-colors"
              >
                retry
              </button>
            )}
          </div>

          {/* Content */}
          <div className="p-6">
            {loading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="animate-pulse flex gap-4 items-start">
                    <div className="size-2 rounded-full bg-white/10 mt-2 shrink-0" />
                    <div className="flex-1 space-y-2">
                      <div className="h-3 bg-white/5 rounded w-1/3" />
                      <div className="h-2 bg-white/[0.03] rounded w-2/3" />
                    </div>
                  </div>
                ))}
              </div>
            ) : error ? (
              <div className="flex flex-col items-center py-8 text-center">
                <div className="size-10 rounded-full bg-red-500/10 flex items-center justify-center mb-4">
                  <span className="text-red-400 text-lg">!</span>
                </div>
                <p className="text-sm text-white/40 mb-1">Could not reach the API</p>
                <p className="text-xs font-mono text-white/20">{error}</p>
              </div>
            ) : items.length === 0 ? (
              <div className="flex flex-col items-center py-10 text-center">
                <div className="size-12 rounded-xl bg-surface-2 flex items-center justify-center mb-4">
                  <svg
                    className="size-5 text-white/20"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5m8.25 3v6.75m0 0l-3-3m3 3l3-3M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z"
                    />
                  </svg>
                </div>
                <p className="text-sm text-white/50 mb-1">No items yet</p>
                <p className="text-xs text-white/25">Connect a PostgreSQL database and seed data to see items here.</p>
              </div>
            ) : (
              <ul className="divide-y divide-white/5">
                {items.map((item) => (
                  <li key={item.id} className="flex items-start gap-4 py-4 first:pt-0 last:pb-0">
                    <span
                      className={`mt-1.5 size-2 rounded-full shrink-0 ${
                        item.status === 'active'
                          ? 'bg-emerald-400'
                          : item.status === 'completed'
                            ? 'bg-accent-400'
                            : 'bg-white/20'
                      }`}
                    />
                    <div className="flex-1 min-w-0">
                      <h3 className="text-sm font-medium text-white/80">{item.title}</h3>
                      <p className="text-xs text-white/35 mt-0.5 line-clamp-1">{item.description}</p>
                    </div>
                    <time className="text-[10px] font-mono text-white/20 shrink-0 mt-0.5">
                      {formatDate(item.createdAt)}
                    </time>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </section>
    </div>
  )
}
