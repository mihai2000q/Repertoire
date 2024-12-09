import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Card,
  Divider,
  Group,
  Image,
  LoadingOverlay,
  Menu,
  Space,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useNavigate, useParams } from 'react-router-dom'
import {
  useDeleteAlbumMutation,
  useGetAlbumQuery,
  useRemoveSongsFromAlbumMutation
} from '../state/albumsApi.ts'
import AlbumLoader from '../components/album/AlbumLoader.tsx'
import albumPlaceholder from '../assets/image-placeholder-1.jpg'
import unknownPlaceholder from '../assets/unknown-placeholder.png'
import AlbumSongCard from '../components/album/AlbumSongCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconDots, IconEdit, IconMusicPlus, IconPlus, IconTrash } from '@tabler/icons-react'
import AddNewAlbumSongModal from '../components/album/modal/AddNewAlbumSongModal.tsx'
import AddExistingAlbumSongsModal from '../components/album/modal/AddExistingAlbumSongsModal.tsx'
import userPlaceholder from '../assets/user-placeholder.jpg'
import { useAppDispatch } from '../state/store.ts'
import { openArtistDrawer } from '../state/globalSlice.ts'
import dayjs from 'dayjs'
import NewHorizontalCard from '../components/card/NewHorizontalCard.tsx'
import HeaderPanelCard from '../components/card/HeaderPanelCard.tsx'
import { toast } from 'react-toastify'
import EditAlbumHeaderModal from '../components/album/modal/EditAlbumHeaderModal.tsx'
import { useGetSongsQuery } from '../state/songsApi.ts'

function Album() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const params = useParams()
  const albumId = params['id'] ?? ''

  const [deleteAlbumMutation] = useDeleteAlbumMutation()

  const isUnknownAlbum = albumId === 'unknown'

  const { data: album, isLoading, isFetching } = useGetAlbumQuery(albumId, { skip: isUnknownAlbum })

  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery(
    {
      orderBy: ['title'],
      searchBy: ['album_id IS NULL']
    },
    { skip: !isUnknownAlbum }
  )

  const [removeSongsFromAlbum] = useRemoveSongsFromAlbumMutation()

  const [openedEditAlbumHeader, { open: openEditAlbumHeader, close: closeEditAlbumHeader }] =
    useDisclosure(false)
  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  function handleDelete() {
    deleteAlbumMutation(album.id)
    navigate(`/albums`, { replace: true })
    toast.success(`${album.title} deleted!`)
  }

  function handleRemoveSongsFromAlbum(songIds: string[]) {
    removeSongsFromAlbum({ id: albumId, songIds })
  }

  if (isLoading || isSongsLoading) return <AlbumLoader />

  return (
    <Stack>
      <HeaderPanelCard
        onEditClick={openEditAlbumHeader}
        menuDropdown={
          <>
            <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditAlbumHeader}>
              Edit
            </Menu.Item>
            <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
              Delete
            </Menu.Item>
          </>
        }
        hideIcons={isUnknownAlbum}
      >
        <Group>
          <AspectRatio>
            <Image
              h={150}
              src={isUnknownAlbum ? unknownPlaceholder : album.imageUrl}
              fallbackSrc={albumPlaceholder}
              radius={'lg'}
              style={(theme) => ({
                boxShadow: theme.shadows.lg
              })}
            />
          </AspectRatio>
          <Stack
            gap={4}
            style={{ ...(!isUnknownAlbum && { alignSelf: 'start', paddingTop: '10px' }) }}
          >
            {!isUnknownAlbum && (
              <Text fw={500} inline>
                Album
              </Text>
            )}
            {isUnknownAlbum ? (
              <Title order={3} fw={200} fs={'italic'}>
                Unknown
              </Title>
            ) : (
              <Title order={1} fw={700}>
                {album.title}
              </Title>
            )}
            <Group gap={4}>
              {album?.artist && (
                <>
                  <Group gap={'xs'}>
                    <Avatar size={35} src={album.artist.imageUrl ?? userPlaceholder} />
                    <Text
                      fw={600}
                      fz={'lg'}
                      sx={{
                        cursor: 'pointer',
                        '&:hover': { textDecoration: 'underline' }
                      }}
                      onClick={handleArtistClick}
                    >
                      {album.artist.name}
                    </Text>
                  </Group>
                  <Text c={'dimmed'}>•</Text>
                </>
              )}
              {album?.releaseDate && (
                <Tooltip
                  label={'Released on ' + dayjs(album.releaseDate).format('DD MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500} c={'dimmed'}>
                    {dayjs(album.releaseDate).format('YYYY')} •
                  </Text>
                </Tooltip>
              )}
              <Text fw={500} c={'dimmed'}>
                {isUnknownAlbum ? songs.totalCount : album.songs.length} songs
              </Text>
            </Group>
          </Stack>
        </Group>
      </HeaderPanelCard>

      <Divider />

      <Card variant={'panel'} h={'100%'} p={0} mx={'xs'}>
        <LoadingOverlay visible={isSongsFetching || isFetching} />

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
            {(isUnknownAlbum ? songs.models : album.songs).map((song) => (
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
      </Card>

      {!isUnknownAlbum && (
        <EditAlbumHeaderModal
          album={album}
          opened={openedEditAlbumHeader}
          onClose={closeEditAlbumHeader}
        />
      )}
      <AddNewAlbumSongModal
        opened={openedAddNewSong}
        onClose={closeAddNewSong}
        albumId={album?.id}
      />
      <AddExistingAlbumSongsModal
        opened={openedAddExistingSongs}
        onClose={closeAddExistingSongs}
        albumId={albumId}
        artistId={album?.artist?.id}
      />
    </Stack>
  )
}

export default Album
