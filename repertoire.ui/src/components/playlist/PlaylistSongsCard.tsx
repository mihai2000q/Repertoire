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
import {
  useGetPlaylistSongsQuery,
  useMoveSongFromPlaylistMutation
} from '../../state/api/playlistsApi.ts'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import Song from '../../types/models/Song.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import SongProperty from '../../types/enums/SongProperty.ts'
import PlaylistSongsLoader from './loader/PlaylistSongsLoader.tsx'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/LocalStorageKeys.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import { useEffect, useRef } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { setSongsTotalCount } from '../../state/slice/playlistSlice.ts'

interface PlaylistSongsCardProps {
  playlistId: string
}

function PlaylistSongsCard({ playlistId }: PlaylistSongsCardProps) {
  const dispatch = useAppDispatch()

  const [moveSongFromPlaylist, { isLoading: isMoveLoading }] = useMoveSongFromPlaylistMutation()

  const [openedAddSongs, { open: openAddSongs, close: closeAddSongs }] = useDisclosure(false)

  const [order, setOrder] = useLocalStorage({
    key: LocalStorageKeys.PlaylistSongsOrder,
    defaultValue: playlistSongsOrders[0]
  })
  const orderBy = useOrderBy([order])

  const { data, isLoading, isFetching } = useGetPlaylistSongsQuery({
    id: playlistId,
    orderBy: orderBy
  })
  const songs = data?.models
  const totalCount = data?.totalCount
  useEffect(() => {
    dispatch(setSongsTotalCount(totalCount))
  }, [totalCount])

  const [internalSongs, { reorder, setState }] = useListState<Song>(songs)
  useDidUpdate(() => setState(songs), [songs])

  function onSongsDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveSongFromPlaylist({
      id: playlistId,
      playlistSongId: internalSongs[source.index].playlistSongId,
      overPlaylistSongId: internalSongs[destination.index].playlistSongId
    })
  }

  if (isLoading) return <PlaylistSongsLoader />

  return (
    <Card variant={'panel'} aria-label={'songs-card'} p={0} mx={'xs'} mb={'lg'}>
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

          <Menu>
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

        <Stack gap={0} align={'center'}>
          <DragDropContext onDragEnd={onSongsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Box ref={provided.innerRef} {...provided.droppableProps}>
                  {internalSongs.map((song, index) => (
                    <Draggable
                      key={song.playlistSongId}
                      index={index}
                      draggableId={song.playlistSongId}
                      isDragDisabled={
                        isMoveLoading || order.property !== SongProperty.PlaylistTrackNo
                      }
                    >
                      {(provided, snapshot) => (
                        <PlaylistSongCard
                          key={song.playlistSongId}
                          song={song}
                          playlistId={playlistId}
                          order={order}
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

          {internalSongs.length === 0 && (
            <NewHorizontalCard ariaLabel={'new-song-card'} onClick={openAddSongs}>
              Add Songs
            </NewHorizontalCard>
          )}
        </Stack>
      </Stack>

      <AddPlaylistSongsModal
        opened={openedAddSongs}
        onClose={closeAddSongs}
        playlistId={playlistId}
      />
    </Card>
  )
}

export default PlaylistSongsCard
