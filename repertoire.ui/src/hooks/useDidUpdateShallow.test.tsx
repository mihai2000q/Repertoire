import { renderHook } from '@testing-library/react'
import { vi } from 'vitest'
import useDidUpdateShallow from './useDidUpdateShallow' // adjust the import path as needed

describe('use Did Update Shallow', () => {
  it('should not call the effect on initial render', () => {
    const effect = vi.fn()
    renderHook(() => useDidUpdateShallow(effect, []))

    expect(effect).not.toHaveBeenCalled()
  })

  it('should call the effect on subsequent renders', () => {
    const effect = vi.fn()
    const { rerender } = renderHook(
      ({ deps }) => useDidUpdateShallow(effect, deps),
      { initialProps: { deps: [1] } }
    )

    expect(effect).not.toHaveBeenCalled()

    rerender({ deps: [2] })
    expect(effect).toHaveBeenCalledOnce()

    rerender({ deps: [3] })
    expect(effect).toHaveBeenCalledTimes(2)
  })

  it('should not call the effect if shallow comparison of dependencies is equal', () => {
    const effect = vi.fn()
    const { rerender } = renderHook(
      ({ deps }) => useDidUpdateShallow(effect, deps),
      { initialProps: { deps: [{ a: 1 }] } }
    )

    expect(effect).not.toHaveBeenCalled()

    // Same object reference
    rerender({ deps: [{ a: 1 }] })
    expect(effect).not.toHaveBeenCalled()

    // Different object reference but same shallow value
    rerender({ deps: [{ a: 1 }] })
    expect(effect).not.toHaveBeenCalled()
  })

  it('should call the effect if shallow comparison of dependencies changes', () => {
    const effect = vi.fn()
    const { rerender } = renderHook(
      ({ deps }) => useDidUpdateShallow(effect, deps),
      { initialProps: { deps: [{ a: 1 }] } }
    )

    expect(effect).not.toHaveBeenCalled()

    rerender({ deps: [{ a: 2 }] })
    expect(effect).toHaveBeenCalledOnce()
  })
})
