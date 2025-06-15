import {
  ActionIcon,
  Avatar,
  Card,
  CardProps,
  Center,
  Group,
  Menu,
  ScrollArea,
  Skeleton,
  Stack,
  Text
} from '@mantine/core'
import Artist from '../../types/models/Artist.ts'
import { useGetArtistsQuery } from '../../state/api/artistsApi.ts'
import { useRef, useState } from 'react'
import { useDidUpdate, useHover, useViewportSize } from '@mantine/hooks'
import { IconChevronLeft, IconChevronRight, IconEye } from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { openArtistDrawer } from '../../state/slice/globalSlice.ts'
import CustomIconUserAlt from '../@ui/icons/CustomIconUserAlt.tsx'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import { useNavigate } from 'react-router-dom'
import useContextMenu from '../../hooks/useContextMenu.ts'

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
  const navigate = useNavigate()

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const isSelected = hovered || openedMenu

  function handleClick() {
    dispatch(openArtistDrawer(artist.id))
  }

  function handleViewDetails() {
    navigate(`/artist/${artist.id}`)
  }

  return (
    <Stack
      align={'center'}
      gap={'xxs'}
      w={60}
      sx={{ transition: '0.2s', ...(isSelected && { transform: 'scale(1.1)' }) }}
    >
      <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
        <Menu.Target>
          <Avatar
            ref={ref}
            size={'lg'}
            src={artist.imageUrl}
            alt={artist.imageUrl && artist.name}
            bg={'gray.0'}
            sx={(theme) => ({
              cursor: 'pointer',
              transition: '0.2s',
              boxShadow: theme.shadows.sm,
              ...(isSelected && { boxShadow: theme.shadows.xl })
            })}
            onClick={handleClick}
            onContextMenu={openMenu}
          >
            <Center c={'gray.7'}>
              <CustomIconUserAlt aria-label={`default-icon-${artist.name}`} size={25} />
            </Center>
          </Avatar>
        </Menu.Target>

        <Menu.Dropdown {...menuDropdownProps}>
          <Menu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
            View Details
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>

      <Text ta={'center'} fw={500} lineClamp={2}>
        {artist.name}
      </Text>
    </Stack>
  )
}

function HomeRecentArtists({ ...others }: CardProps) {
  const orderBy = useOrderBy([
    { property: ArtistProperty.LastModified, type: OrderType.Descending }
  ])

  const { data: artists, isLoading } = useGetArtistsQuery({
    pageSize: 15,
    currentPage: 1,
    orderBy: orderBy
  })

  const viewportRef = useRef<HTMLDivElement>(null)

  const { width } = useViewportSize()

  const [disableBack, setDisableBack] = useState(false)
  const [disableForward, setDisableForward] = useState(false)
  useDidUpdate(() => {
    setDisableBack(viewportRef.current?.scrollLeft === 0)
    setDisableForward(viewportRef.current?.scrollWidth === viewportRef.current?.clientWidth)
  }, [viewportRef.current, width])
  useDidUpdate(() => {
    const frame = requestAnimationFrame(() => {
      setDisableBack(viewportRef.current?.scrollLeft === 0)
      setDisableForward(viewportRef.current?.scrollWidth === viewportRef.current?.clientWidth)
    })
    return () => cancelAnimationFrame(frame)
  }, [artists])

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
            Recent Artists
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

        {artists?.models.length === 0 && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'xl'}>
            There are no artists yet to display
          </Text>
        )}

        <ScrollArea
          h={'100%'}
          scrollbars={'x'}
          scrollbarSize={7}
          viewportRef={viewportRef}
          viewportProps={{ onScroll: handleOnScroll }}
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
                linear-gradient(to right, transparent 97%, ${theme.white}),
                linear-gradient(to left, transparent 97%, ${theme.white})
              `
            }
          })}
          styles={{ viewport: { '> div': { display: 'flex' } } }}
        >
          <Group wrap={'nowrap'} align={'start'} px={'md'} pt={'xs'} pb={'md'}>
            {isLoading || !artists ? (
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

export default HomeRecentArtists
