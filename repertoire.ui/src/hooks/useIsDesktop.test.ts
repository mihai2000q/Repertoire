import { renderHook } from "@testing-library/react"
import useIsDesktop from "./useIsDesktop"

describe('use Is Desktop', () => {
    it('should return true when the environment platform is desktop', () => {
        // Arrange
        import.meta.env.VITE_PLATFORM = 'desktop'

        // Act
        const { result } = renderHook(() => useIsDesktop())

        // Assert
        expect(result.current).toBeTruthy()

        import.meta.env.VITE_PLATFORM = ''
    })

    it('should return false when the environment platform is desktop', () => {
        // Arrange
        import.meta.env.VITE_PLATFORM = 'web'

        // Act
        const { result } = renderHook(() => useIsDesktop())

        // Assert
        expect(result.current).toBeFalsy()

        import.meta.env.VITE_PLATFORM = ''
    })
})