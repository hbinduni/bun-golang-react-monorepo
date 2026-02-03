import {getItems} from '@client/api/items'
import {Button} from '@client/components/ui/Button'
import {Card, CardContent, CardHeader, CardTitle} from '@client/components/ui/Card'
import {SkeletonList} from '@client/components/ui/Skeleton'
import type {Item} from '@client/types'
import {formatDate} from '@client/utils'
import {useCallback, useEffect, useState} from 'react'

export function Home() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>()

  const loadItems = useCallback(() => {
    setLoading(true)
    setError(undefined)

    getItems()
      .then(setItems)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    loadItems()
  }, [loadItems])

  // Loading state
  if (loading) {
    return (
      <div className="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <Card>
          <CardHeader>
            <CardTitle>API Data Example</CardTitle>
          </CardHeader>
          <CardContent>
            <SkeletonList count={3} />
          </CardContent>
        </Card>
      </div>
    )
  }

  // Error state
  if (error) {
    return (
      <div className="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <Card variant="outlined">
          <CardHeader>
            <CardTitle>API Data Example</CardTitle>
          </CardHeader>
          <CardContent>
            <div
              role="alert"
              aria-live="assertive"
              className="bg-red-50 dark:bg-red-900/20 border-l-4 border-red-500 p-4 rounded-md mb-6"
            >
              <div className="flex items-start gap-3">
                <svg
                  className="h-6 w-6 text-red-500 flex-shrink-0 mt-0.5"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z"
                    clipRule="evenodd"
                  />
                </svg>
                <div className="flex-1">
                  <h3 className="text-sm font-semibold text-red-800 dark:text-red-300 mb-1">Error loading data</h3>
                  <p className="text-sm text-red-700 dark:text-red-400">{error}</p>
                </div>
              </div>
            </div>
            <Button variant="primary" onClick={loadItems}>
              Try Again
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  // Success state - data loaded
  return (
    <div className="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <Card>
        <CardHeader>
          <CardTitle>API Data Example</CardTitle>
        </CardHeader>
        <CardContent>
          {items.length === 0 ? (
            <div role="status" className="text-center py-12 sm:py-16">
              <svg
                className="mx-auto h-12 w-12 sm:h-16 sm:w-16 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                />
              </svg>
              <h3 className="mt-4 text-base sm:text-lg font-medium text-gray-900 dark:text-gray-100">No items found</h3>
              <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">Get started by creating a new item.</p>
            </div>
          ) : (
            <ul className="space-y-3 sm:space-y-4">
              {items.map((item) => (
                <li
                  key={item.id}
                  className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 sm:p-5 hover:shadow-md dark:hover:shadow-gray-800/50 transition-all duration-200 hover:border-gray-300 dark:hover:border-gray-600"
                >
                  <h3 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
                    {item.title}
                  </h3>
                  <p className="text-sm sm:text-base text-gray-700 dark:text-gray-300 mb-3 leading-relaxed">
                    {item.description}
                  </p>
                  <time
                    dateTime={item.createdAt.toString()}
                    className="text-xs sm:text-sm text-gray-500 dark:text-gray-400"
                  >
                    {formatDate(item.createdAt)}
                  </time>
                </li>
              ))}
            </ul>
          )}

          <p className="mt-6 sm:mt-8 text-xs sm:text-sm text-gray-600 dark:text-gray-400 text-center leading-relaxed">
            ✨ This data comes from the Go/Fiber backend with TypeScript types!
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
