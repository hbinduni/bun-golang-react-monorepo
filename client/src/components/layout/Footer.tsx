import {config} from '@client/config'
import packageJson from '../../../package.json'

export function Footer() {
  return (
    <footer className="w-full px-4 py-6 sm:px-6 lg:px-8 mt-auto">
      <div className="max-w-7xl mx-auto text-center">
        <p className="read-the-docs text-sm sm:text-base">Monorepo Template with Type Safety</p>
        <div className="mt-6 sm:mt-8 p-4 sm:p-6 bg-white/5 dark:bg-black/10 rounded-lg text-xs sm:text-sm opacity-70 max-w-md mx-auto">
          <div className="mb-2">
            <strong className="font-semibold">Version:</strong>{' '}
            <span className="text-gray-600 dark:text-gray-400">{packageJson.version}</span>
          </div>
          <div>
            <strong className="font-semibold">API URL:</strong>{' '}
            <span className="text-gray-600 dark:text-gray-400 break-all">{config.VITE_API_URL}</span>
          </div>
        </div>
      </div>
    </footer>
  )
}
