import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Box,
  Divider,
  Grid,
  Group,
  HoverCard,
  Image,
  Menu,
  NumberFormatter,
  Progress,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import { useDeleteSongMutation, useGetSongQuery } from '../../state/songsApi.ts'
import { useAppSelector } from '../../state/store.ts'
import SongDrawerLoader from './loader/SongDrawerLoader.tsx'
import imagePlaceholder from '../../assets/image-placeholder-1.jpg'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import Difficulty from '../../utils/enums/Difficulty.ts'
import {
  IconBrandYoutubeFilled,
  IconCheck,
  IconDotsVertical,
  IconEye,
  IconGuitarPickFilled,
  IconTrash
} from '@tabler/icons-react'
import dayjs from 'dayjs'
import useDifficultyInfo from '../../hooks/useDifficultyInfo.ts'
import { useDisclosure } from '@mantine/hooks'
import { useState } from 'react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { toast } from 'react-toastify'
import { useNavigate } from 'react-router-dom'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import RightSideEntityDrawer from '../@ui/drawer/RightSideEntityDrawer.tsx'

const firstColumnSize = 4
const secondColumnSize = 8

interface SongDrawerProps {
  opened: boolean
  onClose: () => void
}

function SongDrawer({ opened, onClose }: SongDrawerProps) {
  const navigate = useNavigate()

  const songId = useAppSelector((state) => state.global.songDrawer.songId)

  const { data: song, isFetching } = useGetSongQuery(songId, { skip: !songId })
  const [deleteSongMutation] = useDeleteSongMutation()

  const [isHovered, setIsHovered] = useState(false)
  const [isMenuOpened, setIsMenuOpened] = useState(false)

  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(song?.difficulty)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleViewDetails() {
    onClose()
    navigate(`/song/${songId}`)
  }

  function handleDelete() {
    deleteSongMutation(song.id)
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
          <AspectRatio ratio={4 / 3}>
            <Image src={song.imageUrl} fallbackSrc={imagePlaceholder} alt={song.title} />
          </AspectRatio>

          <Box pos={'absolute'} top={0} right={0} p={7}>
            <Menu opened={isMenuOpened} onChange={setIsMenuOpened}>
              <Menu.Target>
                <ActionIcon
                  variant={'grey-subtle'}
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

        <Stack px={'md'} pb={'xs'} gap={4}>
          <Title order={5} fw={700}>
            {song.title}
          </Title>

          <Group>
            <Group gap={4}>
              {song.artist && (
                <Group gap={'xs'}>
                  <Avatar size={28} src={song.artist.imageUrl ?? userPlaceholder} />
                  <Text
                    fw={700}
                    fz={'lg'}
                    sx={{
                      cursor: 'pointer',
                      '&:hover': { textDecoration: 'underline' }
                    }}
                    inline
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
                              {dayjs(song.album.releaseDate).format('DD MMM YYYY')}
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
                    label={'Released on ' + dayjs(song.releaseDate).format('DD MMMM YYYY')}
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
          </Group>

          <Text size="sm" c="dimmed" my={'xs'} px={'xs'}>
            {song.description}
          </Text>

          {(song.difficulty ||
            song.guitarTuning ||
            song.bpm ||
            song.isRecorded ||
            song.lastTimePlayed ||
            song.rehearsals != 0 ||
            song.confidence != 0 ||
            song.progress != 0) && <Divider />}

          <Grid align={'center'} gutter={'sm'} p={'xs'}>
            {song.difficulty && (
              <>
                <Grid.Col span={firstColumnSize}>
                  <Text fw={500} c={'dimmed'}>
                    Difficulty:
                  </Text>
                </Grid.Col>
                <Grid.Col span={secondColumnSize}>
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
                    size={'20px'}
                    sx={(theme) => ({
                      cursor: 'default',
                      backgroundColor: theme.colors.cyan[5],
                      '&:hover': { backgroundColor: theme.colors.cyan[5] }
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
                  <Text fw={600}>{dayjs(song.lastTimePlayed).format('DD MMM YYYY')}</Text>
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
                  <Tooltip.Floating label={<><NumberFormatter value={song.confidence} />%</>}>
                    <Progress flex={1} size={7} value={song.confidence} />
                  </Tooltip.Floating>
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
                  <Tooltip.Floating label={<NumberFormatter value={song.progress} />}>
                    <Progress flex={1} size={7} value={song.progress / 10} color={'green'} />
                  </Tooltip.Floating>
                </Grid.Col>
              </>
            )}
          </Grid>

          {(song.youtubeLink || song.songsterrLink) && <Divider my={4} />}

          <Group gap={2} style={{ alignSelf: 'end' }}>
            <Tooltip.Group openDelay={200}>
              {song.songsterrLink && (
                <Tooltip label={'Open Songsterr'}>
                  <ActionIcon variant={'transparent'} c={'blue.7'}>
                    <IconGuitarPickFilled size={23} />
                  </ActionIcon>
                </Tooltip>
              )}

              {song.youtubeLink && (
                <Tooltip label={'Open Youtube'}>
                  <ActionIcon variant={'transparent'} c={'red.7'}>
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
      />
    </RightSideEntityDrawer>
  )
}

export default SongDrawer
