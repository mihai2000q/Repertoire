import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeAll, vi } from 'vitest'
import useLocalStorage from './useLocalStorage.ts'

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString()
    },
    removeItem: (key: string) => {
      delete store[key]
    },
    clear: () => {
      store = {}
    }
  }
})()

describe('use Local Storage', () => {
  beforeAll(() => {
    localStorage.clear()
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock,
      writable: true
    })
  })

  afterEach(() => {
    localStorage.clear()
    vi.restoreAllMocks()
  })


  it('should return the default value when stored JSON is malformed', () => {
    localStorage.setItem('testKey', '{ not valid json')

    const { result } = renderHook(() =>
      useLocalStorage({ key: 'testKey', defaultValue: 'default' })
    )

    expect(result.current[0]).toBe('default')
  })

  it('should return undefined when stored JSON is malformed and no defaultValue is provided', () => {
    localStorage.setItem('testKey', '{ not valid json')

    const { result } = renderHook(() => useLocalStorage({ key: 'testKey' }))

    expect(result.current[0]).toBeUndefined()
  })

  it('should fall back to defaultValue when deserialize throws', () => {
    localStorage.setItem('testKey', 'anything')

    const { result } = renderHook(() =>
      useLocalStorage<number>({
        key: 'testKey',
        defaultValue: 0,
        deserialize: () => {
          throw new Error('boom')
        }
      })
    )

    expect(result.current[0]).toBe(0)
  })

  it('should not throw when localStorage.getItem throws', () => {
    vi.spyOn(localStorage, 'getItem').mockImplementation(() => {
      throw new DOMException('Access denied', 'SecurityError')
    })

    const { result } = renderHook(() =>
      useLocalStorage({ key: 'testKey', defaultValue: 'default' })
    )

    expect(result.current[0]).toBe('default')
  })

  it('should not throw and still update in-memory state when localStorage.setItem throws', () => {
    vi.spyOn(localStorage, 'setItem').mockImplementation(() => {
      throw new DOMException('QuotaExceededError', 'QuotaExceededError')
    })

    const { result } = renderHook(() => useLocalStorage({ key: 'testKey', defaultValue: 0 }))

    expect(() => {
      act(() => {
        result.current[1](42)
      })
    }).not.toThrow()

    expect(result.current[0]).toBe(42)
  })
})
