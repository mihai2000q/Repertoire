import {
  ActionIcon,
  Avatar,
  Card,
  CardProps,
  Group,
  ScrollArea,
  Skeleton,
  Stack,
  Text
} from '@mantine/core'
import Artist from '../../types/models/Artist.ts'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import { useGetArtistsQuery } from '../../state/api/artistsApi.ts'
import { useRef, useState } from 'react'
import { useDidUpdate, useHover, useViewportSize } from '@mantine/hooks'
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react'
import {useAppDispatch} from "../../state/store.ts";
import {openArtistDrawer} from "../../state/slice/globalSlice.ts";

function Loader() {
  return (
    <Group wrap={'nowrap'} gap={'lg'} data-testid={'artists-loader'}>
      {Array.from(Array(15)).map((_, i) => (
        <Stack key={i} gap={6} align={'center'} w={60}>
          <Skeleton
            radius={'50%'}
            h={56}
            w={56}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Stack gap={'xxs'} align={'center'}>
            <Skeleton w={60} h={12} />
            <Skeleton w={40} h={12} />
          </Stack>
        </Stack>
      ))}
    </Group>
  )
}

function LocalArtistCard({ artist }: { artist: Artist }) {
  const dispatch = useAppDispatch()
  const { ref, hovered } = useHover()

  function handleClick() {
    dispatch(openArtistDrawer(artist.id))
  }

  return (
    <Stack
      align={'center'}
      gap={'xxs'}
      w={60}
      sx={{ transition: '0.2s', ...(hovered && { transform: 'scale(1.1)' }) }}
    >
      <Avatar
        ref={ref}
        size={'lg'}
        src={artist.imageUrl ?? artistPlaceholder}
        alt={artist.name}
        sx={(theme) => ({
          cursor: 'pointer',
          transition: '0.2s',
          boxShadow: theme.shadows.sm,
          '&:hover': {
            boxShadow: theme.shadows.xl
          }
        })}
        onClick={handleClick}
      />

      <Text ta={'center'} fw={500} lineClamp={2}>
        {artist.name}
      </Text>
    </Stack>
  )
}

function HomeTopArtists({ ...others }: CardProps) {
  const { data: artists, isLoading } = useGetArtistsQuery({
    pageSize: 15,
    currentPage: 1,
    orderBy: ['updated_at desc', 'name asc']
  })

  const viewportRef = useRef<HTMLDivElement>(null)

  const { width } = useViewportSize()

  const [disableBack, setDisableBack] = useState(false)
  const [disableForward, setDisableForward] = useState(false)
  useDidUpdate(() => {
    setDisableBack(viewportRef.current?.scrollLeft === 0)
    setDisableForward(viewportRef.current?.scrollWidth === viewportRef.current?.clientWidth)
  }, [viewportRef.current, width])

  const handleTopNav = (direction: 'left' | 'right') => {
    if (!viewportRef.current) return
    viewportRef.current.scrollBy({ left: direction === 'left' ? -250 : 250, behavior: 'smooth' })
  }

  const handleOnScroll = () => {
    const viewport = viewportRef.current
    setDisableBack(viewport?.scrollLeft <= 0)
    setDisableForward(viewport?.scrollWidth <= viewport?.clientWidth + viewport?.scrollLeft)
  }

  return (
    <Card aria-label={'top-artists'} variant={'panel'} {...others} p={0} pb={0}>
      <Stack h={'100%'} gap={'xs'}>
        <Group px={'md'} pt={'sm'} justify={'space-between'}>
          <Text c={'gray.7'} fz={'lg'} fw={800}>
            Top Artists
          </Text>

          <Group gap={4}>
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

        <ScrollArea
          h={'100%'}
          scrollbars={'x'}
          scrollbarSize={7}
          viewportRef={viewportRef}
          viewportProps={{ onScroll: handleOnScroll }}
        >
          <Group wrap={'nowrap'} h={'100%'} align={'start'} px={'md'} pt={'xs'} pb={'md'}>
            {isLoading ? (
              <Loader />
            ) : (
              artists.models.map((artist) => <LocalArtistCard key={artist.id} artist={artist} />)
            )}
          </Group>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeTopArtists
