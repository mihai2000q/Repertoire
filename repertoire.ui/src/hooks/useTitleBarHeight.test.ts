import { renderHook } from '@testing-library/react'
import useTitleBarHeight from './useTitleBarHeight'

describe('use Title Bar Height', () => {
  it('should return 45px when the environment is desktop', () => {
    vi.stubEnv('VITE_PLATFORM', 'desktop')

    const { result } = renderHook(() => useTitleBarHeight())

    expect(result.current).toBe('45px')
  })

  it('should return 0px when the environment is not desktop', () => {
    vi.stubEnv('VITE_PLATFORM', 'web')

    const { result } = renderHook(() => useTitleBarHeight())

    expect(result.current).toBe('0px')
  })
})
