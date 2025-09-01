import { Menu } from '@mantine/core'
import { IconTrash } from '@tabler/icons-react'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { ReactNode, useEffect } from 'react'
import DeletePlaylistsModal from '../@ui/modal/delete/DeletePlaylistsModal.tsx'

function PlaylistsContextMenu({ children }: { children: ReactNode }) {
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
          <PerfectRehearsalsMenuItem
            ids={selectedIds}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
            type={'playlists'}
          />
          <Menu.Divider />

          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <DeletePlaylistsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default PlaylistsContextMenu
