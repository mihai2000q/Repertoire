import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import SelectionDrawer from '../../../../../components/drawer/SelectionDrawer.tsx'
import PerfectRehearsalsMenuItem from '../../../../../components/menu/item/PerfectRehearsalsMenuItem.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconChecklist, IconCircleMinus } from '@tabler/icons-react'
import plural from '../../../../../utils/plural.ts'
import AddToPlaylistMenuItem from '../../../../../components/menu/item/AddToPlaylistMenuItem.tsx'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import RemoveSongsFromPlaylistModal from './modal/RemoveSongsFromPlaylistModal.tsx'
import Song from '../../../../../types/models/Song.ts'
import { useEffect, useState } from 'react'
import CustomRehearsalsModal from '../../../../../components/modal/CustomRehearsalsModal.tsx'

interface PlaylistSongsSelectionDrawerProps {
  playlistId: string
  songs: Song[]
}

function PlaylistSongsSelectionDrawer({ playlistId, songs }: PlaylistSongsSelectionDrawerProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const [selectedSongIds, setSelectedSongIds] = useState<string[]>([])
  useEffect(() => {
    setSelectedSongIds(
      songs.filter((s) => selectedIds.some((psId) => psId === s.playlistSongId)).map((s) => s.id)
    )
  }, [selectedIds])

  const [openedMenu, { close: closeMenu, toggle: toggleMenu }] = useDisclosure(false)

  const [openedRemoveWarning, { open: openRemoveWarning, close: closeRemoveWarning }] =
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
            <Tooltip label={'Remove songs from playlist'}>
              <ActionIcon
                aria-label={'remove-from-playlist'}
                variant={'grey-primary'}
                onClick={openRemoveWarning}
              >
                <IconCircleMinus size={18} />
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
                ids={selectedSongIds}
                type={'songs'}
                closeMenu={closeMenu}
                onSuccess={clearSelection}
              />
              <PerfectRehearsalsMenuItem
                ids={selectedIds}
                closeMenu={closeMenu}
                onSuccess={clearSelection}
                type={'playlist-songs'}
              />
              <Menu.Item leftSection={<IconChecklist size={14} />} onClick={openCustomRehearsals}>
                Custom Rehearsals
              </Menu.Item>
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

      <CustomRehearsalsModal
        opened={openedCustomRehearsals}
        onClose={closeCustomRehearsals}
        ids={selectedSongIds}
        onSuccess={clearSelection}
      />
    </>
  )
}

export default PlaylistSongsSelectionDrawer
