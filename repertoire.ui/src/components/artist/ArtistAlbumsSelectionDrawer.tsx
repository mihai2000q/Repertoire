import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconCircleMinus, IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteAlbumsModal from '../@ui/modal/delete/DeleteAlbumsModal.tsx'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import RemoveAlbumsFromArtistModal from './modal/RemoveAlbumsFromArtistModal.tsx'

interface ArtistAlbumsSelectionDrawerProps {
  artistId: string
  isUnknownArtist: boolean
}

function ArtistAlbumsSelectionDrawer({
  artistId,
  isUnknownArtist: isUnknownArtist
}: ArtistAlbumsSelectionDrawerProps) {
  const { selectedIds, clearSelection } = useClickSelect()

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  return (
    <>
      <SelectionDrawer
        aria-label={'albums-selection-drawer'}
        opened={selectedIds.length > 0}
        onClose={clearSelection}
        text={`${selectedIds.length} album${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            {!isUnknownArtist && (
              <Tooltip label={'Remove albums from artist'}>
                <ActionIcon
                  aria-label={'remove-from-artist'}
                  variant={'grey-primary'}
                  onClick={openRemoveWarning}
                >
                  <IconCircleMinus size={18} />
                </ActionIcon>
              </Tooltip>
            )}
            <Tooltip label={'Delete albums'}>
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
            </Menu.Dropdown>
          )
        }}
      />

      <RemoveAlbumsFromArtistModal
        artistId={artistId}
        ids={selectedIds}
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        onRemove={clearSelection}
      />
      <DeleteAlbumsModal
        ids={selectedIds}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default ArtistAlbumsSelectionDrawer
