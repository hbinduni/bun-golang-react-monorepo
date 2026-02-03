// Client-side utility functions

export function formatDate(date: Date | string): string {
  const dateObj = typeof date === 'string' ? new Date(date) : date

  if (Number.isNaN(dateObj.getTime())) {
    return 'Invalid Date'
  }

  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(dateObj)
}

export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}
