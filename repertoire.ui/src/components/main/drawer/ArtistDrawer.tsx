import {
  ActionIcon,
  Avatar,
  Box,
  Center,
  Divider,
  Group,
  Menu,
  SimpleGrid,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { useDeleteArtistMutation, useGetArtistQuery } from '../../../state/api/artistsApi.ts'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import ArtistDrawerLoader from '../loader/ArtistDrawerLoader.tsx'
import { useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { useEffect, useState } from 'react'
import { toast } from 'react-toastify'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import { IconDotsVertical, IconEye, IconTrash, IconUser } from '@tabler/icons-react'
import plural from '../../../utils/plural.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import { useGetAlbumsQuery } from '../../../state/api/albumsApi.ts'
import { useGetSongsQuery } from '../../../state/api/songsApi.ts'
import dayjs from 'dayjs'
import { closeArtistDrawer, deleteArtistDrawer } from '../../../state/slice/globalSlice.ts'
import useDynamicDocumentTitle from '../../../hooks/useDynamicDocumentTitle.ts'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'
import createOrder from '../../../utils/createOrder.ts'
import OrderType from '../../../types/enums/OrderType.ts'
import AlbumProperty from '../../../types/enums/AlbumProperty.ts'
import SongProperty from '../../../types/enums/SongProperty.ts'

function ArtistDrawer() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const setDocumentTitle = useDynamicDocumentTitle()

  const opened = useAppSelector((state) => state.global.artistDrawer.open)
  const artistId = useAppSelector((state) => state.global.artistDrawer.artistId)
  const onClose = () => {
    dispatch(closeArtistDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
  }

  const [deleteArtistMutation, { isLoading: isDeleteLoading }] = useDeleteArtistMutation()

  const { data: artist, isFetching } = useGetArtistQuery(artistId, { skip: !artistId })
  const { data: albums, isFetching: isAlbumsFetching } = useGetAlbumsQuery(
    {
      orderBy: [
        createOrder({
          property: AlbumProperty.ReleaseDate,
          type: OrderType.Descending,
          nullable: true
        }),
        createOrder({ property: AlbumProperty.Title })
      ],
      searchBy: [`artist_id = '${artistId}'`]
    },
    { skip: !artistId }
  )
  const { data: songs, isFetching: isSongsFetching } = useGetSongsQuery(
    {
      orderBy: [
        createOrder({
          property: SongProperty.ReleaseDate,
          type: OrderType.Descending,
          nullable: true
        }),
        createOrder({ property: SongProperty.Title })
      ],
      searchBy: [`songs.artist_id = '${artistId}'`]
    },
    { skip: !artistId }
  )

  useEffect(() => {
    if (artist && opened && !isFetching)
      setDocumentTitle((prevTitle) => prevTitle + ' - ' + artist.name)
  }, [artist, opened, isFetching])

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleViewDetails() {
    onClose()
    navigate(`/artist/${artist.id}`)
  }

  async function handleDelete() {
    await deleteArtistMutation({ id: artist.id }).unwrap()
    dispatch(deleteArtistDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
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
          <Avatar
            w={'100%'}
            h={'unset'}
            radius={0}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            bg={'gray.0'}
            style={{ aspectRatio: 4 / 3 }}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt
                aria-label={`default-icon-${artist.name}`}
                size={'100%'}
                style={{ padding: '26%' }}
              />
            </Center>
          </Avatar>

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

        <Stack px={'md'} pb={'md'} gap={'xxs'}>
          <Title order={5} fw={700} lh={1} lineClamp={2} fz={'max(1.85vw, 24px)'}>
            {artist.name}
          </Title>

          <Group ml={2} gap={'xxs'}>
            <Text fw={500} fz={'sm'} c={'dimmed'}>
              {artist.isBand
                ? artist.bandMembers.length + ` member${plural(artist.bandMembers)} • `
                : ''}
              {albums.totalCount} album{plural(albums.totalCount)} • {songs.totalCount} song
              {plural(songs.totalCount)}
            </Text>
          </Group>

          {artist.isBand && artist.bandMembers.length > 0 && (
            <Stack gap={0} my={6}>
              <Text ml={2} fw={500} fz={'xs'} c={'dimmed'}>
                Band Members
              </Text>
              <Divider />
            </Stack>
          )}

          <Group align={'start'} px={6} gap={'sm'}>
            {artist.isBand &&
              artist.bandMembers.map((bandMember) => (
                <Stack key={bandMember.id} align={'center'} gap={'xxs'} w={53}>
                  <Avatar
                    variant={'light'}
                    size={42}
                    color={bandMember.color}
                    src={bandMember.imageUrl}
                    alt={bandMember.imageUrl && bandMember.name}
                    style={(theme) => ({ boxShadow: theme.shadows.sm })}
                  >
                    <IconUser aria-label={`icon-${bandMember.name}`} size={19} />
                  </Avatar>

                  <Text ta={'center'} fw={500} fz={'sm'} lh={1.1} lineClamp={2}>
                    {bandMember.name}
                  </Text>
                </Stack>
              ))}
          </Group>

          {albums.totalCount > 0 && (
            <Stack gap={0} my={6}>
              <Text ml={2} fw={500} fz={'xs'} c={'dimmed'}>
                Albums
              </Text>
              <Divider />
            </Stack>
          )}

          <SimpleGrid cols={2} px={'xs'}>
            {albums.models.map((album) => (
              <Group key={album.id} wrap={'nowrap'} gap={'xs'}>
                <Avatar
                  radius={'md'}
                  size={28}
                  src={album.imageUrl}
                  alt={album.imageUrl && album.title}
                  bg={'gray.5'}
                  style={(theme) => ({ boxShadow: theme.shadows.sm })}
                >
                  <Center c={'white'}>
                    <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={13} />
                  </Center>
                </Avatar>
                <Stack gap={1} style={{ overflow: 'hidden' }}>
                  <Text fw={500} truncate={'end'} lh={'xxs'}>
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

          {songs.totalCount > 0 && (
            <Stack gap={0} my={6}>
              <Text ml={2} fw={500} fz={'xs'} c={'dimmed'}>
                Songs
              </Text>
              <Divider />
            </Stack>
          )}

          <SimpleGrid cols={2} px={'xs'}>
            {songs.models.map((song) => (
              <Group key={song.id} gap={'xs'} wrap={'nowrap'}>
                <Avatar
                  radius={'md'}
                  size={28}
                  src={song.imageUrl ?? song.album?.imageUrl}
                  alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
                  bg={'gray.5'}
                  style={(theme) => ({ boxShadow: theme.shadows.sm })}
                >
                  <Center c={'white'}>
                    <CustomIconMusicNoteEighth
                      aria-label={`default-icon-${song.title}`}
                      size={16}
                    />
                  </Center>
                </Avatar>
                <Stack gap={1} style={{ overflow: 'hidden' }}>
                  <Text fw={500} truncate={'end'} lh={'xxs'}>
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
        isLoading={isDeleteLoading}
      />
    </RightSideEntityDrawer>
  )
}

export default ArtistDrawer
