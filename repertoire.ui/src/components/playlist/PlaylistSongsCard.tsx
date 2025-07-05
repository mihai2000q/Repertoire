import { ActionIcon, Box, Card, Group, Loader, Menu, Space, Stack, Text } from '@mantine/core'
import { IconDots, IconPlus } from '@tabler/icons-react'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import AddPlaylistSongsModal from './modal/AddPlaylistSongsModal.tsx'
import {
  useGetPlaylistSongsInfiniteQuery,
  useMoveSongFromPlaylistMutation
} from '../../state/api/playlistsApi.ts'
import { useDidUpdate, useDisclosure, useIntersection, useListState } from '@mantine/hooks'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import Song from '../../types/models/Song.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import SongProperty from '../../types/enums/SongProperty.ts'
import PlaylistSongsLoader from './loader/PlaylistSongsLoader.tsx'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/LocalStorageKeys.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import { memo, useEffect } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { setSongsTotalCount } from '../../state/slice/playlistSlice.ts'
import useMainScroll from '../../hooks/useMainScroll.ts'
import Order from '../../types/Order.ts'
import { MoveSongFromPlaylistRequest } from '../../types/requests/PlaylistRequests.ts'

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

  const { data, isLoading, isFetching, isFetchingNextPage, fetchNextPage } =
    useGetPlaylistSongsInfiniteQuery({
      id: playlistId,
      pageSize: 25,
      orderBy: orderBy
    })
  const songs = data?.pages.flatMap((x) => x.models ?? []) ?? []
  const totalCount = data?.pages[0].totalCount
  useEffect(() => {
    dispatch(setSongsTotalCount(totalCount))
    return () => {
      dispatch(setSongsTotalCount(undefined))
    }
  }, [totalCount])

  const { ref: mainScrollRef } = useMainScroll()
  const { ref: lastSongRef, entry } = useIntersection({
    root: mainScrollRef.current,
    threshold: 0.1
  })
  useEffect(() => {
    if (entry?.isIntersecting === true) fetchNextPage()
  }, [entry?.isIntersecting])

  if (isLoading) return <PlaylistSongsLoader />

  return (
    <Card variant={'panel'} aria-label={'songs-card'} p={0} mx={'xs'} mb={'lg'}>
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

        <Stack gap={0}>
          <Songs
            songs={songs}
            order={order}
            playlistId={playlistId}
            moveSongFromPlaylist={moveSongFromPlaylist}
            isMoveLoading={isMoveLoading}
            isFetching={isFetching}
          />

          {songs.length === 0 && (
            <NewHorizontalCard ariaLabel={'new-song-card'} onClick={openAddSongs}>
              Add Songs
            </NewHorizontalCard>
          )}

          <div ref={lastSongRef} />
          {isFetchingNextPage && <Loader size={30} m={'md'} style={{ alignSelf: 'center' }} />}
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

const Songs = memo(
  ({
    songs,
    order,
    playlistId,
    moveSongFromPlaylist,
    isMoveLoading,
    isFetching
  }: {
    songs: Song[]
    order: Order
    playlistId: string
    moveSongFromPlaylist: (request: MoveSongFromPlaylistRequest) => void
    isMoveLoading: boolean
    isFetching: boolean
  }) => {
    const [internalSongs, { setState }] = useListState<Song>(songs)
    useDidUpdate(() => setState(songs), [JSON.stringify(songs)])

    function onSongsDragEnd({ source, destination }) {
      if (!destination || source.index === destination.index) return

      // reorder and change tracking number
      const newSongs = [...songs]
      const song = songs[source.index]
      newSongs.splice(source.index, 1)
      newSongs.splice(destination.index, 0, song)
      setState(newSongs.map((s, i) => ({ ...s, playlistTrackNo: i + 1 })))

      moveSongFromPlaylist({
        id: playlistId,
        playlistSongId: songs[source.index].playlistSongId,
        overPlaylistSongId: songs[destination.index].playlistSongId
      })
    }

    return (
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
                    isFetching || isMoveLoading || order.property !== SongProperty.PlaylistTrackNo
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
    )
  },
  (prevProps, nextProps) => {
    return (
      JSON.stringify(prevProps.songs) === JSON.stringify(nextProps.songs) &&
      JSON.stringify(prevProps.order) === JSON.stringify(nextProps.order) &&
      prevProps.playlistId === nextProps.playlistId &&
      prevProps.isFetching === nextProps.isFetching &&
      prevProps.isMoveLoading === nextProps.isMoveLoading
    )
  }
)

Songs.displayName = 'Songs'

export default PlaylistSongsCard
