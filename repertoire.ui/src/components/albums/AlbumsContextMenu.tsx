import { Menu } from '@mantine/core'
import { IconTrash } from '@tabler/icons-react'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { ReactNode, useEffect } from 'react'
import DeleteAlbumsModal from '../@ui/modal/delete/DeleteAlbumsModal.tsx'

function AlbumsContextMenu({ children }: { children: ReactNode }) {
  const { selectedIds, clearSelection } = useDragSelect()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  useEffect(() => {
    if (selectedIds.length === 0) closeMenu()
  }, [selectedIds])

  return (
    <>
      <ContextMenu
        opened={openedMenu}
        onClose={closeMenu}
        onOpen={openMenu}
        disabled={selectedIds.length === 0}
      >
        <ContextMenu.Target>{children}</ContextMenu.Target>

        <ContextMenu.Dropdown>
          <AddToPlaylistMenuItem
            ids={selectedIds}
            type={'albums'}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
          />
          <PerfectRehearsalsMenuItem
            ids={selectedIds}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
            type={'albums'}
          />
          <Menu.Divider />

          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <DeleteAlbumsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default AlbumsContextMenu
