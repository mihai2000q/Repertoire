import {
  Anchor,
  AspectRatio,
  Avatar,
  Button,
  Card,
  Checkbox,
  Divider,
  Grid,
  Group,
  Image,
  Progress,
  Stack,
  Text,
  Title,
  Tooltip
} from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { useParams } from 'react-router-dom'
import { openAlbumDrawer, openArtistDrawer } from '../../state/globalSlice.ts'
import SongLoader from '../../components/songs/loader/SongLoader.tsx'
import { useGetSongQuery } from '../../state/songsApi.ts'
import useDifficultyInfo from '../../hooks/songs/useDifficultyInfo.ts'
import Difficulty from '../../utils/enums/Difficulty.ts'
import { IconBrandYoutube, IconGuitarPick } from '@tabler/icons-react'
import SongSections from '../../components/songs/SongSections.tsx'

const NotSet = () => (
  <Text fz={'sm'} c={'dimmed'} fs={'oblique'} inline>
    not set
  </Text>
)

function Song() {
  const dispatch = useAppDispatch()

  const params = useParams()
  const songId = params['id'] ?? ''

  const { data: song, isLoading } = useGetSongQuery(songId)

  const { number: difficultyNumber, color: difficultyColor } = useDifficultyInfo(song?.difficulty)

  function handleAlbumClick() {
    dispatch(openAlbumDrawer(song.album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(song.artist.id))
  }

  if (isLoading) return <SongLoader />

  return (
    <Stack>
      <Group>
        <AspectRatio>
          <Image
            h={150}
            src={song.imageUrl}
            fallbackSrc={albumPlaceholder}
            radius={'lg'}
            sx={(theme) => ({
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
                {/*<Avatar radius={'md'} size={40} src={song.album.imageUrl ?? userPlaceholder} />*/}
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

      <Divider />

      <Group align="start">
        <Stack flex={1}>
          <Card variant={'panel'} p={'md'}>
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
                    <Checkbox disabled checked={song.isRecorded} style={{ cursor: 'default' }} />
                  ) : (
                    <Text fw={600}>No</Text>
                  )}
                </Grid.Col>
              </Grid>
            </Stack>
          </Card>

          <Card variant={'panel'} p={'md'}>
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
          </Card>
        </Stack>

        <Stack flex={1.75}>
          <Card variant={'panel'} p={'md'}>
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
          </Card>

          <SongSections songId={songId} sections={song.sections} />
        </Stack>
      </Group>
    </Stack>
  )
}

export default Song
