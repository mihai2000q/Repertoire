import { reduxRenderHook } from '../test-utils.tsx'
import useAuth from './useAuth.ts'

describe('use Auth', () => {
  it('should return true when there is a token, so the user is authenticated', () => {
    const [{ result }] = reduxRenderHook(() => useAuth(), { auth: { token: '' } })

    expect(result.current).toBeTruthy()
  })

  it('should return false when there is no token, so the user is not authenticated', () => {
    const [{ result }] = reduxRenderHook(() => useAuth(), { auth: { token: null } })

    expect(result.current).toBeFalsy()
  })
})
