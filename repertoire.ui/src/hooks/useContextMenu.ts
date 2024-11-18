import { MouseEvent, useState } from 'react'
import { MantineStyleProp } from '@mantine/core'

interface MenuState {
  opened: boolean
  position?: {
    top: number
    left: number
  }
}

interface MenuProps {
  style: MantineStyleProp
}

interface ContextMenuHandlers {
  openMenu: (event: MouseEvent) => void
  onChange: (opened: boolean) => void
}

export default function useContextMenu(): [boolean, MenuProps, ContextMenuHandlers] {
  const [menuState, setMenuState] = useState<MenuState>({ opened: false })
  const menuProps: MenuProps = {
    style: {
      position: 'absolute',
      top: menuState.position?.top,
      left: menuState.position?.left
    }
  }

  function openMenu(e: MouseEvent) {
    e.preventDefault()
    setMenuState({
      opened: true,
      position: {
        top: e.clientY,
        left: e.clientX
      }
    })
  }

  function onChange(opened: boolean) {
    setMenuState({ ...menuState, opened })
  }

  return [menuState.opened, menuProps, { openMenu, onChange }]
}
