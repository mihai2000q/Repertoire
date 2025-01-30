import {
  ActionIcon,
  Box,
  Card,
  Group,
  LoadingOverlay,
  Menu,
  Space,
  Stack,
  Text
} from '@mantine/core'
import { IconDots, IconPlus } from '@tabler/icons-react'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import AddPlaylistSongsModal from './modal/AddPlaylistSongsModal.tsx'
import Playlist from '../../types/models/Playlist.ts'
import { useState } from 'react'
import Order from '../../types/Order.ts'
import {
  useMoveSongFromPlaylistMutation,
  useRemoveSongsFromPlaylistMutation
} from '../../state/api/playlistsApi.ts'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import Song from '../../types/models/Song.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import SongProperty from '../../utils/enums/SongProperty.ts'

interface PlaylistSongsCardProps {
  playlist: Playlist
  isFetching?: boolean
}

function PlaylistSongsCard({ playlist, isFetching }: PlaylistSongsCardProps) {
  const [moveSongFromPlaylist, { isLoading: isMoveLoading }] = useMoveSongFromPlaylistMutation()
  const [removeSongsFromPlaylist] = useRemoveSongsFromPlaylistMutation()

  const [order, setOrder] = useState<Order>(playlistSongsOrders[0])

  const [openedAddSongs, { open: openAddSongs, close: closeAddSongs }] = useDisclosure(false)

  const [internalSongs, { reorder, setState }] = useListState<Song>(playlist.songs)
  useDidUpdate(() => setState(playlist.songs), [playlist])

  function handleRemoveSongsFromPlaylist(songIds: string[]) {
    removeSongsFromPlaylist({ id: playlist.id, songIds })
  }

  function onSongsDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveSongFromPlaylist({
      id: playlist.id,
      songId: playlist.songs[source.index].id,
      overSongId: playlist.songs[destination.index].id
    })
  }

  return (
    <Card variant={'panel'} aria-label={'songs-card'} h={'100%'} p={0} mx={'xs'} mb={'lg'}>
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
          <DragDropContext onDragEnd={onSongsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Box ref={provided.innerRef} {...provided.droppableProps}>
                  {internalSongs.map((song, index) => (
                    <Draggable
                      key={song.id}
                      index={index}
                      draggableId={song.id}
                      isDragDisabled={
                        isMoveLoading || order.property !== SongProperty.PlaylistTrackNo
                      }
                    >
                      {(provided, snapshot) => (
                        <PlaylistSongCard
                          key={song.id}
                          song={song}
                          handleRemove={() => handleRemoveSongsFromPlaylist([song.id])}
                          isDragging={snapshot.isDragging}
                          draggableProvided={provided}
                        />
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </Box>
              )}
            </Droppable>
          </DragDropContext>

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
