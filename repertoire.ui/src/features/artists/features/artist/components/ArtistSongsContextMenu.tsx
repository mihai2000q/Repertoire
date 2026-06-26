import { Menu } from '@mantine/core'
import { IconChecklist, IconCircleMinus, IconTrash } from '@tabler/icons-react'
import AddToPlaylistMenuItem from '../../../../../components/@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../../../../../components/@ui/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import PerfectRehearsalsMenuItem from '../../../../../components/@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { ReactNode, useEffect } from 'react'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import DeleteSongsModal from '../../../../../components/@ui/modal/delete/DeleteSongsModal.tsx'
import RemoveSongsFromArtistModal from './modal/RemoveSongsFromArtistModal.tsx'
import CustomRehearsalsModal from '../../../../../components/@ui/modal/CustomRehearsalsModal.tsx'

interface ArtistSongsContextMenuProps {
  children: ReactNode
  artistId: string
  isUnknownArtist: boolean
}

function ArtistSongsContextMenu({
  children,
  artistId,
  isUnknownArtist: isUnknownArtist
}: ArtistSongsContextMenuProps) {
  const { selectedIds, clearSelection } = useClickSelect()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedCustomRehearsals, { open: openCustomRehearsals, close: closeCustomRehearsals }] =
    useDisclosure(false)
  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
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

          {!isUnknownArtist && (
            <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={openRemoveWarning}>
              Remove from Artist
            </Menu.Item>
          )}
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

export default ArtistSongsContextMenu
