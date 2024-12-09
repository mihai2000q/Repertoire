import {
  ActionIcon,
  Anchor,
  AspectRatio,
  Avatar,
  Button,
  Divider,
  Grid,
  Group,
  HoverCard,
  Image,
  Menu,
  Progress,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import songPlaceholder from '../assets/image-placeholder-1.jpg'
import userPlaceholder from '../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../state/store.ts'
import { useNavigate, useParams } from 'react-router-dom'
import { openAlbumDrawer, openArtistDrawer } from '../state/globalSlice.ts'
import SongLoader from '../components/song/SongLoader.tsx'
import { useDeleteSongMutation, useGetSongQuery } from '../state/songsApi.ts'
import useDifficultyInfo from '../hooks/useDifficultyInfo.ts'
import Difficulty from '../utils/enums/Difficulty.ts'
import {
  IconBrandYoutube,
  IconCheck,
  IconEdit,
  IconGuitarPick,
  IconTrash
} from '@tabler/icons-react'
import SongSections from '../components/song/SongSections.tsx'
import EditPanelCard from '../components/card/EditPanelCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import EditSongDescriptionModal from '../components/song/modal/EditSongDescriptionModal.tsx'
import EditSongInformationModal from '../components/song/modal/EditSongInformationModal.tsx'
import EditSongLinksModal from '../components/song/modal/EditSongLinksModal.tsx'
import EditSongHeaderModal from '../components/song/modal/EditSongHeaderModal.tsx'
import HeaderPanelCard from '../components/card/HeaderPanelCard.tsx'
import { toast } from 'react-toastify'

const NotSet = () => (
  <Text fz={'sm'} c={'dimmed'} fs={'oblique'} inline>
    not set
  </Text>
)

function Song() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const params = useParams()
  const songId = params['id'] ?? ''

  const [deleteSongMutation] = useDeleteSongMutation()

  const { data: song, isLoading } = useGetSongQuery(songId)

  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(song?.difficulty)

  const [openedEditSongHeader, { open: openEditSongHeader, close: closeEditSongHeader }] =
    useDisclosure(false)
  const [
    openedEditSongInformation,
    { open: openEditSongInformation, close: closeEditSongInformation }
  ] = useDisclosure(false)
  const [
    openedEditSongDescription,
    { open: openEditSongDescription, close: closeEditSongDescription }
  ] = useDisclosure(false)
  const [openedEditSongLinks, { open: openEditSongLinks, close: closeEditSongLinks }] =
    useDisclosure(false)

  function handleAlbumClick() {
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(song.artist.id))
  }

  function handleDelete() {
    deleteSongMutation(song.id)
    navigate(`/songs`, { replace: true })
    toast.success(`${song.title} deleted!`)
  }

  if (isLoading) return <SongLoader />

  return (
    <Stack>
      <HeaderPanelCard
        onEditClick={openEditSongHeader}
        menuDropdown={
          <>
            <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditSongHeader}>
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
              src={song.imageUrl}
              fallbackSrc={songPlaceholder}
              radius={'lg'}
              style={(theme) => ({
                boxShadow: theme.shadows.lg
              })}
            />
          </AspectRatio>
          <Stack gap={4} style={{ alignSelf: 'start' }} pt={'xs'}>
            <Text fw={500} inline>
              Song
            </Text>
            <Title order={1} fw={700}>
              {song.title}
            </Title>
            <Group gap={4}>
              {song.artist && (
                <Group gap={'xs'}>
                  <Avatar size={35} src={song.artist.imageUrl ?? userPlaceholder} />
                  <Text
                    fw={700}
                    fz={'lg'}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { textDecoration: 'underline' }
                    }}
                    onClick={handleArtistClick}
                  >
                    {song.artist.name}
                  </Text>
                </Group>
              )}

              {song.album && (
                <Group gap={0}>
                  {song.artist && (
                    <Text fw={500} c={'dimmed'} inline pr={4}>
                      on
                    </Text>
                  )}
                  <HoverCard shadow={'lg'} withArrow>
                    <HoverCard.Target>
                      <Text
                        fw={600}
                        inline
                        c={'dark'}
                        sx={{
                          cursor: 'pointer',
                          '&:hover': { textDecoration: 'underline' }
                        }}
                        onClick={handleAlbumClick}
                      >
                        {song.album.title}
                      </Text>
                    </HoverCard.Target>
                    <HoverCard.Dropdown maw={300}>
                      <Group align={'center'} gap={'xs'} wrap={'nowrap'}>
                        <Avatar
                          size={45}
                          radius={'md'}
                          src={song.album.imageUrl ?? songPlaceholder}
                        />
                        <Stack gap={2}>
                          <Text fw={500} fz={'xs'} inline>
                            Album
                          </Text>
                          <Text fw={600} fz={'md'} inline lineClamp={2}>
                            {song.album.title}
                          </Text>
                          {song.album.releaseDate && (
                            <Text fw={500} c={'dimmed'} fz={'sm'} inline>
                              {dayjs(song.album.releaseDate).format('MMM YYYY')}
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
                  {(song.album || song.artist) && <Text c={'dimmed'}>•</Text>}
                  <Tooltip
                    label={'Released on ' + dayjs(song.releaseDate).format('DD MMMM YYYY')}
                    openDelay={200}
                    position={'bottom'}
                  >
                    <Text fw={500} c={'dimmed'}>
                      {dayjs(song.releaseDate).format('YYYY')}
                    </Text>
                  </Tooltip>
                </>
              )}
            </Group>
          </Stack>
        </Group>
      </HeaderPanelCard>

      <Divider />

      <Group align="start" mb={'lg'}>
        <Stack flex={1}>
          <EditPanelCard p={'md'} onEditClick={openEditSongInformation}>
            <Stack gap={'xs'}>
              <Text fw={600}>Information</Text>
              <Grid align={'center'} gutter={'sm'}>
                <Grid.Col span={6}>
                  <Text fw={500} c={'dimmed'}>
                    Difficulty:
                  </Text>
                </Grid.Col>
                <Grid.Col span={6}>
                  {song.difficulty ? (
                    <Tooltip
                      label={`This song is ${song.difficulty}`}
                      openDelay={400}
                      position={'top'}
                    >
                      <Group grow gap={4}>
                        {Array.from(Array(Object.entries(Difficulty).length)).map((_, i) => (
                          <Progress
                            key={i}
                            size={5}
                            maw={40}
                            value={i + 1 <= difficultyNumber ? 100 : 0}
                            color={difficultyColor}
                          />
                        ))}
                      </Group>
                    </Tooltip>
                  ) : (
                    <NotSet />
                  )}
                </Grid.Col>

                <Grid.Col span={6}>
                  <Text fw={500} c={'dimmed'} truncate={'end'}>
                    Guitar Tuning:
                  </Text>
                </Grid.Col>
                <Grid.Col span={6}>
                  {song.guitarTuning ? <Text fw={600}>{song.guitarTuning.name}</Text> : <NotSet />}
                </Grid.Col>

                <Grid.Col span={6}>
                  <Tooltip label={'Beats Per Minute'} openDelay={200} position={'top-start'}>
                    <Text fw={500} c={'dimmed'}>
                      Bpm:
                    </Text>
                  </Tooltip>
                </Grid.Col>
                <Grid.Col span={6}>
                  {song.bpm ? <Text fw={600}>{song.bpm}</Text> : <NotSet />}
                </Grid.Col>

                <Grid.Col span={6}>
                  <Text fw={500} c={'dimmed'}>
                    Recorded:
                  </Text>
                </Grid.Col>
                <Grid.Col span={6}>
                  {song.isRecorded ? (
                    <ActionIcon
                      size={'sm'}
                      sx={(theme) => ({
                        cursor: 'default',
                        backgroundColor: theme.colors.cyan[5],
                        '&:hover': { backgroundColor: theme.colors.cyan[5] }
                      })}
                    >
                      <IconCheck size={17} />
                    </ActionIcon>
                  ) : (
                    <Text fw={600}>No</Text>
                  )}
                </Grid.Col>
              </Grid>
            </Stack>
          </EditPanelCard>

          <EditPanelCard p={'md'} onEditClick={openEditSongLinks}>
            <Stack>
              <Text fw={600}>Links</Text>
              <Stack gap={'xs'}>
                {!song.youtubeLink && !song.songsterrLink && (
                  <Text fw={500} c={'dimmed'} ta={'center'}>
                    No links to display
                  </Text>
                )}
                {song.youtubeLink && (
                  <Anchor
                    underline={'never'}
                    href={song.youtubeLink}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <Button
                      fullWidth
                      variant={'gradient'}
                      size={'md'}
                      radius={'lg'}
                      leftSection={<IconBrandYoutube size={30} />}
                      fz={'h6'}
                      fw={700}
                      gradient={{ from: 'red.7', to: 'red.1', deg: 90 }}
                    >
                      Youtube
                    </Button>
                  </Anchor>
                )}
                {song.songsterrLink && (
                  <Anchor
                    underline={'never'}
                    href={song.songsterrLink}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <Button
                      fullWidth
                      variant={'gradient'}
                      size={'md'}
                      radius={'lg'}
                      leftSection={<IconGuitarPick size={30} />}
                      fz={'h6'}
                      fw={700}
                      gradient={{ from: 'blue.7', to: 'blue.1', deg: 90 }}
                    >
                      Songsterr
                    </Button>
                  </Anchor>
                )}
              </Stack>
            </Stack>
          </EditPanelCard>
        </Stack>

        <Stack flex={1.75}>
          <EditPanelCard p={'md'} onEditClick={openEditSongDescription}>
            <Stack gap={'xs'}>
              <Text fw={600}>Description</Text>
              {song.description ? (
                <Text>{song.description}</Text>
              ) : (
                <Text fw={500} c={'dimmed'}>
                  No Description
                </Text>
              )}
            </Stack>
          </EditPanelCard>

          <SongSections songId={songId} sections={song.sections} />
        </Stack>
      </Group>

      <EditSongHeaderModal
        song={song}
        opened={openedEditSongHeader}
        onClose={closeEditSongHeader}
      />
      <EditSongInformationModal
        song={song}
        opened={openedEditSongInformation}
        onClose={closeEditSongInformation}
      />
      <EditSongDescriptionModal
        song={song}
        opened={openedEditSongDescription}
        onClose={closeEditSongDescription}
      />
      <EditSongLinksModal song={song} opened={openedEditSongLinks} onClose={closeEditSongLinks} />
    </Stack>
  )
}

export default Song
