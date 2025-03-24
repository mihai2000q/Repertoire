import {
  AspectRatio,
  Card,
  CardProps,
  Group,
  Image,
  ScrollArea,
  SimpleGrid,
  Skeleton,
  Stack,
  Text
} from '@mantine/core'
import { useGetPlaylistsQuery } from '../../state/api/playlistsApi.ts'
import playlistPlaceholder from '../../assets/image-placeholder-1.jpg'
import Playlist from '../../types/models/Playlist.ts'

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
  return (
    <Group wrap={'nowrap'}>
      <AspectRatio>
        <Image
          w={60}
          src={playlist.imageUrl}
          alt={playlist.title}
          fallbackSrc={playlistPlaceholder}
          radius={'lg'}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.2s',
            boxShadow: theme.shadows.sm,
            '&:hover': {
              boxShadow: theme.shadows.xl,
              transform: 'scale(1.1)'
            }
          })}
        />
      </AspectRatio>

      <Text fw={500} lineClamp={2}>
        {playlist.title}
      </Text>
    </Group>
  )
}

function HomePlaylists({ ...others }: CardProps) {
  const { data: playlists, isLoading } = useGetPlaylistsQuery({
    pageSize: 20,
    currentPage: 1,
    orderBy: ['updated_at desc']
  })

  return (
    <Card aria-label={'playlists'} variant={'panel'} {...others} p={0}>
      <Stack h={'100%'} gap={'xs'}>
        <Text c={'gray.7'} fz={'lg'} fw={800} px={'md'} pt={'sm'}>
          Playlists
        </Text>

        {playlists?.models.length === 0 && (
          <Text ta={'center'} c={'gray.6'} fw={500} pt={'lg'}>
            There are no playlists yet to display
          </Text>
        )}

        <ScrollArea h={'100%'} scrollbars={'y'} scrollbarSize={7}>
          {/*DO NOT Change the Max Height, it helps with the responsive layout (somehow for some reason)*/}
          {/*Also the value 100 is randomly chosen, it has no effect whatsoever*/}
          <SimpleGrid cols={2} px={'md'} pt={'xs'} mah={100}>
            {(isLoading || !playlists) ? (
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

export default HomePlaylists
