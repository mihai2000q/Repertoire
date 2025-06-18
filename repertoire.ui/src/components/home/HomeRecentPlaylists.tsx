import {
  Avatar,
  Card,
  CardProps,
  Center,
  Group,
  ScrollArea,
  SimpleGrid,
  Skeleton,
  Space,
  Stack,
  Text
} from '@mantine/core'
import { useGetPlaylistsQuery } from '../../state/api/playlistsApi.ts'
import Playlist from '../../types/models/Playlist.ts'
import { IconEye, IconPlaylist } from '@tabler/icons-react'
import OrderType from '../../types/enums/OrderType.ts'
import PlaylistProperty from '../../types/enums/PlaylistProperty.ts'
import useOrderBy from '../../hooks/api/useOrderBy.ts'
import { openPlaylistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'
import { useDisclosure, useHover } from '@mantine/hooks'
import { useNavigate } from 'react-router-dom'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'

function Loader() {
  return (
    <>
      {Array.from(Array(20)).map((_, i) => (
        <Group key={i} wrap={'nowrap'}>
          <Skeleton
            radius={'lg'}
            h={60}
            w={60}
            style={(theme) => ({ boxShadow: theme.shadows.md })}
          />
          <Stack gap={'xxs'}>
            <Skeleton w={100} h={13} />
            <Skeleton w={75} h={13} />
          </Stack>
        </Group>
      ))}
    </>
  )
}

function LocalPlaylistCard({ playlist }: { playlist: Playlist }) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  const isSelected = hovered || openedMenu

  function handleClick() {
    dispatch(openPlaylistDrawer(playlist.id))
  }

  function handleViewDetails() {
    navigate(`/playlist/${playlist.id}`)
  }

  return (
    <Group wrap={'nowrap'} gap={0}>
      <ContextMenu shadow={'lg'} opened={openedMenu} onChange={toggleMenu}>
        <ContextMenu.Target>
          <Avatar
            ref={ref}
            radius={'28%'}
            size={60}
            src={playlist.imageUrl}
            alt={playlist.imageUrl && playlist.title}
            bg={'gray.5'}
            sx={(theme) => ({
              aspectRatio: 1,
              cursor: 'pointer',
              transition: '0.2s',
              boxShadow: theme.shadows.sm,
              ...(isSelected && {
                boxShadow: theme.shadows.xl,
                transform: 'scale(1.1)'
              })
            })}
            onClick={handleClick}
          >
            <Center c={'white'}>
              <IconPlaylist
                aria-label={`default-icon-${playlist.title}`}
                size={'100%'}
                style={{ padding: '27%' }}
              />
            </Center>
          </Avatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item leftSection={<IconEye size={14} />} onClick={handleViewDetails}>
            View Details
          </ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <Space ml={{ base: 'xs', xl: 'sm', xxl: 'md' }} style={{ transition: '0.16s' }} />

      <Text fw={500} lineClamp={2}>
        {playlist.title}
      </Text>
    </Group>
  )
}

function HomeRecentPlaylists({ ...others }: CardProps) {
  const orderBy = useOrderBy([
    { property: PlaylistProperty.LastModified, type: OrderType.Descending }
  ])

  const { data: playlists, isLoading } = useGetPlaylistsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: orderBy
  })

  return (
    <Card aria-label={'playlists'} variant={'panel'} {...others} p={0}>
      <Stack h={'100%'} gap={'xs'}>
        <Text c={'gray.7'} fz={'lg'} fw={800} px={'md'} pt={'sm'}>
          Recent Playlists
        </Text>

        {playlists?.models.length === 0 && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no playlists yet to display
          </Text>
        )}

        <ScrollArea
          h={'100%'}
          scrollbars={'y'}
          scrollbarSize={7}
          styles={{
            root: {
              height: '100%'
            },
            viewport: {
              '> div': {
                height: 0,
                minHeight: '100%',
                minWidth: '100%',
                width: 0
              }
            }
          }}
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
                linear-gradient(to top, transparent 99%, ${theme.white}),
                linear-gradient(to bottom, transparent 95%, ${theme.white})
              `
            }
          })}
        >
          <SimpleGrid cols={2} px={'md'} py={'xs'}>
            {isLoading || !playlists ? (
              <Loader />
            ) : (
              playlists.models.map((playlist) => (
                <LocalPlaylistCard key={playlist.id} playlist={playlist} />
              ))
            )}
          </SimpleGrid>
        </ScrollArea>
      </Stack>
    </Card>
  )
}

export default HomeRecentPlaylists
