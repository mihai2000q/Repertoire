import { act, renderHook } from '@testing-library/react'
import useContextMenu from './useContextMenu.ts'
import React from 'react'

describe('useAuth', () => {
  it('should return change menu state on open', () => {
    const expectedClientX = 10
    const expectedClientY = 10

    const { result, rerender } = renderHook(() => useContextMenu())

    const [opened, props, { openMenu }] = result.current

    // open menu
    expect(opened).toBeFalsy()
    expect(props).toStrictEqual({
      style: {
        position: 'absolute',
        top: undefined,
        left: undefined
      }
    })

    act(() =>
      openMenu(
        new MouseEvent('mousemove', {
          clientX: expectedClientX,
          clientY: expectedClientY
        }) as unknown as React.MouseEvent
      )
    )
    rerender()

    // on menu change
    const [opened2, props2, { onMenuChange }] = result.current

    expect(opened2).toBeTruthy()
    expect(props2).toStrictEqual({
      style: {
        position: 'absolute',
        top: expectedClientY,
        left: expectedClientX
      }
    })

    act(() => onMenuChange(false))
    rerender()

    const [opened3] = result.current

    expect(opened3).toBeFalsy()
  })
})
