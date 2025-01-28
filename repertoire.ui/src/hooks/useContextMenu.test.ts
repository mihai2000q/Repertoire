import {act, renderHook} from '@testing-library/react'
import useContextMenu from './useContextMenu.ts'
import React from 'react'

describe('use Context Menu', () => {
  it('should return change menu state on open', () => {
    const expectedPageX = 10
    const expectedPageY = 10

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

    // noinspection JSUnusedGlobalSymbols
    const event = new class {
      pageX = expectedPageX
      pageY = expectedPageY
      preventDefault() {}
    }

    act(() => openMenu(event as React.MouseEvent))
    rerender()

    // close menu
    const [opened2, props2, { closeMenu }] = result.current

    expect(opened2).toBeTruthy()
    expect(props2).toStrictEqual({
      style: {
        position: 'absolute',
        top: expectedPageY,
        left: expectedPageX
      }
    })

    act(() => closeMenu())
    rerender()

    const [opened3] = result.current

    expect(opened3).toBeFalsy()
  })
})
