import {act, renderHook} from '@testing-library/react'
import useContextMenu from './useContextMenu.ts'
import React from 'react'

describe('use Context Menu', () => {
  it('should return change menu state on open', () => {
    const expectedPageX = 10
    const expectedPageY = 10

    const { result, rerender } = renderHook(() => useContextMenu())

    const [opened, props, { openMenu }] = result.current

    expect(opened).toBeFalsy()
    expect(props).toStrictEqual({
      style: {
        position: 'absolute',
        top: undefined,
        left: undefined
      }
    })

    // open menu
    // noinspection JSUnusedGlobalSymbols
    const event = new class {
      pageX = expectedPageX
      pageY = expectedPageY
    } as React.MouseEvent
    event.preventDefault = vitest.fn()
    act(() => openMenu(event))
    rerender()

    const [opened2, props2, { closeMenu }] = result.current

    expect(event.preventDefault).toHaveBeenCalledOnce()
    expect(opened2).toBeTruthy()
    expect(props2).toStrictEqual({
      style: {
        position: 'absolute',
        top: expectedPageY,
        left: expectedPageX
      }
    })

    // close menu
    act(() => closeMenu())
    rerender()

    const [opened3] = result.current

    expect(opened3).toBeFalsy()
  })
})
