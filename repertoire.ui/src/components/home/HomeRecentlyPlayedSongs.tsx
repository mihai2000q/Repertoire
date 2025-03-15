import { useGetSongsQuery } from '../../state/api/songsApi.ts'
import {
  alpha,
  Avatar,
  Card,
  Center,
  Grid,
  Group,
  ScrollArea,
  Skeleton,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { IconClock } from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { useHover } from '@mantine/hooks'
import { openArtistDrawer, openSongDrawer } from '../../state/slice/globalSlice.ts'
import { MouseEvent } from 'react'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import SongProgressBar from '../@ui/misc/SongProgressBar.tsx'
import dayjs from 'dayjs'
import Song from '../../types/models/Song.ts'

function Loader() {
  return (
    <>
      {Array.from(Array(20)).map((_, i) => (
        <Group key={i} pl={'lg'} pr={'xxs'} py={'xs'}>
          <Skeleton
            radius={'md'}
            h={38}
            w={38}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Grid flex={1} columns={12} align={'center'}>
            <Grid.Col span={{ base: 5, md: 8, xxl: 5 }}>
              <Stack gap={0}>
                <Skeleton w={125} h={15} mb={4} />
                <Skeleton w={70} h={10} />
              </Stack>
            </Grid.Col>
            <Grid.Col span={4} display={{ base: 'block', md: 'none', xxl: 'block' }}>
              <Skeleton flex={1} h={12} px={'xs'} />
            </Grid.Col>
            <Grid.Col span={{ base: 3, md: 4, xxl: 3 }}>
              <Center>
                <Skeleton w={50} h={15} px={'md'} />
              </Center>
            </Grid.Col>
          </Grid>
        </Group>
      ))}
    </>
  )
}

function LocalSongCard({ song }: { song: Song }) {
  const dispatch = useAppDispatch()

  const { ref, hovered } = useHover()

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleArtistClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openArtistDrawer(song.artist.id))
  }

  return (
    <Group
      ref={ref}
      wrap={'nowrap'}
      sx={(theme) => ({
        transition: '0.3s',
        border: '1px solid transparent',
        ...(hovered && {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.primary[0], 0.15)
        })
      })}
      pl={'lg'}
      pr={'xxs'}
      py={'xs'}
      onClick={handleClick}
    >
      <Avatar
        radius={'md'}
        src={song.imageUrl ?? song.album?.imageUrl ?? songPlaceholder}
        alt={song.title}
        style={(theme) => ({ boxShadow: theme.shadows.sm })}
      />

      <Grid flex={1} columns={12} align={'center'}>
        <Grid.Col span={{ base: 5, md: 8, xxl: 5 }}>
          <Stack gap={0} style={{ overflow: 'hidden' }}>
            <Text fw={600} lineClamp={1}>
              {song.title}
            </Text>
            {song.artist && (
              <Group>
                <Text
                  fz={'sm'}
                  c={'dimmed'}
                  lineClamp={1}
                  sx={{ '&:hover': { textDecoration: 'underline' } }}
                  style={{ cursor: 'pointer' }}
                  onClick={handleArtistClick}
                >
                  {song.artist.name}
                </Text>
              </Group>
            )}
          </Stack>
        </Grid.Col>
        <Grid.Col span={4} display={{ base: 'block', md: 'none', xxl: 'block' }}>
          <SongProgressBar progress={song.progress} mx={'xs'} />
        </Grid.Col>
        <Grid.Col span={{ base: 3, md: 4, xxl: 3 }} px={'md'}>
          <Tooltip
            label={`Song was played last time on ${dayjs(song.lastTimePlayed).format('D MMMM YYYY [at] hh:mm A')}`}
            openDelay={400}
            disabled={!song.lastTimePlayed}
          >
            <Text ta={'center'} fz={'sm'} fw={500} c={'dimmed'} truncate={'end'}>
              {song.lastTimePlayed ? dayjs(song.lastTimePlayed).format('DD MMM') : 'never'}
            </Text>
          </Tooltip>
        </Grid.Col>
      </Grid>
    </Group>
  )
}

function HomeRecentlyPlayedSongs() {
  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: ['last_time_played desc', 'progress desc', 'title desc'],
    searchBy: ['last_time_played IS NOT NULL']
  })

  return (
    <Card aria-label={'recently-played-songs'} radius={'lg'} variant={'panel'} p={0}>
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Text c={'gray.7'} fz={'lg'} fw={800} mx={'lg'} my={'md'}>
          Recently Played
        </Text>

        {songs?.models.length !== 0 && (
          <Grid columns={12} pl={'lg'} pr={'sm'}>
            <Grid.Col span={{ base: 4.5, md: 10, xxl: 4.5 }}>
              <Text fz={'sm'} fw={500} c={'gray.5'}>
                Title
              </Text>
            </Grid.Col>
            <Grid.Col span={6} display={{ base: 'block', md: 'none', xxl: 'block' }}>
              <Text ta={'center'} fz={'sm'} fw={500} c={'gray.5'}>
                Progress
              </Text>
            </Grid.Col>
            <Grid.Col span={{ base: 1.5, md: 2, xxl: 1.5 }} c={'gray.5'}>
              <IconClock size={15} aria-label={'last-time-played-icon'} />
            </Grid.Col>
          </Grid>
        )}

        {songs?.models.length === 0 && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no songs yet to display
          </Text>
        )}

        <ScrollArea scrollbars={'y'} scrollbarSize={7}>
          <Stack gap={'xxs'} h={'100%'}>
            {isLoading ? (
              <Loader />
            ) : (
              songs.models.map((song) => <LocalSongCard key={song.id} song={song} />)
            )}
          </Stack>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeRecentlyPlayedSongs
