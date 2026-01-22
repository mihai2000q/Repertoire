import { RefObject, useRef, useState } from 'react'
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
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import useLocalStorage from '../../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../../types/enums/keys/LocalStorageKeys.ts'
import HomeTopEntity from '../../types/enums/HomeTopEntity.ts'
import { defaultHomeTopOrderEntities } from '../../data/home/homeTopOrderEntities.ts'
import HomeTopEntityOrderButton from './top/HomeTopEntityOrderButton.tsx'

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

function HomeTop({ ref }: { ref: RefObject<HTMLDivElement> }) {
  const [topEntity, setHomeTopEntity] = useLocalStorage({
    key: LocalStorageKeys.HomeTopEntity,
    defaultValue: HomeTopEntity.Albums
  })
  const [orderEntities, setOrderEntities] = useLocalStorage({
    key: LocalStorageKeys.HomeOrderTopEntities,
    defaultValue: defaultHomeTopOrderEntities,
    serialize: (val) => JSON.stringify([...val]),
    deserialize: (val) => new Map(JSON.parse(val))
  })

  const [artistsOrder, setArtistsOrder] = useState(orderEntities.get(HomeTopEntity.Artists))
  const artistsOrderBy = useOrderBy([artistsOrder])
  const { data: artists, isLoading: isArtistsLoading } = useGetArtistsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: artistsOrderBy
    },
    { skip: topEntity !== HomeTopEntity.Artists }
  )

  const [albumsOrder, setAlbumsOrder] = useState(orderEntities.get(HomeTopEntity.Albums))
  const albumsOrderBy = useOrderBy([albumsOrder])
  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: albumsOrderBy
    },
    { skip: topEntity !== HomeTopEntity.Albums }
  )

  const [songsOrder, setSongsOrder] = useState(orderEntities.get(HomeTopEntity.Songs))
  const songsOrderBy = useOrderBy([songsOrder])
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery(
    {
      pageSize: 20,
      currentPage: 1,
      orderBy: songsOrderBy
    },
    { skip: topEntity !== HomeTopEntity.Songs }
  )

  useDidUpdate(() => {
    setArtistsOrder(orderEntities.get(HomeTopEntity.Artists))
    setAlbumsOrder(orderEntities.get(HomeTopEntity.Albums))
    setSongsOrder(orderEntities.get(HomeTopEntity.Songs))
  }, [JSON.stringify([...orderEntities])])

  const totalCount =
    (topEntity === HomeTopEntity.Artists
      ? artists?.totalCount
      : topEntity === HomeTopEntity.Albums
        ? albums?.totalCount
        : songs?.totalCount) ?? 0

  // Navigation Buttons
  const topRef = useRef<HTMLDivElement>(null)

  const { width } = useViewportSize()

  const [disableBack, setDisableBack] = useState(false)
  const [disableForward, setDisableForward] = useState(false)
  useDidUpdate(() => {
    setDisableBack(topRef.current?.scrollLeft === 0)
    setDisableForward(topRef.current?.scrollWidth === topRef.current?.clientWidth)
  }, [topRef.current, width, topEntity])
  useDidUpdate(() => {
    const frame = requestAnimationFrame(() => {
      setDisableBack(topRef.current?.scrollLeft === 0)
      setDisableForward(topRef.current?.scrollWidth === topRef.current?.clientWidth)
    })
    return () => cancelAnimationFrame(frame)
  }, [albums, songs, artists])

  const handleTopNav = (direction: 'left' | 'right') => {
    topRef.current?.scrollBy({ left: direction === 'left' ? -400 : 400, behavior: 'smooth' })
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
            selected={topEntity === HomeTopEntity.Artists}
            onClick={() => setHomeTopEntity(HomeTopEntity.Artists)}
          >
            Artists
          </TabButton>
          <TabButton
            selected={topEntity === HomeTopEntity.Albums}
            onClick={() => setHomeTopEntity(HomeTopEntity.Albums)}
          >
            Albums
          </TabButton>
          <TabButton
            selected={topEntity === HomeTopEntity.Songs}
            onClick={() => setHomeTopEntity(HomeTopEntity.Songs)}
          >
            Songs
          </TabButton>
        </Group>

        <Group gap={'xxs'}>
          <HomeTopEntityOrderButton
            topEntity={topEntity}
            orderEntities={orderEntities}
            setOrderEntities={setOrderEntities}
            disabled={totalCount === 0}
          />

          <ActionIcon
            aria-label={'back'}
            variant={'grey'}
            radius={'50%'}
            disabled={disableBack}
            onClick={() => handleTopNav('left')}
          >
            <IconChevronLeft size={20} />
          </ActionIcon>

          <ActionIcon
            aria-label={'forward'}
            variant={'grey'}
            radius={'50%'}
            disabled={disableForward}
            onClick={() => handleTopNav('right')}
          >
            <IconChevronRight size={20} />
          </ActionIcon>
        </Group>
      </Group>

      {topEntity === HomeTopEntity.Artists && artists?.models.length === 0 && (
        <Center h={158}>
          <Text c={'gray.6'} fw={500} pt={12}>
            There are no artists yet to display
          </Text>
        </Center>
      )}
      {topEntity === HomeTopEntity.Albums && albums?.models.length === 0 && (
        <Center h={170}>
          <Text c={'gray.6'} fw={500}>
            There are no albums yet to display
          </Text>
        </Center>
      )}
      {topEntity === HomeTopEntity.Songs && songs?.models.length === 0 && (
        <Center h={170}>
          <Text c={'gray.6'} fw={500}>
            There are no songs yet to display
          </Text>
        </Center>
      )}

      <ScrollArea.Autosize
        viewportRef={topRef}
        viewportProps={{ onScroll: handleOnScroll }}
        scrollbars={'x'}
        offsetScrollbars={'x'}
        scrollbarSize={7}
        styles={{
          viewport: {
            '> div': {
              display: 'flex !important',
              minWidth: '100%',
              width: 0
            }
          }
        }}
      >
        <Group
          wrap={'nowrap'}
          align={'start'}
          pl={'xl'}
          pr={'5vw'}
          pt={'lg'}
          pb={topEntity === HomeTopEntity.Artists ? 'md' : 'xxs'}
          gap={topEntity === HomeTopEntity.Artists ? 'sm' : 'lg'}
          style={{ transition: 'padding-bottom 0.3s' }}
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
          {topEntity === HomeTopEntity.Artists &&
            (isArtistsLoading || !artists ? (
              <HomeArtistsLoader />
            ) : (
              artists.models.map((artist) => <HomeArtistCard key={artist.id} artist={artist} />)
            ))}
          {topEntity === HomeTopEntity.Albums &&
            (isAlbumsLoading || !albums ? (
              <HomeAlbumsLoader />
            ) : (
              albums.models.map((album) => <HomeAlbumCard key={album.id} album={album} />)
            ))}
          {topEntity === HomeTopEntity.Songs &&
            (isSongsLoading || !songs ? (
              <HomeSongsLoader />
            ) : (
              songs.models.map((song) => <HomeSongCard key={song.id} song={song} />)
            ))}
        </Group>
      </ScrollArea.Autosize>
    </Stack>
  )
}

export default HomeTop
