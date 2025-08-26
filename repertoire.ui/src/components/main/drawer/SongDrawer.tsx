import {
  ActionIcon,
  Anchor,
  Avatar,
  Box,
  Center,
  Divider,
  Grid,
  Group,
  HoverCard,
  Menu,
  NumberFormatter,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useDeleteSongMutation, useGetSongQuery } from '../../../state/api/songsApi.ts'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import SongDrawerLoader from '../loader/SongDrawerLoader.tsx'
import {
  IconBrandYoutubeFilled,
  IconCheck,
  IconDotsVertical,
  IconEye,
  IconGuitarPickFilled,
  IconTrash
} from '@tabler/icons-react'
import dayjs from 'dayjs'
import { useDisclosure } from '@mantine/hooks'
import { useEffect, useRef, useState } from 'react'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import { toast } from 'react-toastify'
import { useNavigate } from 'react-router-dom'
import RightSideEntityDrawer from '../../@ui/drawer/RightSideEntityDrawer.tsx'
import { closeSongDrawer } from '../../../state/slice/globalSlice.ts'
import DifficultyBar from '../../@ui/bar/DifficultyBar.tsx'
import YoutubeModal from '../../@ui/modal/YoutubeModal.tsx'
import useDynamicDocumentTitle from '../../../hooks/useDynamicDocumentTitle.ts'
import ConfidenceBar from '../../@ui/bar/ConfidenceBar.tsx'
import ProgressBar from '../../@ui/bar/ProgressBar.tsx'
import PerfectRehearsalMenuItem from '../../@ui/menu/item/PerfectRehearsalMenuItem.tsx'
import PartialRehearsalMenuItem from '../../@ui/menu/item/song/PartialRehearsalMenuItem.tsx'
import CustomIconMusicNote from '../../@ui/icons/CustomIconMusicNote.tsx'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'
import AddToPlaylistMenuItem from '../../@ui/menu/item/AddToPlaylistMenuItem.tsx'

const firstColumnSize = 4
const secondColumnSize = 8

function SongDrawer() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const setDocumentTitle = useDynamicDocumentTitle()

  const isDocumentTitleSet = useRef(false)

  const { songId, open: opened } = useAppSelector((state) => state.global.songDrawer)
  const onClose = () => {
    dispatch(closeSongDrawer())
    setDocumentTitle((prevTitle) => prevTitle.split(' - ')[0])
    isDocumentTitleSet.current = false
  }

  const { data: song, isFetching } = useGetSongQuery(songId, { skip: !songId })
  const [deleteSongMutation, { isLoading: isDeleteLoading }] = useDeleteSongMutation()

  useEffect(() => {
    if (song && opened && songId === song.id && !isDocumentTitleSet.current) {
      setDocumentTitle((prevTitle) => prevTitle + ' - ' + song.title)
      isDocumentTitleSet.current = true
    }
  }, [song, opened])

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const showInfo =
    song &&
    (song.difficulty ||
      song.guitarTuning ||
      song.bpm ||
      song.isRecorded ||
      song.lastTimePlayed ||
      song.rehearsals != 0 ||
      song.confidence != 0 ||
      song.progress != 0)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)
  const [openedYoutube, { open: openYoutube, close: closeYoutube }] = useDisclosure(false)

  function handleArtistClick() {
    onClose()
    navigate(`/artist/${song.artist.id}`)
  }

  function handleAlbumClick() {
    onClose()
    navigate(`/album/${song.album.id}`)
  }

  function handleViewDetails() {
    onClose()
    navigate(`/song/${songId}`)
  }

  async function handleDelete() {
    await deleteSongMutation(song.id).unwrap()
    onClose()
    toast.success(`${song.title} deleted!`)
  }

  if (!song)
    return (
      <RightSideEntityDrawer
        opened={opened}
        onClose={onClose}
        isLoading={true}
        loader={<SongDrawerLoader />}
      />
    )

  return (
    <RightSideEntityDrawer
      opened={opened}
      onClose={onClose}
      isLoading={isFetching}
      loader={<SongDrawerLoader />}
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
            src={song.imageUrl ?? song.album?.imageUrl}
            alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
            bg={'gray.5'}
            style={{ aspectRatio: 4 / 3 }}
          >
            <Center c={'white'}>
              <CustomIconMusicNote
                aria-label={`default-icon-${song.title}`}
                size={'100%'}
                style={{ padding: '35%' }}
              />
            </Center>
          </Avatar>

          <Box pos={'absolute'} top={0} right={0} p={7}>
            <Menu opened={isMenuOpened} onOpen={openMenu} onClose={closeMenu} zIndex={3000}>
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

                <Menu.Divider />
                <AddToPlaylistMenuItem ids={[song.id]} type={'song'} closeMenu={closeMenu} />
                <PartialRehearsalMenuItem songId={song.id} closeMenu={closeMenu} />
                <PerfectRehearsalMenuItem id={song.id} closeMenu={closeMenu} type={'song'} />
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

        <Stack px={'md'} pb={'xs'} gap={'xxs'}>
          <Title order={5} fw={700} lh={'xs'} lineClamp={2} fz={'max(1.85vw, 24px)'}>
            {song.title}
          </Title>

          <Group gap={'xxs'} wrap={'nowrap'}>
            {song.artist && (
              <Group gap={'xs'} wrap={'nowrap'}>
                <Avatar
                  size={28}
                  src={song.artist.imageUrl}
                  alt={song.artist.imageUrl && song.artist.name}
                  style={(theme) => ({ boxShadow: theme.shadows.sm })}
                  bg={'gray.0'}
                >
                  <Center c={'gray.7'}>
                    <CustomIconUserAlt aria-label={`default-icon-${song.artist.name}`} size={13} />
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
                  {song.artist.name}
                </Text>
              </Group>
            )}

            {song.album && (
              <Group gap={0} wrap={'nowrap'}>
                {song.artist && (
                  <Text fw={500} c={'dimmed'} lh={'xxs'} pr={4}>
                    on
                  </Text>
                )}
                <HoverCard>
                  <HoverCard.Target>
                    <Text
                      fw={600}
                      lh={'xxs'}
                      c={'dark'}
                      sx={{
                        cursor: 'pointer',
                        '&:hover': { textDecoration: 'underline' }
                      }}
                      lineClamp={1}
                      onClick={handleAlbumClick}
                    >
                      {song.album.title}
                    </Text>
                  </HoverCard.Target>
                  <HoverCard.Dropdown maw={300}>
                    <Group gap={'xs'} wrap={'nowrap'}>
                      <Avatar
                        radius={'md'}
                        size={45}
                        src={song.album.imageUrl}
                        alt={song.album.imageUrl && song.album.title}
                        bg={'gray.5'}
                      >
                        <Center c={'white'}>
                          <CustomIconAlbumVinyl
                            aria-label={`default-icon-${song.album.title}`}
                            size={18}
                          />
                        </Center>
                      </Avatar>
                      <Stack gap={2}>
                        <Text fw={500} fz={'xs'} inline>
                          Album
                        </Text>
                        <Text fw={600} fz={'md'} lh={'xxs'} lineClamp={2}>
                          {song.album.title}
                        </Text>
                        {song.album.releaseDate && (
                          <Text fw={500} c={'dimmed'} fz={'sm'} inline>
                            {dayjs(song.album.releaseDate).format('D MMM YYYY')}
                          </Text>
                        )}
                      </Stack>
                    </Group>
                  </HoverCard.Dropdown>
                </HoverCard>
              </Group>
            )}

            {song.releaseDate && (
              <>
                {(song.album || song.artist) && (
                  <Text c={'dimmed'} inline>
                    •
                  </Text>
                )}
                <Tooltip
                  label={'Released on ' + dayjs(song.releaseDate).format('D MMMM YYYY')}
                  openDelay={200}
                  position={'bottom'}
                >
                  <Text fw={500} c={'dimmed'} inline>
                    {dayjs(song.releaseDate).format('YYYY')}
                  </Text>
                </Tooltip>
              </>
            )}
          </Group>

          {song.description !== '' && (
            <Text size="sm" c="dimmed" mt={'xs'} px={'xs'} lineClamp={3}>
              {song.description}
            </Text>
          )}

          {showInfo && <Divider mt={'xs'} />}

          <Grid align={'center'} gutter={'sm'} p={showInfo ? 'xs' : 0}>
            {song.difficulty && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Difficulty:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <DifficultyBar difficulty={song.difficulty} size={7} />
                </Grid.Col>
              </>
            )}

            {song.guitarTuning && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'} truncate={'end'}>
                    Guitar Tuning:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <Text fw={600}>{song.guitarTuning.name}</Text>
                </Grid.Col>
              </>
            )}

            {song.bpm && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Bpm:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <Text fw={600}>{song.bpm}</Text>
                </Grid.Col>
              </>
            )}

            {song.isRecorded && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Recorded:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <ActionIcon
                    component={'div'}
                    size={'20px'}
                    aria-label={'recorded-icon'}
                    sx={(theme) => ({
                      cursor: 'default',
                      backgroundColor: theme.colors.primary[5],
                      '&:hover': { backgroundColor: theme.colors.primary[5] },
                      '&:active': { transform: 'none' }
                    })}
                  >
                    <IconCheck size={14} />
                  </ActionIcon>
                </Grid.Col>
              </>
            )}

            {song.lastTimePlayed && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Last Played On:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <Tooltip
                    label={`Last time played on ${dayjs(song.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
                    openDelay={400}
                    disabled={!song.lastTimePlayed}
                  >
                    <Text fw={600}>{dayjs(song.lastTimePlayed).format('D MMMM YYYY')}</Text>
                  </Tooltip>
                </Grid.Col>
              </>
            )}

            {song.rehearsals !== 0 && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Rehearsals:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <Text fw={600}>
                    <NumberFormatter value={song.rehearsals} />
                  </Text>
                </Grid.Col>
              </>
            )}

            {song.confidence !== 0 && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'} truncate={'end'}>
                    Confidence:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <ConfidenceBar confidence={song.confidence} size={7} />
                </Grid.Col>
              </>
            )}

            {song.progress !== 0 && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'} truncate={'end'}>
                    Progress:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
                  <ProgressBar progress={song.progress} size={7} />
                </Grid.Col>
              </>
            )}
          </Grid>

          {(song.youtubeLink || song.songsterrLink) && <Divider my={4} />}

          <Group gap={2} style={{ alignSelf: 'end' }}>
            <Tooltip.Group openDelay={200}>
              {song.songsterrLink && (
                <Tooltip label={'Open Songsterr'}>
                  <Anchor
                    aria-label={'songsterr'}
                    underline={'never'}
                    href={song.songsterrLink}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <ActionIcon variant={'transparent'} c={'blue.7'} aria-label={'songsterr'}>
                      <IconGuitarPickFilled size={23} />
                    </ActionIcon>
                  </Anchor>
                </Tooltip>
              )}

              {song.youtubeLink && (
                <Tooltip label={'Open Youtube'}>
                  <ActionIcon
                    mb={3}
                    variant={'transparent'}
                    c={'red.7'}
                    aria-label={'youtube'}
                    onClick={openYoutube}
                  >
                    <IconBrandYoutubeFilled size={25} />
                  </ActionIcon>
                </Tooltip>
              )}
            </Tooltip.Group>
          </Group>
        </Stack>
      </Stack>

      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={'Delete Song'}
        description={`Are you sure you want to delete this song?`}
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
      <YoutubeModal
        title={song.title}
        link={song.youtubeLink}
        opened={openedYoutube}
        onClose={closeYoutube}
      />
    </RightSideEntityDrawer>
  )
}

export default SongDrawer
