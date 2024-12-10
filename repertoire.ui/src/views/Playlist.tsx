import {
  ActionIcon,
  AspectRatio,
  Button,
  Card,
  Divider,
  Group,
  Image,
  LoadingOverlay,
  Menu,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { useNavigate, useParams } from 'react-router-dom'
import {
  useDeletePlaylistMutation,
  useGetPlaylistQuery,
  useRemoveSongsFromPlaylistMutation
} from '../state/playlistsApi.ts'
import PlaylistLoader from '../components/playlist/PlaylistLoader.tsx'
import playlistPlaceholder from '../assets/image-placeholder-1.jpg'
import PlaylistSongCard from '../components/playlist/PlaylistSongCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import {
  IconCaretDownFilled,
  IconCheck,
  IconDots,
  IconEdit,
  IconPlus,
  IconTrash
} from '@tabler/icons-react'
import NewHorizontalCard from '../components/card/NewHorizontalCard.tsx'
import HeaderPanelCard from '../components/card/HeaderPanelCard.tsx'
import { toast } from 'react-toastify'
import EditPlaylistHeaderModal from '../components/playlist/modal/EditPlaylistHeaderModal.tsx'
import AddExistingPlaylistSongsModal from '../components/playlist/modal/AddExistingPlaylistSongsModal.tsx'
import { useState } from 'react'
import Order from '../types/Order.ts'
import playlistSongsOrders from '../data/playlist/playlistSongsOrders.ts'
import plural from "../utils/plural.ts";

function Playlist() {
  const navigate = useNavigate()

  const params = useParams()
  const playlistId = params['id'] ?? ''

  const [deletePlaylistMutation] = useDeletePlaylistMutation()

  const [songsOrder, setSongsOrder] = useState<Order>(playlistSongsOrders[0])

  const { data: playlist, isLoading, isFetching } = useGetPlaylistQuery(playlistId)

  const [removeSongsFromPlaylist] = useRemoveSongsFromPlaylistMutation()
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  const [
    openedEditPlaylistHeader,
    { open: openEditPlaylistHeader, close: closeEditPlaylistHeader }
  ] = useDisclosure(false)

  function handleDelete() {
    deletePlaylistMutation(playlist.id)
    navigate(`/playlists`, { replace: true })
    toast.success(`${playlist.title} deleted!`)
  }

  function handleRemoveSongsFromPlaylist(songIds: string[]) {
    removeSongsFromPlaylist({ id: playlistId, songIds })
  }

  if (isLoading) return <PlaylistLoader />

  return (
    <Stack>
      <HeaderPanelCard
        onEditClick={openEditPlaylistHeader}
        menuDropdown={
          <>
            <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditPlaylistHeader}>
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
              src={playlist.imageUrl}
              fallbackSrc={playlistPlaceholder}
              radius={'lg'}
              style={(theme) => ({
                boxShadow: theme.shadows.lg
              })}
            />
          </AspectRatio>
          <Stack gap={4} style={{ alignSelf: 'start', paddingTop: '10px' }}>
            <Text fw={500} inline>
              Playlist
            </Text>

            <Title order={1} fw={700}>
              {playlist.title}
            </Title>

            <Text fw={500} fz={'sm'} c={'dimmed'} inline>
              {playlist.songs.length} song{plural(playlist.songs)}
            </Text>

            <Text fz={'sm'} c={'dimmed'} lineClamp={3}>
              {playlist.description}
            </Text>
          </Stack>
        </Group>
      </HeaderPanelCard>

      <Divider />

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
                  {songsOrder.label}
                </Button>
              </Menu.Target>

              <Menu.Dropdown>
                {playlistSongsOrders.map((o) => (
                  <Menu.Item
                    key={o.value}
                    leftSection={songsOrder === o && <IconCheck size={12} />}
                    onClick={() => setSongsOrder(o)}
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
                <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingSongs}>
                  Add Existing Songs
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
              <NewHorizontalCard onClick={openAddExistingSongs}>Add New Song</NewHorizontalCard>
            )}
          </Stack>
        </Stack>
      </Card>

      <EditPlaylistHeaderModal
        playlist={playlist}
        opened={openedEditPlaylistHeader}
        onClose={closeEditPlaylistHeader}
      />
      <AddExistingPlaylistSongsModal
        opened={openedAddExistingSongs}
        onClose={closeAddExistingSongs}
        playlistId={playlistId}
      />
    </Stack>
  )
}

export default Playlist
