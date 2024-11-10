import { renderHook } from '@testing-library/react'
import useTitleBarHeight from './useTitleBarHeight'

describe('use Title Bar Height', () => {
  it('should return 45px when the environment is desktop', () => {
    // Arrange
    import.meta.env.VITE_PLATFORM = 'desktop'

    // Act
    const { result } = renderHook(() => useTitleBarHeight())

    // Assert
    expect(result.current).toBe('45px')

    import.meta.env.VITE_PLATFORM = ''
  })

  it('should return 0px when the environment is not desktop', () => {
    // Arrange
    import.meta.env.VITE_PLATFORM = 'web'

    // Act
    const { result } = renderHook(() => useTitleBarHeight())

    // Assert
    expect(result.current).toBe('0px')

    import.meta.env.VITE_PLATFORM = ''
  })
})
