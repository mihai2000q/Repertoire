import { MouseEvent, useState } from 'react'
import { MantineStyleProp } from '@mantine/core'

interface MenuState {
  opened: boolean
  position?: {
    top: number
    left: number
  }
}

interface MenuDropdownProps {
  style: MantineStyleProp
}

interface ContextMenuHandlers {
  openMenu: (event: MouseEvent) => void
  onMenuChange: (opened: boolean) => void
}

export default function useContextMenu(): [boolean, MenuDropdownProps, ContextMenuHandlers] {
  const [menuState, setMenuState] = useState<MenuState>({ opened: false })
  const menuProps: MenuDropdownProps = {
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

  function onMenuChange(opened: boolean) {
    setMenuState({ ...menuState, opened })
  }

  return [menuState.opened, menuProps, { openMenu, onMenuChange }]
}
