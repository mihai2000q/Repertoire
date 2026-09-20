import { Divider, Group, Skeleton, Stack } from '@mantine/core'
import PlaylistSongsLoader from './PlaylistSongsLoader.tsx'

function PlaylistLoader() {
  return (
    <Stack px={'xl'} data-testid={'playlist-loader'}>
      <Group>
        <Skeleton radius={'lg'} w={'max(12vw, 150px)'} h={'max(12vw, 150px)'} />
        <Stack gap={'xxs'} pt={'md'}>
          <Skeleton w={80} h={15} />
          <Skeleton w={'max(25vw, 200px)'} h={'max(4vw, 48px)'} my={'xs'} />
          <Skeleton w={85} h={15} />
        </Stack>
      </Group>

      <Divider />

      <PlaylistSongsLoader />
    </Stack>
  )
}

export default PlaylistLoader
