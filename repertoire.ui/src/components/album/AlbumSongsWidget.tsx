import { ActionIcon, Box, Card, Group, Menu, Space, Stack, Text } from '@mantine/core'
import { IconDots, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import AlbumSongCard from './AlbumSongCard.tsx'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import AddNewAlbumSongModal from './modal/AddNewAlbumSongModal.tsx'
import AddExistingAlbumSongsModal from './modal/AddExistingAlbumSongsModal.tsx'
import { useMoveSongFromAlbumMutation } from '../../state/api/albumsApi.ts'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import Album from '../../types/models/Album.ts'
import Song from '../../types/models/Song.ts'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import albumSongsOrders from '../../data/album/albumSongsOrders.ts'
import Order from '../../types/Order.ts'
import { Dispatch, SetStateAction } from 'react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import LoadingOverlayDebounced from '../@ui/loader/LoadingOverlayDebounced.tsx'
import PerfectRehearsalMenuItem from '../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import { ClickSelectProvider, useClickSelect } from '../../context/ClickSelectContext.tsx'
import AlbumSongsContextMenu from './AlbumSongsContextMenu.tsx'
import AlbumSongsSelectionDrawer from './AlbumSongsSelectionDrawer.tsx'

interface AlbumSongsWidgetProps {
  album: Album | undefined
  songs: Song[] | undefined
  isUnknownAlbum: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  isFetching?: boolean
}

function AlbumSongsWidget({
  album,
  songs,
  isUnknownAlbum,
  order,
  setOrder,
  isFetching
}: AlbumSongsWidgetProps) {
  const [moveSongFromAlbum, { isLoading: isMoveLoading }] = useMoveSongFromAlbumMutation()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  const [internalSongs, { setState }] = useListState<Song>(isUnknownAlbum ? [] : album.songs)
  useDidUpdate(() => setState(album.songs), [album])

  function onSongsDragEnd({ source, destination }) {
    if (!destination || source.index === destination.index) return

    // reorder and change tracking number
    const newSongs = [...album.songs]
    const song = album.songs[source.index]
    newSongs.splice(source.index, 1)
    newSongs.splice(destination.index, 0, song)
    setState(newSongs.map((s, i) => ({ ...s, albumTrackNo: i + 1 })))

    moveSongFromAlbum({
      id: album.id,
      songId: album.songs[source.index].id,
      overSongId: album.songs[destination.index].id
    })
  }

  return (
    <ClickSelectProvider data={songs}>
      <Card aria-label={'songs-widget'} variant={'widget'} h={'100%'} p={0} mx={'xs'} mb={'lg'}>
        <Stack gap={0}>
          <LoadingOverlayDebounced visible={isFetching || isMoveLoading} timeout={750} />

          <Group px={'md'} pt={'md'} pb={'xs'} gap={'xs'}>
            <Text fw={600}>Songs</Text>

            <CompactOrderButton
              availableOrders={albumSongsOrders}
              order={order}
              setOrder={setOrder}
              disabledOrders={isUnknownAlbum ? [albumSongsOrders[0]] : []}
            />

            <Space flex={1} />

            <Menu opened={openedMenu} onOpen={openMenu} onClose={closeMenu}>
              <Menu.Target>
                <ActionIcon aria-label={'songs-more-menu'} size={'md'} variant={'grey'}>
                  <IconDots size={15} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown>
                {!isUnknownAlbum && (
                  <PerfectRehearsalMenuItem id={album.id} closeMenu={closeMenu} type={'album'} />
                )}
                {!isUnknownAlbum && (
                  <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingSongs}>
                    Add Existing Songs
                  </Menu.Item>
                )}
                <Menu.Item leftSection={<IconMusicPlus size={15} />} onClick={openAddNewSong}>
                  Add New Song
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>

          <Stack gap={0}>
            <AlbumSongsContextMenu albumId={album?.id} isUnknownAlbum={isUnknownAlbum}>
              <span style={{ display: 'contents' }}>
                <DragDropContext onDragEnd={onSongsDragEnd}>
                  <Droppable droppableId="dnd-list" direction="vertical">
                    {(provided) => (
                      <Box ref={provided.innerRef} {...provided.droppableProps}>
                        {(isUnknownAlbum ? songs : internalSongs).map((song, index) => {
                          const { isClickSelectionActive } = useClickSelect()
                          return (
                            <Draggable
                              key={song.id}
                              index={index}
                              draggableId={song.id}
                              isDragDisabled={
                                isFetching ||
                                isMoveLoading ||
                                isUnknownAlbum ||
                                order.property !== SongProperty.AlbumTrackNo ||
                                isClickSelectionActive
                              }
                            >
                              {(provided, snapshot) => (
                                <AlbumSongCard
                                  key={song.id}
                                  song={song}
                                  albumId={album?.id}
                                  isUnknownAlbum={isUnknownAlbum}
                                  order={order}
                                  isDragging={snapshot.isDragging}
                                  draggableProvided={provided}
                                  albumImageUrl={album?.imageUrl}
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
              </span>
            </AlbumSongsContextMenu>
            <AlbumSongsSelectionDrawer albumId={album?.id} isUnknownAlbum={isUnknownAlbum} />

            {(isUnknownAlbum || album.songs.length === 0) && (
              <NewHorizontalCard
                ariaLabel={`new-song-card`}
                onClick={isUnknownAlbum ? openAddNewSong : openAddExistingSongs}
              >
                Add New Song{isUnknownAlbum ? '' : 's'}
              </NewHorizontalCard>
            )}
          </Stack>
        </Stack>

        <AddNewAlbumSongModal opened={openedAddNewSong} onClose={closeAddNewSong} album={album} />
        <AddExistingAlbumSongsModal
          opened={openedAddExistingSongs}
          onClose={closeAddExistingSongs}
          albumId={album?.id}
          artistId={album?.artist?.id}
        />
      </Card>
    </ClickSelectProvider>
  )
}

export default AlbumSongsWidget
