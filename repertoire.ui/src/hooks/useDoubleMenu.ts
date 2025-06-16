import { useDisclosure } from '@mantine/hooks'

export default function useDoubleMenu() {
  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)
  const [openedContextMenu, { open: openContextMenu, close: closeContextMenu }] =
    useDisclosure(false)

  function handleOpenMenu() {
    if (openedContextMenu) closeContextMenu()
    openMenu()
  }

  function handleOpenContextMenu() {
    if (openedMenu) closeMenu()
    openContextMenu()
  }

  function handleToggleMenu(opened: boolean) {
    if (opened) handleOpenMenu()
    else closeMenu()
  }

  function handleToggleContextMenu(opened: boolean) {
    if (opened) handleOpenContextMenu()
    else closeContextMenu()
  }

  function closeMenus() {
    closeMenu()
    closeContextMenu()
  }

  return {
    openedMenu,
    openMenu: handleOpenMenu,
    closeMenu,
    toggleMenu: handleToggleMenu,
    openedContextMenu,
    openContextMenu: handleOpenContextMenu,
    closeContextMenu,
    toggleContextMenu: handleToggleContextMenu,
    closeMenus
  }
}
