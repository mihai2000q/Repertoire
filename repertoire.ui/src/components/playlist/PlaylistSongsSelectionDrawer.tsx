import { ActionIcon, Menu } from '@mantine/core'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../@ui/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconCircleMinus } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import AddToPlaylistMenuItem from '../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import RemoveSongsFromPlaylistModal from './modal/RemoveSongsFromPlaylistModal.tsx'
import Song from '../../types/models/Song.ts'
import { useEffect, useState } from 'react'

interface PlaylistSongsSelectionDrawerProps {
  playlistId: string
  songs: Song[]
}

function PlaylistSongsSelectionDrawer({ playlistId, songs }: PlaylistSongsSelectionDrawerProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const [selectedSongsIds, setSelectedSongsIds] = useState<string[]>([])
  useEffect(() => {
    setSelectedSongsIds(
      songs.filter((s) => selectedIds.some((psId) => psId === s.playlistSongId)).map((s) => s.id)
    )
  }, [selectedIds])

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
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
            <ActionIcon
              aria-label={'remove-from-playlist'}
              variant={'grey-primary'}
              onClick={openRemoveWarning}
            >
              <IconCircleMinus size={18} />
            </ActionIcon>
          </>
        }
        menu={{
          opened: openedMenu,
          toggle: toggleMenu,
          dropdown: (
            <Menu.Dropdown>
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
            </Menu.Dropdown>
          )
        }}
      />

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

export default PlaylistSongsSelectionDrawer
