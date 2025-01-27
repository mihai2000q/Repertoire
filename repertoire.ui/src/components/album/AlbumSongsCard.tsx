import { ActionIcon, Card, Group, LoadingOverlay, Menu, Space, Stack, Text } from '@mantine/core'
import { IconDots, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import AlbumSongCard from './AlbumSongCard.tsx'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import AddNewAlbumSongModal from './modal/AddNewAlbumSongModal.tsx'
import AddExistingAlbumSongsModal from './modal/AddExistingAlbumSongsModal.tsx'
import { useRemoveSongsFromAlbumMutation } from '../../state/albumsApi.ts'
import { useDisclosure } from '@mantine/hooks'
import Album from '../../types/models/Album.ts'
import Song from '../../types/models/Song.ts'
import CompactOrderButton from '../@ui/button/CompactOrderButton.tsx'
import albumSongsOrders from '../../data/album/albumSongsOrders.ts'
import Order from '../../types/Order.ts'
import { Dispatch, SetStateAction } from 'react'

interface AlbumSongsCardProps {
  album: Album | undefined
  songs: Song[] | undefined
  isUnknownAlbum: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  isFetching?: boolean
}

function AlbumSongsCard({
  album,
  songs,
  isUnknownAlbum,
  order,
  setOrder,
  isFetching
}: AlbumSongsCardProps) {
  const [removeSongsFromAlbum] = useRemoveSongsFromAlbumMutation()

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  function handleRemoveSongsFromAlbum(songIds: string[]) {
    removeSongsFromAlbum({ id: album?.id, songIds })
  }

  return (
    <Card aria-label={'songs-card'} variant={'panel'} h={'100%'} p={0} mx={'xs'} mb={'lg'}>
      <LoadingOverlay visible={isFetching} />

      <Stack gap={0}>
        <Group px={'md'} pt={'md'} pb={'xs'} gap={'xs'}>
          <Text fw={600}>Songs</Text>

          <CompactOrderButton
            availableOrders={albumSongsOrders}
            order={order}
            setOrder={setOrder}
            disabledOrders={isUnknownAlbum ? [albumSongsOrders[0]] : []}
          />

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon aria-label={'songs-more-menu'} size={'md'} variant={'grey'}>
                <IconDots size={15} />
              </ActionIcon>
            </Menu.Target>
            <Menu.Dropdown>
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
          {(isUnknownAlbum ? songs : album.songs).map((song) => (
            <AlbumSongCard
              key={song.id}
              song={song}
              handleRemove={() => handleRemoveSongsFromAlbum([song.id])}
              isUnknownAlbum={isUnknownAlbum}
              order={order}
            />
          ))}
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
  )
}

export default AlbumSongsCard
