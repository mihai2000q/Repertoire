import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeAll } from 'vitest'
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

  afterEach(() => localStorage.clear())

  it('should return undefined when no value is stored and no defaultValue provided', () => {
    const { result } = renderHook(() => useLocalStorage({ key: 'testKey' }))

    expect(result.current[0]).toBeUndefined()
    expect(localStorage.getItem('testKey')).toBeNull()
  })

  it('should return the default value when no value is stored', () => {
    const { result } = renderHook(() =>
      useLocalStorage({ key: 'testKey', defaultValue: 'default' })
    )

    expect(result.current[0]).toBe('default')
  })

  it('should return the stored value when it exists, ignoring defaultValue', () => {
    localStorage.setItem('testKey', JSON.stringify('storedValue'))

    const { result } = renderHook(() =>
      useLocalStorage({ key: 'testKey', defaultValue: 'default' })
    )

    expect(result.current[0]).toBe('storedValue')
  })

  it('should update localStorage when the value changes', () => {
    const { result } = renderHook(() =>
      useLocalStorage({ key: 'testKey', defaultValue: 'default' })
    )

    act(() => {
      result.current[1]('newValue')
    })

    expect(result.current[0]).toBe('newValue')
    expect(JSON.parse(localStorage.getItem('testKey')!)).toBe('newValue')
  })

  it('should remove item from localStorage when set to undefined', () => {
    localStorage.setItem('testKey', JSON.stringify('initialValue'))

    const { result } = renderHook(() => useLocalStorage<string | undefined>({ key: 'testKey' }))

    act(() => {
      result.current[1](undefined)
    })

    expect(result.current[0]).toBeUndefined()
    expect(localStorage.getItem('testKey')).toBeNull()
  })

  it('should handle complex objects with optional defaultValue', () => {
    const defaultValue = { name: 'John', age: 30 }
    const newValue = { name: 'Jane', age: 25 }

    const { result } = renderHook(() => useLocalStorage({ key: 'testKey', defaultValue }))

    expect(result.current[0]).toStrictEqual(defaultValue)

    act(() => {
      result.current[1](newValue)
    })

    expect(result.current[0]).toStrictEqual(newValue)
    expect(JSON.parse(localStorage.getItem('testKey')!)).toStrictEqual(newValue)
  })

  it('should handle setting value when no defaultValue was provided', () => {
    const { result } = renderHook(() => useLocalStorage<number | undefined>({ key: 'testKey' }))

    act(() => {
      result.current[1](42)
    })

    expect(result.current[0]).toBe(42)
    expect(JSON.parse(localStorage.getItem('testKey')!)).toBe(42)
  })

  it('should accept serialize and deserialize functions', () => {
    const { result } = renderHook(() =>
      useLocalStorage<number | undefined>({
        key: 'testKey',
        serialize: (val) => (val - 10).toString(),
        deserialize: (val) => (parseInt(val) + 10)
      })
    )

    act(() => {
      result.current[1](42)
    })

    expect(result.current[0]).toBe(42)
    expect(JSON.parse(localStorage.getItem('testKey')!)).toBe(32)
  })
})
