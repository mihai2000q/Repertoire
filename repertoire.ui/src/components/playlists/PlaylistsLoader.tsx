import { Skeleton, Stack } from '@mantine/core'

function PlaylistsLoader() {
  return Array.from(Array(40)).map((_, i) => (
    <Stack key={i} gap={'xs'} align={'center'}>
      <Skeleton
        radius={'lg'}
        w={'100%'}
        pb={'100%'}
        style={(theme) => ({ boxShadow: theme.shadows.sm })}
      />
      <Skeleton w={100} h={15} mb={4} />
    </Stack>
  ))
}

export default PlaylistsLoader
