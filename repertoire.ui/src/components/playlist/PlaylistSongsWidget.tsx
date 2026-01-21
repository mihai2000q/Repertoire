import { ActionIcon, Box, Card, Group, Loader, Menu, Space, Stack, Text } from '@mantine/core'
import { IconArrowsShuffle, IconDots, IconPlus } from '@tabler/icons-react'
import playlistSongsOrders from '../../data/playlist/playlistSongsOrders.ts'
import PlaylistSongCard from './PlaylistSongCard.tsx'
import AddPlaylistSongsModal from './modal/AddPlaylistSongsModal.tsx'
import {
  useGetInfinitePlaylistSongsInfiniteQuery,
  useMoveSongFromPlaylistMutation,
  useShufflePlaylistMutation
} from '../../state/api/playlistsApi.ts'
import { useDidUpdate, useDisclosure, useIntersection, useListState } from '@mantine/hooks'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import Song from '../../types/models/Song.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import PlaylistSongsLoader from './loader/PlaylistSongsLoader.tsx'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/keys/LocalStorageKeys.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import { memo, useEffect, useRef } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { setSongsTotalCount } from '../../state/slice/playlistSlice.ts'
import Order from '../../types/Order.ts'
import { MoveSongFromPlaylistRequest } from '../../types/requests/PlaylistRequests.ts'
import LoadingOverlayDebounced from '../@ui/loader/LoadingOverlayDebounced.tsx'
import MenuItemConfirmation from '../@ui/menu/item/MenuItemConfirmation.tsx'
import { Id, toast } from 'react-toastify'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import { useMain } from '../../context/MainContext.tsx'
import PlaylistSongsContextMenu from './PlaylistSongsContextMenu.tsx'
import PlaylistSongsSelectionDrawer from './PlaylistSongsSelectionDrawer.tsx'
import { ClickSelectProvider, useClickSelect } from '../../context/ClickSelectContext.tsx'

interface PlaylistSongsWidgetProps {
  playlistId: string
}

function PlaylistSongsWidget({ playlistId }: PlaylistSongsWidgetProps) {
  const dispatch = useAppDispatch()

  const [moveSongFromPlaylist, { isLoading: isMoveLoading }] = useMoveSongFromPlaylistMutation()
  const [shufflePlaylistSongs, { isLoading: isShuffleLoading }] = useShufflePlaylistMutation()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedAddSongs, { open: openAddSongs, close: closeAddSongs }] = useDisclosure(false)

  const [order, setOrder] = useLocalStorage({
    key: LocalStorageKeys.PlaylistSongsOrder,
    defaultValue: playlistSongsOrders[0]
  })
  const orderBy = useOrderBy([order])

  const { data, isLoading, isFetching, isFetchingNextPage, fetchNextPage } =
    useGetInfinitePlaylistSongsInfiniteQuery({
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

  const { mainScroll } = useMain()
  const { ref: lastSongRef, entry } = useIntersection({
    root: mainScroll.ref.current,
    threshold: 0.1
  })
  useEffect(() => {
    if (entry?.isIntersecting === true) fetchNextPage()
  }, [entry?.isIntersecting])

  const shuffleToastId = useRef<Id>()
  async function handleShuffle() {
    await shufflePlaylistSongs({ id: playlistId }).unwrap()
    if (shuffleToastId.current) toast.dismiss(shuffleToastId.current)
    shuffleToastId.current = toast.info('Playlist shuffled!')
  }

  if (isLoading) return <PlaylistSongsLoader />

  return (
    <ClickSelectProvider data={songs}>
      <Card variant={'widget'} aria-label={'songs-widget'} p={0} mx={'xs'} mb={'lg'}>
        <Stack gap={0}>
          <LoadingOverlayDebounced visible={isFetching || isMoveLoading} timeout={750} />

          <Group px={'md'} pt={'md'} pb={'xs'} gap={'xs'}>
            <Text fw={600}>Songs</Text>

            <CompactOrderButton
              availableOrders={playlistSongsOrders}
              order={order}
              setOrder={setOrder}
            />

            <Space flex={1} />

            <Menu opened={openedMenu} onOpen={openMenu} onClose={closeMenu}>
              <Menu.Target>
                <ActionIcon size={'md'} variant={'grey'} aria-label={'songs-more-menu'}>
                  <IconDots size={15} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown miw={125}>
                <MenuItemConfirmation
                  leftSection={<IconArrowsShuffle size={15} />}
                  isLoading={isShuffleLoading}
                  onConfirm={handleShuffle}
                >
                  Shuffle
                </MenuItemConfirmation>
                <PerfectRehearsalMenuItem id={playlistId} closeMenu={closeMenu} type={'playlist'} />
                <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddSongs}>
                  Add Songs
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>

          <Stack gap={0}>
            <PlaylistSongsContextMenu playlistId={playlistId} songs={songs}>
              <span style={{ display: 'contents' }}>
                <Songs
                  songs={songs}
                  order={order}
                  playlistId={playlistId}
                  moveSongFromPlaylist={moveSongFromPlaylist}
                  isMoveLoading={isMoveLoading}
                  isFetching={isFetching}
                />
              </span>
            </PlaylistSongsContextMenu>
            <PlaylistSongsSelectionDrawer playlistId={playlistId} songs={songs} />

            {songs.length === 0 && (
              <NewHorizontalCard ariaLabel={'new-song-card'} onClick={openAddSongs}>
                Add Songs
              </NewHorizontalCard>
            )}

            <Stack gap={0} align={'center'}>
              <div ref={lastSongRef} />
              {isFetchingNextPage && <Loader size={30} m={'md'} />}
            </Stack>
          </Stack>
        </Stack>

        <AddPlaylistSongsModal
          opened={openedAddSongs}
          onClose={closeAddSongs}
          playlistId={playlistId}
        />
      </Card>
    </ClickSelectProvider>
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
              {internalSongs.map((song, index) => {
                const { isClickSelectionActive } = useClickSelect()
                return (
                  <Draggable
                    key={song.playlistSongId}
                    index={index}
                    draggableId={song.playlistSongId}
                    isDragDisabled={
                      isFetching ||
                      isMoveLoading ||
                      order.property !== SongProperty.PlaylistTrackNo ||
                      isClickSelectionActive
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
                )
              })}
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

export default PlaylistSongsWidget
