import { renderHook } from '@testing-library/react'
import useIsProduction from './useIsProduction.ts'

describe('use Is Production', () => {
  it('should return true when the environment variable environment is production', () => {
    import.meta.env.VITE_ENVIRONMENT = 'production'

    const { result } = renderHook(() => useIsProduction())

    expect(result.current).toBeTruthy()

    import.meta.env.VITE_ENVIRONMENT = ''
  })

  it('should return false when the environment variable environment is not production', () => {
    import.meta.env.VITE_ENVIRONMENT = 'development'

    const { result } = renderHook(() => useIsProduction())

    expect(result.current).toBeFalsy()

    import.meta.env.VITE_ENVIRONMENT = ''
  })
})
