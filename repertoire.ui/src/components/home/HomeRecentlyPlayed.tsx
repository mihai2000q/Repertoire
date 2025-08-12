import { useGetSongsQuery } from '../../state/api/songsApi.ts'
import {
  ActionIcon,
  ActionIconProps,
  alpha,
  Card,
  FloatingIndicator,
  Grid,
  Group,
  ScrollArea,
  Stack,
  Text
} from '@mantine/core'
import { IconClock } from '@tabler/icons-react'
import SongProperty from '../../types/enums/SongProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import useSearchBy from '../../hooks/api/useSearchBy.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import HomeRecentlyPlayedSongCard from './recently-played/HomeRecentlyPlayedSongCard.tsx'
import HomeRecentlyPlayedLoader from './recently-played/HomeRecentlyPlayedLoader.tsx'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import CustomIconAlbumVinyl from '../@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconMusicNoteEighth from '../@ui/icons/CustomIconMusicNoteEighth.tsx'
import { Dispatch, SetStateAction, useState } from 'react'
import { useHover } from '@mantine/hooks'
import { useGetArtistsQuery } from '../../state/api/artistsApi.ts'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import { useGetAlbumsQuery } from '../../state/api/albumsApi.ts'
import HomeRecentlyPlayedAlbumCard from './recently-played/HomeRecentlyPlayedAlbumCard.tsx'
import HomeRecentlyPlayedArtistCard from './recently-played/HomeRecentlyPlayedArtistCard.tsx'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/LocalStorageKeys.ts'

enum RecentlyPlayedTab {
  Artists,
  Albums,
  Songs
}

interface TabButtonProps extends ActionIconProps {
  tabValue: RecentlyPlayedTab
  tab: RecentlyPlayedTab
  setTab: Dispatch<SetStateAction<RecentlyPlayedTab>>
  setIndicatorRef: (val: RecentlyPlayedTab) => (node: HTMLButtonElement) => void
}

const TabButton = ({ tabValue, tab, setTab, setIndicatorRef, ...others }: TabButtonProps) => (
  <ActionIcon
    variant={'subtle'}
    size={'sm'}
    aria-selected={tab === tabValue}
    ref={setIndicatorRef(tabValue)}
    onClick={() => setTab(tabValue)}
    sx={(theme) => ({
      color: theme.colors.gray[5],
      // disable active so that the floating indicator can work smoothly
      '&:active': { transform: 'scale(1)' },
      '&:hover': {
        color: theme.colors.gray[6],
        backgroundColor: theme.colors.gray[2],
        shadows: theme.shadows.lg
      },
      '&[aria-selected="true"]': { color: theme.colors.primary[4] }
    })}
    {...others}
  />
)

function HomeRecentlyPlayed() {
  const { ref, hovered } = useHover()

  const artistsOrderBy = useOrderBy([
    { property: ArtistProperty.LastPlayed, type: OrderType.Descending },
    { property: ArtistProperty.Progress, type: OrderType.Descending },
    { property: ArtistProperty.Name }
  ])
  const artistsSearchBy = useSearchBy([
    { property: ArtistProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: artists, isLoading: isArtistsLoading } = useGetArtistsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: artistsOrderBy,
    searchBy: artistsSearchBy
  })

  const albumsOrderBy = useOrderBy([
    { property: AlbumProperty.LastPlayed, type: OrderType.Descending },
    { property: AlbumProperty.Progress, type: OrderType.Descending },
    { property: AlbumProperty.Title }
  ])
  const albumsSearchBy = useSearchBy([
    { property: AlbumProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: albumsOrderBy,
    searchBy: albumsSearchBy
  })

  const songsOrderBy = useOrderBy([
    { property: SongProperty.LastPlayed, type: OrderType.Descending },
    { property: SongProperty.Progress, type: OrderType.Descending },
    { property: SongProperty.Title }
  ])
  const songsSearchBy = useSearchBy([
    { property: SongProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: songsOrderBy,
    searchBy: songsSearchBy
  })

  const [tab, setTab] = useLocalStorage<RecentlyPlayedTab>({
    key: LocalStorageKeys.HomeRecentlyPlayedEntity,
    defaultValue: RecentlyPlayedTab.Songs
  })
  const [indicatorRootRef, setIndicatorRootRef] = useState<HTMLDivElement>()
  const [indicatorsRefs, setIndicatorsRefs] = useState<
    Record<RecentlyPlayedTab, HTMLButtonElement>
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
  >({})
  const setIndicatorRef = (tab: RecentlyPlayedTab) => (node: HTMLButtonElement) => {
    indicatorsRefs[tab] = node
    setIndicatorsRefs(indicatorsRefs)
  }

  return (
    <Card ref={ref} aria-label={'recently-played'} radius={'lg'} variant={'panel'} p={0}>
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Group justify={'space-between'} mx={'lg'} my={'md'}>
          <Text c={'gray.7'} fz={'lg'} fw={800}>
            Recently Played
          </Text>

          <Group
            ref={setIndicatorRootRef}
            gap={2}
            pos={'relative'}
            style={{ transition: '0.25s', opacity: hovered ? 1 : 0 }}
          >
            <TabButton
              aria-label={'recently-played-artists'}
              tabValue={RecentlyPlayedTab.Artists}
              tab={tab}
              setTab={setTab}
              setIndicatorRef={setIndicatorRef}
            >
              <CustomIconUserAlt size={11} />
            </TabButton>
            <TabButton
              aria-label={'recently-played-albums'}
              tabValue={RecentlyPlayedTab.Albums}
              tab={tab}
              setTab={setTab}
              setIndicatorRef={setIndicatorRef}
            >
              <CustomIconAlbumVinyl size={11} />
            </TabButton>
            <TabButton
              aria-label={'recently-played-songs'}
              tabValue={RecentlyPlayedTab.Songs}
              tab={tab}
              setTab={setTab}
              setIndicatorRef={setIndicatorRef}
            >
              <CustomIconMusicNoteEighth size={13} />
            </TabButton>

            <FloatingIndicator
              parent={indicatorRootRef}
              target={indicatorsRefs[tab]}
              sx={(theme) => ({
                transition: '0.2s ease',
                borderRadius: theme.radius.md,
                backgroundColor: alpha(theme.colors.primary[4], 0.13),
                boxShadow: theme.shadows.md
              })}
            />
          </Group>
        </Group>

        {((artists?.models.length !== 0 && tab === RecentlyPlayedTab.Artists) ||
          (albums?.models.length !== 0 && tab === RecentlyPlayedTab.Albums) ||
          (songs?.models.length !== 0 && tab === RecentlyPlayedTab.Songs)) && (
          <Grid columns={12} pl={'lg'} pr={'sm'}>
            <Grid.Col span={{ base: 4.5, md: 10, xxl: 4.5 }}>
              <Text fz={'sm'} fw={500} c={'gray.5'}>
                {tab === RecentlyPlayedTab.Artists ? 'Name' : 'Title'}
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

        {artists?.models.length === 0 && tab === RecentlyPlayedTab.Artists && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no artists yet to display
          </Text>
        )}
        {albums?.models.length === 0 && tab === RecentlyPlayedTab.Albums && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no albums yet to display
          </Text>
        )}
        {songs?.models.length === 0 && tab === RecentlyPlayedTab.Songs && (
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
            {tab === RecentlyPlayedTab.Artists &&
              (isArtistsLoading ? (
                <HomeRecentlyPlayedLoader />
              ) : (
                artists.models.map((artist) => (
                  <HomeRecentlyPlayedArtistCard key={artist.id} artist={artist} />
                ))
              ))}
            {tab === RecentlyPlayedTab.Albums &&
              (isAlbumsLoading ? (
                <HomeRecentlyPlayedLoader />
              ) : (
                albums.models.map((album) => (
                  <HomeRecentlyPlayedAlbumCard key={album.id} album={album} />
                ))
              ))}
            {tab === RecentlyPlayedTab.Songs &&
              (isSongsLoading ? (
                <HomeRecentlyPlayedLoader />
              ) : (
                songs.models.map((song) => <HomeRecentlyPlayedSongCard key={song.id} song={song} />)
              ))}
          </Stack>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeRecentlyPlayed
