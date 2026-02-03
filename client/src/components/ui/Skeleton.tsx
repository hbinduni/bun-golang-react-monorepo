interface SkeletonProps {
  className?: string
}

export function Skeleton({className = ''}: SkeletonProps) {
  return <div className={`animate-pulse bg-gray-200 rounded ${className}`} aria-hidden="true" />
}

interface SkeletonListProps {
  count?: number
}

export function SkeletonList({count = 3}: SkeletonListProps) {
  return (
    <div className="space-y-4" role="status" aria-label="Loading content">
      <span className="sr-only">Loading...</span>
      {Array.from({length: count}).map((_, i) => (
        <div key={i} className="border border-gray-200 rounded-lg p-4 space-y-3">
          <Skeleton className="h-6 w-3/4" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-5/6" />
          <Skeleton className="h-3 w-24" />
        </div>
      ))}
    </div>
  )
}

interface SkeletonTextProps {
  lines?: number
  className?: string
}

export function SkeletonText({lines = 3, className = ''}: SkeletonTextProps) {
  return (
    <div className={`space-y-2 ${className}`} aria-hidden="true">
      {Array.from({length: lines}).map((_, i) => (
        <Skeleton key={i} className={`h-4 ${i === lines - 1 ? 'w-3/4' : 'w-full'}`} />
      ))}
    </div>
  )
}

export function SkeletonCard() {
  return (
    <div className="border border-gray-200 rounded-xl p-6" aria-hidden="true">
      <div className="space-y-4">
        <Skeleton className="h-8 w-1/2" />
        <SkeletonText lines={3} />
        <div className="flex gap-2 pt-2">
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-24" />
        </div>
      </div>
    </div>
  )
}
