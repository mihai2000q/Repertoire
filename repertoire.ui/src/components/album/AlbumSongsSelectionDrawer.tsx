import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconCircleMinus, IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import DeleteSongsModal from '../@ui/modal/delete/DeleteSongsModal.tsx'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import RemoveSongsFromAlbumModal from './modal/RemoveSongsFromAlbumModal.tsx'

interface AlbumSongsSelectionDrawerProps {
  albumId: string
  isUnknownAlbum: boolean
}

function AlbumSongsSelectionDrawer({
  albumId,
  isUnknownAlbum: isUnknownAlbum
}: AlbumSongsSelectionDrawerProps) {
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
          <Tooltip.Group openDelay={200}>
            {!isUnknownAlbum && (
              <Tooltip label={'Remove songs from album'}>
                <ActionIcon
                  aria-label={'remove-from-album'}
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
            </Menu.Dropdown>
          )
        }}
      />

      <RemoveSongsFromAlbumModal
        albumId={albumId}
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

export default AlbumSongsSelectionDrawer
