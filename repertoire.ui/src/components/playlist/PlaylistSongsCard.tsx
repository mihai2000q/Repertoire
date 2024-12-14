import {
  ActionIcon,
  Button,
  Card,
  Group,
  LoadingOverlay,
  Menu,
  Space,
  Stack,
  Text
} from '@mantine/core'
import { IconCaretDownFilled, IconCheck, IconDots, IconPlus } from '@tabler/icons-react'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import NewHorizontalCard from '../card/NewHorizontalCard.tsx'
import AddPlaylistSongsModal from './modal/AddPlaylistSongsModal.tsx'
import Playlist from '../../types/models/Playlist.ts'
import { useState } from 'react'
import Order from '../../types/Order.ts'
import { useRemoveSongsFromPlaylistMutation } from '../../state/playlistsApi.ts'
import { useDisclosure } from '@mantine/hooks'

interface PlaylistSongsCardProps {
  playlist: Playlist
  isFetching: boolean
}

function PlaylistSongsCard({ playlist, isFetching }: PlaylistSongsCardProps) {
  const [removeSongsFromPlaylist] = useRemoveSongsFromPlaylistMutation()

  const [order, setOrder] = useState<Order>(playlistSongsOrders[0])

  const [openedAddSongs, { open: openAddSongs, close: closeAddSongs }] = useDisclosure(false)

  function handleRemoveSongsFromPlaylist(songIds: string[]) {
    removeSongsFromPlaylist({ id: playlist.id, songIds })
  }

  return (
    <Card variant={'panel'} h={'100%'} p={0} mx={'xs'}>
      <LoadingOverlay visible={isFetching} />

      <Stack gap={0}>
        <Group px={'md'} pt={'md'} pb={'xs'} align={'center'}>
          <Text fw={600}>Songs</Text>

          <Menu shadow={'sm'}>
            <Menu.Target>
              <Button
                variant={'subtle'}
                size={'compact-xs'}
                rightSection={<IconCaretDownFilled size={11} />}
                styles={{ section: { marginLeft: 4 } }}
              >
                {order.label}
              </Button>
            </Menu.Target>

            <Menu.Dropdown>
              {playlistSongsOrders.map((o) => (
                <Menu.Item
                  key={o.value}
                  leftSection={order === o && <IconCheck size={12} />}
                  onClick={() => setOrder(o)}
                >
                  {o.label}
                </Menu.Item>
              ))}
            </Menu.Dropdown>
          </Menu>

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'}>
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
            <NewHorizontalCard onClick={openAddSongs}>Add Song</NewHorizontalCard>
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
