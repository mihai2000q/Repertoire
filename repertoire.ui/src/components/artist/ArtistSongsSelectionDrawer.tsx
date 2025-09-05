import { ActionIcon, Menu } from '@mantine/core'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconCircleMinus, IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteSongsModal from '../@ui/modal/delete/DeleteSongsModal.tsx'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import RemoveSongsFromArtistModal from './modal/RemoveSongsFromArtistModal.tsx'

interface ArtistSongsSelectionDrawerProps {
  artistId: string
  isUnknownArtist: boolean
}

function ArtistSongsSelectionDrawer({
  artistId,
  isUnknownArtist: isUnknownArtist
}: ArtistSongsSelectionDrawerProps) {
  const { selectedIds, clearSelection } = useClickSelect()

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  return (
    <>
      <SelectionDrawer
        aria-label={'songs-selection-drawer'}
        opened={selectedIds.length > 0}
        onClose={clearSelection}
        text={`${selectedIds.length} song${plural(selectedIds)} selected`}
        actionIcons={
          <>
            {!isUnknownArtist && (
              <ActionIcon
                aria-label={'remove-from-artist'}
                variant={'grey-primary'}
                onClick={openRemoveWarning}
              >
                <IconCircleMinus size={18} />
              </ActionIcon>
            )}
            <ActionIcon aria-label={'delete'} variant={'grey-primary'} onClick={openDeleteWarning}>
              <IconTrash size={18} />
            </ActionIcon>
          </>
        }
        menu={{
          opened: openedMenu,
          toggle: toggleMenu,
          dropdown: (
            <Menu.Dropdown>
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
            </Menu.Dropdown>
          )
        }}
      />

      <RemoveSongsFromArtistModal
        artistId={artistId}
        ids={selectedIds}
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        onRemove={clearSelection}
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

export default ArtistSongsSelectionDrawer
