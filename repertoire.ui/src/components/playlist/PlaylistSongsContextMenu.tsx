import { Menu } from '@mantine/core'
import { IconCircleMinus } from '@tabler/icons-react'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { ReactNode, useEffect, useState } from 'react'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import RemoveSongsFromPlaylistModal from './modal/RemoveSongsFromPlaylistModal.tsx'
import Song from '../../types/models/Song.ts'

interface PlaylistSongsContextMenuProps {
  children: ReactNode
  playlistId: string
  songs: Song[]
}

function PlaylistSongsContextMenu({ children, playlistId, songs }: PlaylistSongsContextMenuProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const [selectedSongsIds, setSelectedSongsIds] = useState<string[]>([])
  useEffect(() => {
    setSelectedSongsIds(
      songs.filter((s) => selectedIds.some((psId) => psId === s.playlistSongId)).map((s) => s.id)
    )
  }, [selectedIds])

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
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
            ids={selectedSongsIds}
            type={'songs'}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
          />
          <PerfectRehearsalsMenuItem
            ids={selectedSongsIds}
            closeMenu={closeMenu}
            onSuccess={clearSelection}
            type={'songs'}
          />
          <Menu.Divider />
          <Menu.Item leftSection={<IconCircleMinus size={14} />} onClick={openRemoveWarning}>
            Remove from Playlist
          </Menu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <RemoveSongsFromPlaylistModal
        playlistId={playlistId}
        ids={selectedIds}
        opened={openedRemoveWarning}
        onClose={closeRemoveWarning}
        onRemove={clearSelection}
      />
    </>
  )
}

export default PlaylistSongsContextMenu
