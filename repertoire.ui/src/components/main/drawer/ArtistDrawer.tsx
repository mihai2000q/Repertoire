import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Box,
  Divider,
  Group,
  Image,
  Menu,
  SimpleGrid,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { useDeleteArtistMutation, useGetArtistQuery } from '../../../state/artistsApi.ts'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import ArtistDrawerLoader from '../loader/ArtistDrawerLoader.tsx'
import imagePlaceholder from '../../../assets/user-placeholder.jpg'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { useState } from 'react'
import { toast } from 'react-toastify'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import { IconDotsVertical, IconEye, IconTrash } from '@tabler/icons-react'
import plural from '../../../utils/plural.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import { useGetAlbumsQuery } from '../../../state/albumsApi.ts'
import { useGetSongsQuery } from '../../../state/songsApi.ts'
import dayjs from 'dayjs'
import { closeArtistDrawer, deleteArtistDrawer } from '../../../state/globalSlice.ts'

function ArtistDrawer() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const opened = useAppSelector((state) => state.global.artistDrawer.open)
  const artistId = useAppSelector((state) => state.global.artistDrawer.artistId)
  const onClose = () => dispatch(closeArtistDrawer())

  const [deleteArtistMutation] = useDeleteArtistMutation()

  const { data: artist, isFetching } = useGetArtistQuery(artistId, { skip: !artistId })
  const { data: albums, isFetching: isAlbumsFetching } = useGetAlbumsQuery(
    {
      orderBy: ['release_date desc', 'title asc'],
      searchBy: [`artist_id = '${artistId}'`]
    },
    { skip: !artistId }
  )
  const { data: songs, isFetching: isSongsFetching } = useGetSongsQuery(
    {
      orderBy: ['release_date desc', 'title asc'],
      searchBy: [`songs.artist_id = '${artistId}'`]
    },
    { skip: !artistId }
  )

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleViewDetails() {
    onClose()
    navigate(`/artist/${artist.id}`)
  }

  function handleDelete() {
    deleteArtistMutation(artist.id)
    dispatch(deleteArtistDrawer())
    toast.success(`${artist.name} deleted!`)
  }

  if (!artist || !songs || !albums)
    return (
      <RightSideEntityDrawer
        opened={opened}
        onClose={onClose}
        isLoading={true}
        loader={<ArtistDrawerLoader />}
      />
    )

  return (
    <RightSideEntityDrawer
      opened={opened}
      onClose={onClose}
      isLoading={isFetching || isSongsFetching || isAlbumsFetching}
      loader={<ArtistDrawerLoader />}
    >
      <Stack gap={'xs'}>
        <Box
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
          pos={'relative'}
        >
          <AspectRatio ratio={4 / 3}>
            <Image src={artist.imageUrl} fallbackSrc={imagePlaceholder} alt={artist.name} />
          </AspectRatio>

          <Box pos={'absolute'} top={0} right={0} p={7}>
            <Menu opened={isMenuOpened} onChange={setIsMenuOpened}>
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

        <Stack px={'md'} pb={'md'} gap={4}>
          <Title order={5} fw={700} lh={1}>
            {artist.name}
          </Title>

          <Group ml={2} gap={4}>
            <Text fw={500} fz={'sm'} c={'dimmed'}>
              {albums.totalCount} album{plural(albums.totalCount)} • {songs.totalCount} song
              {plural(songs.totalCount)}
            </Text>
          </Group>

          <Stack gap={0} my={6}>
            <Text ml={2} fw={500} fz={'xs'} c={'dimmed'}>
              Albums
            </Text>
            <Divider />
          </Stack>

          <SimpleGrid cols={2} px={'xs'}>
            {albums.models.map((album) => (
              <Group key={album.id} align={'center'} wrap={'nowrap'} gap={'xs'}>
                <Avatar
                  radius={'8px'}
                  size={28}
                  src={album.imageUrl ?? albumPlaceholder}
                  alt={album.title}
                />
                <Stack gap={1} style={{ overflow: 'hidden' }}>
                  <Text fw={500} truncate={'end'} inline>
                    {album.title}
                  </Text>
                  {album.releaseDate && (
                    <Text fw={500} fz={'xs'} c={'dimmed'} inline>
                      {dayjs(album.releaseDate).format('D MMM YYYY')}
                    </Text>
                  )}
                </Stack>
              </Group>
            ))}
          </SimpleGrid>

          <Stack gap={0} my={6}>
            <Text ml={2} fw={500} fz={'xs'} c={'dimmed'}>
              Songs
            </Text>
            <Divider />
          </Stack>

          <SimpleGrid cols={2} px={'xs'}>
            {songs.models.map((song) => (
              <Group key={song.id} align={'center'} gap={'xs'} wrap={'nowrap'}>
                <Avatar
                  radius={'8px'}
                  size={28}
                  src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
                  alt={song.title}
                />
                <Stack gap={1}>
                  <Text fw={500} truncate={'end'} inline>
                    {song.title}
                  </Text>
                  {song.album && (
                    <Text fz={'xxs'} c={'dimmed'} fw={500} truncate={'end'} inline>
                      {song.album.title}
                    </Text>
                  )}
                </Stack>
              </Group>
            ))}
          </SimpleGrid>
        </Stack>
      </Stack>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Artist'}
        description={`Are you sure you want to delete this artist?`}
        onYes={handleDelete}
      />
    </RightSideEntityDrawer>
  )
}

export default ArtistDrawer
