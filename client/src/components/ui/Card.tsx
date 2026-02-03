import type {HTMLAttributes, ReactNode} from 'react'

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'outlined' | 'elevated'
  children: ReactNode
}

export function Card({variant = 'default', className = '', children, ...props}: CardProps) {
  const variants = {
    default: 'bg-white border border-gray-200',
    outlined: 'bg-transparent border-2 border-gray-300',
    elevated: 'bg-white shadow-lg',
  }

  return (
    <div className={`rounded-xl p-6 transition-shadow ${variants[variant]} ${className}`} {...props}>
      {children}
    </div>
  )
}

interface CardHeaderProps {
  className?: string
  children: ReactNode
}

export function CardHeader({className = '', children}: CardHeaderProps) {
  return <div className={`mb-4 ${className}`}>{children}</div>
}

interface CardTitleProps {
  className?: string
  children: ReactNode
}

export function CardTitle({className = '', children}: CardTitleProps) {
  return <h3 className={`text-2xl font-bold text-gray-900 ${className}`}>{children}</h3>
}

interface CardDescriptionProps {
  className?: string
  children: ReactNode
}

export function CardDescription({className = '', children}: CardDescriptionProps) {
  return <p className={`text-sm text-gray-600 mt-1 ${className}`}>{children}</p>
}

interface CardContentProps {
  className?: string
  children: ReactNode
}

export function CardContent({className = '', children}: CardContentProps) {
  return <div className={`text-gray-700 ${className}`}>{children}</div>
}

interface CardFooterProps {
  className?: string
  children: ReactNode
}

export function CardFooter({className = '', children}: CardFooterProps) {
  return <div className={`mt-6 pt-4 border-t border-gray-200 ${className}`}>{children}</div>
}
