import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import DeleteArtistsModal from '../@ui/modal/delete/DeleteArtistsModal.tsx'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'

function ArtistsSelectionDrawer() {
  const { selectedIds, clearSelection } = useDragSelect()

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  return (
    <>
      <SelectionDrawer
        aria-label={'artists-selection-drawer'}
        opened={selectedIds.length > 0}
        onClose={clearSelection}
        text={`${selectedIds.length} artist${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            <Tooltip label={'Delete artists'}>
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
              <AddToPlaylistMenuItem
                ids={selectedIds}
                type={'artists'}
                closeMenu={closeMenu}
                onSuccess={clearSelection}
              />
              <PerfectRehearsalsMenuItem
                ids={selectedIds}
                closeMenu={closeMenu}
                onSuccess={clearSelection}
                type={'artists'}
              />
            </Menu.Dropdown>
          )
        }}
      />

      <DeleteArtistsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default ArtistsSelectionDrawer
