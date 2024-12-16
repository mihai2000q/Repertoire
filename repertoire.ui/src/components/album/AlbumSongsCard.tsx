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

interface AlbumSongsCardProps {
  album: Album | undefined
  songs: Song[] | undefined
  isFetching: boolean
  isUnknownAlbum: boolean
}

function AlbumSongsCard({ album, songs, isFetching, isUnknownAlbum }: AlbumSongsCardProps) {
  const [removeSongsFromAlbum] = useRemoveSongsFromAlbumMutation()

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  function handleRemoveSongsFromAlbum(songIds: string[]) {
    removeSongsFromAlbum({ id: album?.id, songIds })
  }

  return (
    <Card variant={'panel'} h={'100%'} p={0} mx={'xs'} mb={'lg'}>
      <LoadingOverlay visible={isFetching} />

      <Stack gap={0}>
        <Group px={'md'} pt={'md'} pb={'xs'}>
          <Text fw={600}>Songs</Text>
          <Space flex={1} />
          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'}>
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
            />
          ))}
          {(isUnknownAlbum || album.songs.length === 0) && (
            <NewHorizontalCard onClick={isUnknownAlbum ? openAddNewSong : openAddExistingSongs}>
              Add New Song{isUnknownAlbum ? '' : 's'}
            </NewHorizontalCard>
          )}
        </Stack>
      </Stack>

      <AddNewAlbumSongModal
        opened={openedAddNewSong}
        onClose={closeAddNewSong}
        albumId={album?.id}
      />
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
