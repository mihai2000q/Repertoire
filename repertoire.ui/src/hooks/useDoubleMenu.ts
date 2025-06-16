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

  function handleMenuOnChange(opened: boolean) {
    if (opened) handleOpenMenu()
    else closeMenu()
  }

  function handleContextMenuOnChange(opened: boolean) {
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
    onChangeMenu: handleMenuOnChange,
    openedContextMenu,
    openContextMenu: handleOpenContextMenu,
    closeContextMenu,
    onChangeContextMenu: handleContextMenuOnChange,
    closeMenus
  }
}
