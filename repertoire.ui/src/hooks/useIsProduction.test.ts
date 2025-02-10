import { renderHook } from '@testing-library/react'
import useIsProduction from './useIsProduction.ts'

describe('use Is Production', () => {
  it('should return true when the environment variable environment is production', () => {
    vi.stubEnv('VITE_ENVIRONMENT', 'production')

    const { result } = renderHook(() => useIsProduction())

    expect(result.current).toBeTruthy()
  })

  it('should return false when the environment variable environment is not production', () => {
    vi.stubEnv('VITE_ENVIRONMENT', 'development')

    const { result } = renderHook(() => useIsProduction())

    expect(result.current).toBeFalsy()
  })
})
