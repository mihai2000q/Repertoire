import {
  ActionIcon,
  Avatar,
  Box,
  Center,
  Divider,
  Grid,
  Group,
  Menu,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useGetAlbumQuery } from '../../../state/api/albumsApi.ts'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import AlbumDrawerLoader from '../loader/AlbumDrawerLoader.tsx'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import { IconDotsVertical, IconEye, IconTrash } from '@tabler/icons-react'
import { useDisclosure } from '@mantine/hooks'
import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import plural from '../../../utils/plural.ts'
import { closeAlbumDrawer } from '../../../state/slice/globalSlice.ts'
import useDynamicDocumentTitle from '../../../hooks/useDynamicDocumentTitle.ts'
import CustomIconMusicNoteEighth from '../../@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'
import AddToPlaylistMenuItem from '../../@ui/menu/item/AddToPlaylistMenuItem.tsx'
import Song from '../../../types/models/Song.ts'
import DeleteAlbumModal from '../../@ui/modal/DeleteAlbumModal.tsx'

function AlbumDrawerSongCard({
  song,
  albumImageUrl,
  onClose
}: {
  song: Song
  albumImageUrl: string | undefined | null
  onClose: () => void
}) {
  const navigate = useNavigate()

  function onClick() {
    onClose()
    navigate(`/song/${song.id}`)
  }

  return (
    <Grid align={'center'} gutter={'xs'} px={'xs'}>
      <Grid.Col span={1}>
        <Text fw={500} ta={'center'}>
          {song.albumTrackNo}
        </Text>
      </Grid.Col>

      <Grid.Col span={1.2}>
        <Avatar
          radius={'md'}
          size={28}
          src={song.imageUrl ?? albumImageUrl}
          alt={(song.imageUrl ?? albumImageUrl) && song.title}
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
        <Text fw={500} truncate={'end'}>
          {song.title}
        </Text>
      </Grid.Col>
    </Grid>
  )
}

function AlbumDrawer() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const setDocumentTitle = useDynamicDocumentTitle()

  const isDocumentTitleSet = useRef(false)

  const { albumId, open: opened } = useAppSelector((state) => state.global.albumDrawer)
  const onClose = () => {
    dispatch(closeAlbumDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
    isDocumentTitleSet.current = false
  }

  const { data: album, isFetching } = useGetAlbumQuery({ id: albumId }, { skip: !albumId })

  useEffect(() => {
    if (album && opened && albumId === album.id && !isDocumentTitleSet.current) {
      setDocumentTitle((prevTitle) => prevTitle + ' - ' + album.title)
      isDocumentTitleSet.current = true
    }
  }, [album, opened])

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleArtistClick() {
    onClose()
    navigate(`/artist/${album.artist.id}`)
  }

  function handleViewDetails() {
    onClose()
    navigate(`/album/${album.id}`)
  }

  if (!album)
    return (
      <RightSideEntityDrawer
        opened={opened}
        onClose={onClose}
        isLoading={true}
        loader={<AlbumDrawerLoader />}
      />
    )

  return (
    <RightSideEntityDrawer
      opened={opened}
      onClose={onClose}
      isLoading={isFetching}
      loader={<AlbumDrawerLoader />}
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
            src={album.imageUrl}
            alt={album.imageUrl && album.title}
            bg={'gray.5'}
            style={{ aspectRatio: 4 / 3 }}
          >
            <Center c={'white'}>
              <CustomIconAlbumVinyl
                aria-label={`default-icon-${album.title}`}
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
                <AddToPlaylistMenuItem
                  ids={[album.id]}
                  type={'album'}
                  closeMenu={closeMenu}
                  disabled={album.songs.length === 0}
                />
                <Menu.Divider />
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
          <Title order={5} fw={700} lh={'xs'} lineClamp={2} fz={'max(1.85vw, 24px)'}>
            {album.title}
          </Title>

          <Group gap={'xxs'} wrap={'nowrap'}>
            {album.artist && (
              <>
                <Group gap={6} wrap={'nowrap'}>
                  <Avatar
                    size={28}
                    src={album.artist.imageUrl}
                    alt={album.artist.imageUrl && album.artist.name}
                    style={(theme) => ({ boxShadow: theme.shadows.sm })}
                    bg={'gray.0'}
                  >
                    <Center c={'gray.7'}>
                      <CustomIconUserAlt
                        aria-label={`default-icon-${album.artist.name}`}
                        size={13}
                      />
                    </Center>
                  </Avatar>
                  <Text
                    fw={700}
                    fz={'lg'}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { textDecoration: 'underline' }
                    }}
                    lh={'xxs'}
                    lineClamp={1}
                    onClick={handleArtistClick}
                  >
                    {album.artist.name}
                  </Text>
                </Group>
                <Text c={'dimmed'}>•</Text>
              </>
            )}
            {album.releaseDate && (
              <>
                <Tooltip
                  label={'Released on ' + dayjs(album.releaseDate).format('D MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500}>{dayjs(album.releaseDate).format('YYYY')}</Text>
                </Tooltip>
                <Text fw={500}>•</Text>
              </>
            )}
            <Text fw={500} c={'dimmed'} truncate={'end'}>
              {album.songs.length} song{plural(album.songs)}
            </Text>
          </Group>

          {album.songs.length > 0 && <Divider my={'xs'} />}

          <Stack gap={'md'}>
            {album.songs.map((song) => (
              <AlbumDrawerSongCard
                key={song.id}
                song={song}
                albumImageUrl={album.imageUrl}
                onClose={onClose}
              />
            ))}
          </Stack>
        </Stack>
      </Stack>

      <DeleteAlbumModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        album={album}
        onDelete={onClose}
      />
    </RightSideEntityDrawer>
  )
}

export default AlbumDrawer
