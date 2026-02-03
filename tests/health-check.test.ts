import {describe, expect, test} from 'bun:test'

describe('Health Check', () => {
  test('should pass basic health check', () => {
    expect(true).toBe(true)
  })

  test('environment is set up correctly', () => {
    expect(process.env.NODE_ENV).toBeDefined()
  })

  test('can perform basic arithmetic', () => {
    const sum = 1 + 1
    expect(sum).toBe(2)
  })
})
