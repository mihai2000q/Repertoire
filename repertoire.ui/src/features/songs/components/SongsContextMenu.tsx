import { Menu } from '@mantine/core'
import { IconChecklist, IconTrash } from '@tabler/icons-react'
import AddToPlaylistMenuItem from '../../../components/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../../../components/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useDragSelect } from '../../../context/DragSelectContext.tsx'
import PerfectRehearsalsMenuItem from '../../../components/menu/item/PerfectRehearsalsMenuItem.tsx'
import { ReactNode, useEffect } from 'react'
import DeleteSongsModal from '../../../components/modal/delete/DeleteSongsModal.tsx'
import CustomRehearsalsModal from '../../../components/modal/CustomRehearsalsModal.tsx'

function SongsContextMenu({ children }: { children: ReactNode }) {
  const { selectedIds, clearSelection } = useDragSelect()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedCustomRehearsals, { open: openCustomRehearsals, close: closeCustomRehearsals }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  useEffect(() => {
    if (selectedIds.length === 0) closeMenu()
  }, [selectedIds])

  return (
    <>
      <ContextMenu
        aria-label={'songs-context-menu'}
        opened={openedMenu}
        onClose={closeMenu}
        onOpen={openMenu}
        disabled={selectedIds.length === 0}
      >
        <ContextMenu.Target>{children}</ContextMenu.Target>

        <ContextMenu.Dropdown>
          <AddToPlaylistMenuItem
            ids={selectedIds}
            type={'songs'}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
          />
          <PerfectRehearsalsMenuItem
            ids={selectedIds}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
            type={'songs'}
          />
          <Menu.Item leftSection={<IconChecklist size={14} />} onClick={openCustomRehearsals}>
            Custom Rehearsals
          </Menu.Item>
          <Menu.Divider />

          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <CustomRehearsalsModal
        opened={openedCustomRehearsals}
        onClose={closeCustomRehearsals}
        ids={selectedIds}
        onSuccess={clearSelection}
      />
      <DeleteSongsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default SongsContextMenu
