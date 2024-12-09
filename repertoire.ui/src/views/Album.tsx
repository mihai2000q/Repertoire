import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Card,
  Divider,
  Group,
  Image,
  Menu,
  Space,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useNavigate, useParams } from 'react-router-dom'
import { useDeleteAlbumMutation, useGetAlbumQuery } from '../state/albumsApi.ts'
import AlbumLoader from '../components/album/AlbumLoader.tsx'
import albumPlaceholder from '../assets/image-placeholder-1.jpg'
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

function Album() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const params = useParams()
  const albumId = params['id'] ?? ''

  const [deleteAlbumMutation] = useDeleteAlbumMutation()

  const { data: album, isLoading } = useGetAlbumQuery(albumId)

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

  if (isLoading) return <AlbumLoader />

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
      >
        <Group>
          <AspectRatio>
            <Image
              h={150}
              src={album.imageUrl}
              fallbackSrc={albumPlaceholder}
              radius={'lg'}
              style={(theme) => ({
                boxShadow: theme.shadows.lg
              })}
            />
          </AspectRatio>
          <Stack gap={4} style={{ alignSelf: 'start' }} pt={'xs'}>
            <Text fw={500} inline>
              Album
            </Text>
            <Title order={1} fw={700}>
              {album.title}
            </Title>
            <Group gap={4}>
              {album.artist && (
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
              {album.releaseDate && (
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
                {album.songs.length} songs
              </Text>
            </Group>
          </Stack>
        </Group>
      </HeaderPanelCard>

      <Divider />

      <Card variant={'panel'} h={'100%'} p={0} mx={'xs'}>
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
                <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingSongs}>
                  Add Existing Songs
                </Menu.Item>
                <Menu.Item leftSection={<IconMusicPlus size={15} />} onClick={openAddNewSong}>
                  Add New Song
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>

          <Stack gap={0}>
            {album.songs.map((song) => (
              <AlbumSongCard key={song.id} song={song} />
            ))}
            {album.songs.length === 0 && (
              <NewHorizontalCard onClick={openAddExistingSongs}>Add New Songs</NewHorizontalCard>
            )}
          </Stack>
        </Stack>
      </Card>

      <EditAlbumHeaderModal
        album={album}
        opened={openedEditAlbumHeader}
        onClose={closeEditAlbumHeader}
      />
      <AddNewAlbumSongModal opened={openedAddNewSong} onClose={closeAddNewSong} albumId={albumId} />
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
