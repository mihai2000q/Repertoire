import { forwardRef, useRef, useState } from 'react'
import { ActionIcon, Button, Center, Group, ScrollArea, Stack, Text, Title } from '@mantine/core'
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react'
import { useDidUpdate, useViewportSize } from '@mantine/hooks'
import { useGetAlbumsQuery } from '../../state/api/albumsApi.ts'
import { useGetArtistsQuery } from '../../state/api/artistsApi.ts'
import { useGetSongsQuery } from '../../state/api/songsApi.ts'
import HomeAlbumsLoader from './top/loader/HomeAlbumsLoader.tsx'
import HomeAlbumCard from './top/HomeAlbumCard.tsx'
import HomeSongCard from './top/HomeSongCard.tsx'
import HomeSongsLoader from './top/loader/HomeSongsLoader.tsx'
import HomeArtistsLoader from './top/loader/HomeArtistsLoader.tsx'
import HomeArtistCard from './top/HomeArtistCard.tsx'
import SongProperty from '../../types/enums/SongProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/LocalStorageKeys.ts'

enum TopEntity {
  Songs,
  Albums,
  Artists
}

const TabButton = ({
  children,
  selected,
  onClick
}: {
  children: string
  selected: boolean
  onClick: () => void
}) => (
  <Button
    size={'compact-sm'}
    aria-selected={selected}
    px={'xs'}
    sx={(theme) => ({
      color: theme.colors.gray[6],
      backgroundColor: 'transparent',

      '&:hover': {
        color: theme.colors.gray[9],
        backgroundColor: theme.colors.gray[1],
        transform: 'scale(1)'
      },

      ...(selected && {
        color: theme.colors.primary[4],
        backgroundColor: 'transparent',
        '&:hover': {
          color: theme.colors.primary[4],
          backgroundColor: theme.colors.primary[0],
          transform: 'scale(1)'
        }
      }),

      '&:active': {
        transform: 'unset'
      }
    })}
    style={{ transition: '0.25s' }}
    onClick={onClick}
  >
    {children}
  </Button>
)

const HomeTop = forwardRef<HTMLDivElement>((_, ref) => {
  const songsOrderBy = useOrderBy([
    { property: SongProperty.Progress, type: OrderType.Descending },
    { property: SongProperty.Title }
  ])
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: songsOrderBy
  })

  const albumsOrderBy = useOrderBy([
    { property: AlbumProperty.Progress, type: OrderType.Descending },
    { property: AlbumProperty.Title }
  ])
  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: albumsOrderBy
  })

  const artistsOrderBy = useOrderBy([
    { property: ArtistProperty.Progress, type: OrderType.Descending },
    { property: ArtistProperty.Name }
  ])
  const { data: artists, isLoading: isArtistsLoading } = useGetArtistsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: artistsOrderBy
  })

  const topRef = useRef<HTMLDivElement>(null)

  const { width } = useViewportSize()

  const [topEntity, setTopEntity] = useLocalStorage({
    key: LocalStorageKeys.HomeTopEntity,
    defaultValue: TopEntity.Albums
  })

  const [disableBack, setDisableBack] = useState(false)
  const [disableForward, setDisableForward] = useState(false)
  useDidUpdate(() => {
    setDisableBack(topRef.current?.scrollLeft === 0)
    setDisableForward(topRef.current?.scrollWidth === topRef.current?.clientWidth)
  }, [topRef.current, width, topEntity])

  const handleTopNav = (direction: 'left' | 'right') => {
    if (!topRef.current) return
    topRef.current.scrollBy({ left: direction === 'left' ? -400 : 400, behavior: 'smooth' })
  }

  const handleOnScroll = () => {
    const viewport = topRef.current
    setDisableBack(viewport?.scrollLeft <= 0)
    setDisableForward(viewport?.scrollWidth <= viewport?.clientWidth + viewport?.scrollLeft)
  }

  return (
    <Stack ref={ref} gap={0} aria-label={'top'}>
      <Title px={'xl'} order={2} fw={800} lh={1} mb={'xs'} fz={'max(3vw, 36px)'}>
        Welcome Back
      </Title>

      <Group px={'xl'} gap={0} justify={'space-between'}>
        <Group gap={0}>
          <TabButton
            selected={topEntity === TopEntity.Songs}
            onClick={() => setTopEntity(TopEntity.Songs)}
          >
            Songs
          </TabButton>
          <TabButton
            selected={topEntity === TopEntity.Albums}
            onClick={() => setTopEntity(TopEntity.Albums)}
          >
            Albums
          </TabButton>
          <TabButton
            selected={topEntity === TopEntity.Artists}
            onClick={() => setTopEntity(TopEntity.Artists)}
          >
            Artists
          </TabButton>
        </Group>

        <Group gap={4}>
          <ActionIcon
            aria-label={'back'}
            size={'lg'}
            variant={'grey'}
            radius={'50%'}
            disabled={disableBack}
            onClick={() => handleTopNav('left')}
          >
            <IconChevronLeft size={20} />
          </ActionIcon>

          <ActionIcon
            aria-label={'forward'}
            size={'lg'}
            variant={'grey'}
            radius={'50%'}
            disabled={disableForward}
            onClick={() => handleTopNav('right')}
          >
            <IconChevronRight size={20} />
          </ActionIcon>
        </Group>
      </Group>

      {topEntity === TopEntity.Songs && songs?.models.length === 0 && (
        <Center h={170}>
          <Text c={'gray.6'} fw={500}>
            There are no songs yet to display
          </Text>
        </Center>
      )}
      {topEntity === TopEntity.Albums && albums?.models.length === 0 && (
        <Center h={170}>
          <Text c={'gray.6'} fw={500}>
            There are no albums yet to display
          </Text>
        </Center>
      )}
      {topEntity === TopEntity.Artists && artists?.models.length === 0 && (
        <Center h={158}>
          <Text c={'gray.6'} fw={500} pt={12}>
            There are no artists yet to display
          </Text>
        </Center>
      )}

      <ScrollArea
        viewportRef={topRef}
        viewportProps={{ onScroll: handleOnScroll }}
        scrollbars={'x'}
        offsetScrollbars={'x'}
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
              linear-gradient(to right, transparent 85%, ${theme.white}),
              linear-gradient(to left, transparent 97%, ${theme.white})
            `
          }
        })}
      >
        <Group
          wrap={'nowrap'}
          align={'start'}
          pl={'xl'}
          pr={'5vw'}
          pt={'lg'}
          pb={topEntity === TopEntity.Artists && 'md'}
          gap={topEntity === TopEntity.Artists ? 'sm' : 'lg'}
          style={{ transition: 'padding-bottom 0.3s' }}
        >
          {topEntity === TopEntity.Songs &&
            (isSongsLoading || !songs ? (
              <HomeSongsLoader />
            ) : (
              songs.models.map((song) => <HomeSongCard key={song.id} song={song} />)
            ))}
          {topEntity === TopEntity.Albums &&
            (isAlbumsLoading || !albums ? (
              <HomeAlbumsLoader />
            ) : (
              albums.models.map((album) => <HomeAlbumCard key={album.id} album={album} />)
            ))}
          {topEntity === TopEntity.Artists &&
            (isArtistsLoading || !artists ? (
              <HomeArtistsLoader />
            ) : (
              artists.models.map((artist) => <HomeArtistCard key={artist.id} artist={artist} />)
            ))}
        </Group>
      </ScrollArea>
    </Stack>
  )
})

HomeTop.displayName = 'HomeTop'

export default HomeTop
