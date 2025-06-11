import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import useDynamicDocumentTitle from '../../../hooks/useDynamicDocumentTitle.ts'
import { closePlaylistDrawer, deletePlaylistDrawer } from '../../../state/slice/globalSlice.ts'
import { useEffect, useState } from 'react'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import {
  ActionIcon,
  Avatar,
  Box,
  Center,
  Divider,
  Grid,
  Menu,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { IconDotsVertical, IconEye, IconPlaylist, IconTrash } from '@tabler/icons-react'
import plural from '../../../utils/plural.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import Song from '../../../types/models/Song.ts'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'
import { useDeletePlaylistMutation, useGetPlaylistQuery } from '../../../state/api/playlistsApi.ts'
import PlaylistDrawerLoader from '../loader/PlaylistDrawerLoader.tsx'

function PlaylistDrawerSongCard({ song, onClose }: { song: Song; onClose: () => void }) {
  const navigate = useNavigate()

  function onClick() {
    onClose()
    navigate(`/song/${song.id}`)
  }

  return (
    <Grid align={'center'} gutter={'xs'} px={'xs'}>
      <Grid.Col span={1}>
        <Text fw={500} ta={'center'}>
          {song.playlistTrackNo}
        </Text>
      </Grid.Col>

      <Grid.Col span={1.2}>
        <Avatar
          radius={'md'}
          size={28}
          src={song.imageUrl ?? song.album?.imageUrl}
          alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
          bg={'gray.5'}
          sx={(theme) => ({
            transition: '0.18s',
            cursor: 'pointer',
            boxShadow: theme.shadows.sm,
            '&:hover': {
              transform: 'scale(1.2)'
            }
          })}
          onClick={onClick}
        >
          <Center c={'white'}>
            <CustomIconMusicNoteEighth aria-label={`default-icon-${song.title}`} size={16} />
          </Center>
        </Avatar>
      </Grid.Col>

      <Grid.Col span={9.6}>
        <Stack gap={0} style={{ overflow: 'hidden' }}>
          <Text fw={500} truncate={'end'}>
            {song.title}
          </Text>
          {song.artist && (
            <Text fz={'xs'} c={'dimmed'} fw={500} lh={'xxs'} truncate={'end'}>
              {song.artist.name}
            </Text>
          )}
        </Stack>
      </Grid.Col>
    </Grid>
  )
}

function PlaylistDrawer() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const setDocumentTitle = useDynamicDocumentTitle()

  const { playlistId, open: opened } = useAppSelector((state) => state.global.playlistDrawer)
  const onClose = () => {
    dispatch(closePlaylistDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
  }

  const [deletePlaylistMutation, { isLoading: isDeleteLoading }] = useDeletePlaylistMutation()

  const { data: playlist, isFetching } = useGetPlaylistQuery(
    { id: playlistId },
    { skip: !playlistId }
  )

  useEffect(() => {
    if (playlist && opened && !isFetching)
      setDocumentTitle((prevTitle) => prevTitle + ' - ' + playlist.title)
  }, [playlist, opened, isFetching])

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleViewDetails() {
    onClose()
    navigate(`/playlist/${playlist.id}`)
  }

  async function handleDelete() {
    await deletePlaylistMutation(playlist.id).unwrap()
    dispatch(deletePlaylistDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
    toast.success(`${playlist.title} deleted!`)
  }

  if (!playlist)
    return (
      <RightSideEntityDrawer
        opened={opened}
        onClose={onClose}
        isLoading={true}
        loader={<PlaylistDrawerLoader />}
      />
    )

  return (
    <RightSideEntityDrawer
      opened={opened}
      onClose={onClose}
      isLoading={isFetching}
      loader={<PlaylistDrawerLoader />}
    >
      <Stack gap={'xs'}>
        <Box
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
          pos={'relative'}
        >
          <Avatar
            radius={0}
            w={'100%'}
            h={'unset'}
            src={playlist.imageUrl}
            alt={playlist.imageUrl && playlist.title}
            bg={'gray.5'}
            style={{ aspectRatio: 4 / 3 }}
          >
            <Center c={'white'}>
              <IconPlaylist
                aria-label={`default-icon-${playlist.title}`}
                size={'100%'}
                style={{ padding: '35%' }}
              />
            </Center>
          </Avatar>

          <Box pos={'absolute'} top={0} right={0} p={7}>
            <Menu opened={isMenuOpened} onOpen={openMenu} onClose={closeMenu}>
              <Menu.Target>
                <ActionIcon
                  variant={'grey-subtle'}
                  aria-label={'more-menu'}
                  style={{ transition: '0.25s', opacity: isHovered || isMenuOpened ? 1 : 0 }}
                >
                  <IconDotsVertical size={20} />
                </ActionIcon>
              </Menu.Target>

              <Menu.Dropdown>
                <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
                  View Details
                </Menu.Item>
                <Menu.Item
                  leftSection={<IconTrash size={14} />}
                  c={'red.5'}
                  onClick={openDeleteWarning}
                >
                  Delete
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Box>
        </Box>

        <Stack px={'md'} pb={'md'} gap={'xxs'}>
          <Title order={5} fw={700} lineClamp={2} fz={'max(1.85vw, 24px)'}>
            {playlist.title}
          </Title>

          <Text fw={500} c={'dimmed'} lh={'xs'} truncate={'end'}>
            {playlist.songs.length} song{plural(playlist.songs)}
          </Text>

          {playlist.description !== '' && (
            <Text size="sm" c="dimmed" my={'xxs'} px={'xs'} lineClamp={3}>
              {playlist.description}
            </Text>
          )}

          {playlist.songs.length > 0 && (
            <Divider mt={playlist.description === '' ? 'xs' : 'xxs'} mb={'xs'} />
          )}

          <Stack gap={'md'}>
            {playlist.songs.map((song) => (
              <PlaylistDrawerSongCard key={song.playlistSongId} song={song} onClose={onClose} />
            ))}
          </Stack>
        </Stack>
      </Stack>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Playlist'}
        description={`Are you sure you want to delete this playlist?`}
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </RightSideEntityDrawer>
  )
}

export default PlaylistDrawer
