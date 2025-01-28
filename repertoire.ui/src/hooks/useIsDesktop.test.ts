import { renderHook } from '@testing-library/react'
import useIsDesktop from './useIsDesktop'

describe('use Is Desktop', () => {
  it('should return true when the environment platform is desktop', () => {
    vi.stubEnv('VITE_PLATFORM', 'desktop')

    const { result } = renderHook(() => useIsDesktop())

    expect(result.current).toBeTruthy()
  })

  it('should return false when the environment platform is not desktop', () => {
    vi.stubEnv('VITE_PLATFORM', 'web')

    const { result } = renderHook(() => useIsDesktop())

    expect(result.current).toBeFalsy()
  })
})
