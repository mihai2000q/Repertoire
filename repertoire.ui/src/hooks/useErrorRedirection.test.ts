import { reduxRouterRenderHook } from '../test-utils.tsx'
import useErrorRedirection from './useErrorRedirection.ts'
import { RootState } from '../state/store.ts'

describe('useErrorRedirection', () => {
  it('should not navigate when the error path is undefined', () => {
    reduxRouterRenderHook(() => useErrorRedirection(), { global: { errorPath: undefined } })

    expect(window.location.pathname).toBe('/')
  })

  it('should navigate to the error path when it is not undefined', () => {
    // Arrange
    const errorPath = '404'

    // Act
    const [_, store] = reduxRouterRenderHook(() => useErrorRedirection(), {
      global: { errorPath: errorPath }
    })

    // Assert
    expect(window.location.pathname).toBe(`/${errorPath}`)
    expect((store.getState() as RootState).global.errorPath).toBeUndefined()
  })
})
