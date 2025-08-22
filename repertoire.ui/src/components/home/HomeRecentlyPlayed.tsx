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
import SongProperty from '../../types/enums/properties/SongProperty.ts'
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
import ArtistProperty from '../../types/enums/properties/ArtistProperty.ts'
import AlbumProperty from '../../types/enums/properties/AlbumProperty.ts'
import { useGetAlbumsQuery } from '../../state/api/albumsApi.ts'
import HomeRecentlyPlayedAlbumCard from './recently-played/HomeRecentlyPlayedAlbumCard.tsx'
import HomeRecentlyPlayedArtistCard from './recently-played/HomeRecentlyPlayedArtistCard.tsx'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/keys/LocalStorageKeys.ts'
import HomeRecentlyPlayedEntity from '../../types/enums/HomeRecentlyPlayedEntity.ts'

interface TabButtonProps extends ActionIconProps {
  tabValue: HomeRecentlyPlayedEntity
  tab: HomeRecentlyPlayedEntity
  setTab: Dispatch<SetStateAction<HomeRecentlyPlayedEntity>>
  setIndicatorRef: (val: HomeRecentlyPlayedEntity) => (node: HTMLButtonElement) => void
}

const TabButton = ({ tabValue, tab, setTab, setIndicatorRef, ...others }: TabButtonProps) => (
  <ActionIcon
    variant={'grey'}
    size={'sm'}
    aria-selected={tab === tabValue}
    ref={setIndicatorRef(tabValue)}
    onClick={() => setTab(tabValue)}
    sx={{
      // disable active so that the floating indicator can work smoothly
      '&:active': { transform: 'scale(1)' }
    }}
    style={(theme) => ({ ...(tab === tabValue && { color: theme.colors.primary[4] }) })}
    {...others}
  />
)

function HomeRecentlyPlayed() {
  const { ref, hovered } = useHover()

  const [tab, setTab] = useLocalStorage<HomeRecentlyPlayedEntity>({
    key: LocalStorageKeys.HomeRecentlyPlayedEntity,
    defaultValue: HomeRecentlyPlayedEntity.Songs
  })
  const [indicatorRootRef, setIndicatorRootRef] = useState<HTMLDivElement>()
  const [indicatorsRefs, setIndicatorsRefs] = useState<
    Record<HomeRecentlyPlayedEntity, HTMLButtonElement>
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
  >({})
  const setIndicatorRef = (tab: HomeRecentlyPlayedEntity) => (node: HTMLButtonElement) => {
    indicatorsRefs[tab] = node
    setIndicatorsRefs(indicatorsRefs)
  }

  const artistsOrderBy = useOrderBy([
    { property: ArtistProperty.LastPlayed, type: OrderType.Descending },
    { property: ArtistProperty.Progress, type: OrderType.Descending },
    { property: ArtistProperty.Name }
  ])
  const artistsSearchBy = useSearchBy([
    { property: ArtistProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: artists, isLoading: isArtistsLoading } = useGetArtistsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: artistsOrderBy,
      searchBy: artistsSearchBy
    },
    { skip: tab !== HomeRecentlyPlayedEntity.Artists }
  )

  const albumsOrderBy = useOrderBy([
    { property: AlbumProperty.LastPlayed, type: OrderType.Descending },
    { property: AlbumProperty.Progress, type: OrderType.Descending },
    { property: AlbumProperty.Title }
  ])
  const albumsSearchBy = useSearchBy([
    { property: AlbumProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: albumsOrderBy,
      searchBy: albumsSearchBy
    },
    { skip: tab !== HomeRecentlyPlayedEntity.Albums }
  )

  const songsOrderBy = useOrderBy([
    { property: SongProperty.LastPlayed, type: OrderType.Descending },
    { property: SongProperty.Progress, type: OrderType.Descending },
    { property: SongProperty.Title }
  ])
  const songsSearchBy = useSearchBy([
    { property: SongProperty.LastPlayed, operator: FilterOperator.IsNotNull }
  ])
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: songsOrderBy,
      searchBy: songsSearchBy
    },
    { skip: tab !== HomeRecentlyPlayedEntity.Songs }
  )

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
              tabValue={HomeRecentlyPlayedEntity.Artists}
              tab={tab}
              setTab={setTab}
              setIndicatorRef={setIndicatorRef}
            >
              <CustomIconUserAlt size={11} />
            </TabButton>
            <TabButton
              aria-label={'recently-played-albums'}
              tabValue={HomeRecentlyPlayedEntity.Albums}
              tab={tab}
              setTab={setTab}
              setIndicatorRef={setIndicatorRef}
            >
              <CustomIconAlbumVinyl size={11} />
            </TabButton>
            <TabButton
              aria-label={'recently-played-songs'}
              tabValue={HomeRecentlyPlayedEntity.Songs}
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

        {((artists?.models.length !== 0 && tab === HomeRecentlyPlayedEntity.Artists) ||
          (albums?.models.length !== 0 && tab === HomeRecentlyPlayedEntity.Albums) ||
          (songs?.models.length !== 0 && tab === HomeRecentlyPlayedEntity.Songs)) && (
          <Grid columns={12} pl={'lg'} pr={'sm'}>
            <Grid.Col span={{ base: 4.5, md: 10, xxl: 4.5 }}>
              <Text fz={'sm'} fw={500} c={'gray.5'}>
                {tab === HomeRecentlyPlayedEntity.Artists ? 'Name' : 'Title'}
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

        {artists?.models.length === 0 && tab === HomeRecentlyPlayedEntity.Artists && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no artists yet to display
          </Text>
        )}
        {albums?.models.length === 0 && tab === HomeRecentlyPlayedEntity.Albums && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no albums yet to display
          </Text>
        )}
        {songs?.models.length === 0 && tab === HomeRecentlyPlayedEntity.Songs && (
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
            {tab === HomeRecentlyPlayedEntity.Artists &&
              (isArtistsLoading ? (
                <HomeRecentlyPlayedLoader />
              ) : (
                artists.models.map((artist) => (
                  <HomeRecentlyPlayedArtistCard key={artist.id} artist={artist} />
                ))
              ))}
            {tab === HomeRecentlyPlayedEntity.Albums &&
              (isAlbumsLoading ? (
                <HomeRecentlyPlayedLoader />
              ) : (
                albums.models.map((album) => (
                  <HomeRecentlyPlayedAlbumCard key={album.id} album={album} />
                ))
              ))}
            {tab === HomeRecentlyPlayedEntity.Songs &&
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
