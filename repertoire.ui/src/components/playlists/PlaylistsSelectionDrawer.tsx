import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import DeletePlaylistsModal from '../@ui/modal/delete/DeletePlaylistsModal.tsx'

function PlaylistsSelectionDrawer() {
  const { selectedIds, clearSelection } = useDragSelect()

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  return (
    <>
      <SelectionDrawer
        aria-label={'playlists-selection-drawer'}
        opened={selectedIds.length > 0}
        onClose={clearSelection}
        text={`${selectedIds.length} playlist${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            <Tooltip label={'Delete playlists'}>
              <ActionIcon
                aria-label={'delete'}
                variant={'grey-primary'}
                onClick={openDeleteWarning}
              >
                <IconTrash size={18} />
              </ActionIcon>
            </Tooltip>
          </Tooltip.Group>
        }
        menu={{
          opened: openedMenu,
          toggle: toggleMenu,
          dropdown: (
            <Menu.Dropdown>
              <PerfectRehearsalsMenuItem
                ids={selectedIds}
                closeMenu={closeMenu}
                onSuccess={clearSelection}
                type={'playlists'}
              />
            </Menu.Dropdown>
          )
        }}
      />

      <DeletePlaylistsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default PlaylistsSelectionDrawer
