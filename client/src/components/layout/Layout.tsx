import type {ReactNode} from 'react'
import {Footer} from './Footer'
import {Header} from './Header'

interface LayoutProps {
  children: ReactNode
}

export function Layout({children}: LayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      <a href="#main-content" className="sr-only focus:not-sr-only">
        Skip to main content
      </a>
      <Header />
      {/* biome-ignore lint/correctness/useUniqueElementIds: Skip link target - only one main element per page */}
      <main id="main-content" className="flex-1">
        {children}
      </main>
      <Footer />
    </div>
  )
}
