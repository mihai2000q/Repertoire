import { renderHook } from '@testing-library/react'
import useTopbarHeight from './useTopbarHeight.ts'

describe('use Topbar Height', () => {
  it('should return 65 pixels', () => {
    const { result } = renderHook(() => useTopbarHeight())
    expect(result.current).toBe('65px')
  })
})
