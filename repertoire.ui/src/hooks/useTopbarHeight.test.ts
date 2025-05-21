import useTopbarHeight from './useTopbarHeight.ts'
import { reduxRenderHook } from '../test-utils.tsx'

describe('use Topbar Height', () => {
  it('should return 65 pixels when logged in', () => {
    const [{ result }] = reduxRenderHook(() => useTopbarHeight(), {
      auth: {
        token: 'something',
        historyOnSignIn: undefined
      }
    })
    expect(result.current).toBe('65px')
  })

  it('should return 0 pixels when not logged in', () => {
    const [{ result }] = reduxRenderHook(() => useTopbarHeight())
    expect(result.current).toBe('0px')
  })
})
