import { useGetSongsQuery } from '../../state/api/songsApi.ts'
import { Card, Grid, ScrollArea, Stack, Text } from '@mantine/core'
import { IconClock } from '@tabler/icons-react'
import SongProperty from '../../types/enums/SongProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import useSearchBy from '../../hooks/api/useSearchBy.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import HomeRecentlyPlayedSongCard from './recently-played/HomeRecentlyPlayedSongCard.tsx'
import HomeRecentlyPlayedLoader from './recently-played/HomeRecentlyPlayedLoader.tsx'

function HomeRecentlyPlayed() {
  const orderBy = useOrderBy([
    { property: SongProperty.LastPlayed, type: OrderType.Descending },
    { property: SongProperty.Progress, type: OrderType.Descending },
    { property: SongProperty.Title }
  ])
  const searchBy = useSearchBy([
    { property: SongProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])

  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: orderBy,
    searchBy: searchBy
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

        <ScrollArea
          scrollbars={'y'}
          scrollbarSize={7}
          sx={(theme) => ({
            '&::after': {
              content: '""',
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: '100%',
              pointerEvents: 'none',
              background: `
                linear-gradient(to top, transparent 98%, ${theme.white}),
                linear-gradient(to bottom, transparent 96%, ${theme.white})
              `
            }
          })}
        >
          <Stack gap={'xxs'}>
            {isLoading ? (
              <HomeRecentlyPlayedLoader />
            ) : (
              songs.models.map((song) => <HomeRecentlyPlayedSongCard key={song.id} song={song} />)
            )}
          </Stack>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeRecentlyPlayed
