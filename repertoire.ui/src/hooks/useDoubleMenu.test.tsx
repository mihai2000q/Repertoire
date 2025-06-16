import { act, renderHook } from '@testing-library/react'
import useDoubleMenu from './useDoubleMenu.ts'

describe('use Double Menu', () => {
  it('should initialize with both menus closed', () => {
    const { result } = renderHook(() => useDoubleMenu())

    expect(result.current.openedMenu).toBe(false)
    expect(result.current.openedContextMenu).toBe(false)
  })

  it('should open menu and close context menu when opening menu', () => {
    const { result } = renderHook(() => useDoubleMenu())

    // Open context menu first
    act(() => {
      result.current.openContextMenu()
    })
    expect(result.current.openedContextMenu).toBe(true)

    // Then open the menu
    act(() => {
      result.current.openMenu()
    })

    expect(result.current.openedMenu).toBe(true)
    expect(result.current.openedContextMenu).toBe(false)
  })

  it('should open context menu and close menu when opening context menu', () => {
    const { result } = renderHook(() => useDoubleMenu())

    // Open menu first
    act(() => {
      result.current.openMenu()
    })
    expect(result.current.openedMenu).toBe(true)

    // Then open the context menu
    act(() => {
      result.current.openContextMenu()
    })

    expect(result.current.openedContextMenu).toBe(true)
    expect(result.current.openedMenu).toBe(false)
  })

  it('should handle menu onChange events correctly', () => {
    const { result } = renderHook(() => useDoubleMenu())

    // Test opening via onChange
    act(() => {
      result.current.onChangeMenu(true)
    })
    expect(result.current.openedMenu).toBe(true)

    // Test closing via onChange
    act(() => {
      result.current.onChangeMenu(false)
    })
    expect(result.current.openedMenu).toBe(false)
  })

  it('should handle context menu onChange events correctly', () => {
    const { result } = renderHook(() => useDoubleMenu())

    // Test opening via onChange
    act(() => {
      result.current.onChangeContextMenu(true)
    })
    expect(result.current.openedContextMenu).toBe(true)

    // Test closing via onChange
    act(() => {
      result.current.onChangeContextMenu(false)
    })
    expect(result.current.openedContextMenu).toBe(false)
  })

  it('should close both menus with closeMenus', () => {
    const { result } = renderHook(() => useDoubleMenu())

    // Open both menus
    act(() => {
      result.current.openMenu()
      result.current.openContextMenu()
    })

    // Close both
    act(() => {
      result.current.closeMenus()
    })

    expect(result.current.openedMenu).toBe(false)
    expect(result.current.openedContextMenu).toBe(false)
  })

  it('should not throw when closing already closed menus', () => {
    const { result } = renderHook(() => useDoubleMenu())

    expect(result.current.openedMenu).toBe(false)
    expect(result.current.openedContextMenu).toBe(false)

    act(() => {
      expect(() => result.current.closeMenu()).not.toThrow()
      expect(() => result.current.closeContextMenu()).not.toThrow()
      expect(() => result.current.closeMenus()).not.toThrow()
    })
  })
})
