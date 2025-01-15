import { Group, Skeleton, Stack } from '@mantine/core'

function PlaylistsLoader() {
  return (
    <Group data-testid={'playlists-loader'}>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'xs'} align={'center'}>
          <Skeleton
            radius={'lg'}
            h={150}
            w={150}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Skeleton w={100} h={15} mb={4} />
        </Stack>
      ))}
    </Group>
  )
}

export default PlaylistsLoader
