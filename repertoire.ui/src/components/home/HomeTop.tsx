import { useRef, useState } from 'react'
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

function HomeTop() {
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: ['progress desc', 'title asc']
  })

  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: ['updated_at desc', 'title asc']
  })

  const { data: artists, isLoading: isArtistsLoading } = useGetArtistsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: ['updated_at desc', 'name asc']
  })

  const topRef = useRef<HTMLDivElement>(null)

  const { width } = useViewportSize()

  const [topEntity, setTopEntity] = useState(TopEntity.Albums)

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
    <Stack gap={0} aria-label={'top'}>
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
          <Text c={'gray.6'} fw={500}>
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
      >
        <Group
          wrap={'nowrap'}
          align={'start'}
          px={'xl'}
          pt={'lg'}
          pb={topEntity === TopEntity.Artists && 'sm'}
          gap={'lg'}
        >
          {topEntity === TopEntity.Songs &&
            (isSongsLoading ? (
              <HomeSongsLoader />
            ) : (
              songs.models.map((song) => <HomeSongCard key={song.id} song={song} />)
            ))}
          {topEntity === TopEntity.Albums &&
            (isAlbumsLoading ? (
              <HomeAlbumsLoader />
            ) : (
              albums.models.map((album) => <HomeAlbumCard key={album.id} album={album} />)
            ))}
          {topEntity === TopEntity.Artists &&
            (isArtistsLoading ? (
              <HomeArtistsLoader />
            ) : (
              artists.models.map((artist) => <HomeArtistCard key={artist.id} artist={artist} />)
            ))}
        </Group>
      </ScrollArea>
    </Stack>
  )
}

export default HomeTop
