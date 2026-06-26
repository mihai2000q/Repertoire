import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import SelectionDrawer from '../../../../../components/@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../../../../../components/@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconChecklist, IconCircleMinus, IconTrash } from '@tabler/icons-react'
import plural from '../../../../../utils/plural.ts'
import AddToPlaylistMenuItem from '../../../../../components/@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteSongsModal from '../../../../../components/@ui/modal/delete/DeleteSongsModal.tsx'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import RemoveSongsFromArtistModal from './modal/RemoveSongsFromArtistModal.tsx'
import CustomRehearsalsModal from '../../../../../components/@ui/modal/CustomRehearsalsModal.tsx'

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

  const [openedCustomRehearsals, { open: openCustomRehearsals, close: closeCustomRehearsals }] =
    useDisclosure(false)

  return (
    <>
      <SelectionDrawer
        aria-label={'songs-selection-drawer'}
        opened={selectedIds.length > 0}
        onClose={clearSelection}
        text={`${selectedIds.length} song${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            {!isUnknownArtist && (
              <Tooltip label={'Remove songs from artist'}>
                <ActionIcon
                  aria-label={'remove-from-artist'}
                  variant={'grey-primary'}
                  onClick={openRemoveWarning}
                >
                  <IconCircleMinus size={18} />
                </ActionIcon>
              </Tooltip>
            )}
            <Tooltip label={'Delete songs'}>
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

      <CustomRehearsalsModal
        opened={openedCustomRehearsals}
        onClose={closeCustomRehearsals}
        ids={selectedIds}
        onSuccess={clearSelection}
      />
    </>
  )
}

export default ArtistSongsSelectionDrawer
