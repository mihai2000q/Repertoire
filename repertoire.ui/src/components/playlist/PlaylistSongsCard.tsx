import { ActionIcon, Card, Group, LoadingOverlay, Menu, Space, Stack, Text } from '@mantine/core'
import { IconDots, IconPlus } from '@tabler/icons-react'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import AddPlaylistSongsModal from './modal/AddPlaylistSongsModal.tsx'
import Playlist from '../../types/models/Playlist.ts'
import { useState } from 'react'
import Order from '../../types/Order.ts'
import { useRemoveSongsFromPlaylistMutation } from '../../state/playlistsApi.ts'
import { useDisclosure } from '@mantine/hooks'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'

interface PlaylistSongsCardProps {
  playlist: Playlist
  isFetching?: boolean
}

function PlaylistSongsCard({ playlist, isFetching }: PlaylistSongsCardProps) {
  const [removeSongsFromPlaylist] = useRemoveSongsFromPlaylistMutation()

  const [order, setOrder] = useState<Order>(playlistSongsOrders[0])

  const [openedAddSongs, { open: openAddSongs, close: closeAddSongs }] = useDisclosure(false)

  function handleRemoveSongsFromPlaylist(songIds: string[]) {
    removeSongsFromPlaylist({ id: playlist.id, songIds })
  }

  return (
    <Card variant={'panel'} aria-label={'songs-card'} h={'100%'} p={0} mx={'xs'}>
      <LoadingOverlay visible={isFetching} />

      <Stack gap={0}>
        <Group px={'md'} pt={'md'} pb={'xs'} gap={'xs'}>
          <Text fw={600}>Songs</Text>

          <CompactOrderButton
            availableOrders={playlistSongsOrders}
            order={order}
            setOrder={setOrder}
          />

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'} aria-label={'songs-more-menu'}>
                <IconDots size={15} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddSongs}>
                Add Songs
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>

        <Stack gap={0}>
          {playlist.songs.map((song) => (
            <PlaylistSongCard
              key={song.id}
              song={song}
              handleRemove={() => handleRemoveSongsFromPlaylist([song.id])}
            />
          ))}
          {playlist.songs.length === 0 && (
            <NewHorizontalCard ariaLabel={'new-song-card'} onClick={openAddSongs}>
              Add Songs
            </NewHorizontalCard>
          )}
        </Stack>
      </Stack>

      <AddPlaylistSongsModal
        opened={openedAddSongs}
        onClose={closeAddSongs}
        playlistId={playlist.id}
      />
    </Card>
  )
}

export default PlaylistSongsCard
