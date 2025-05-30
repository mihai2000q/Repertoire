import { Skeleton, Stack } from '@mantine/core'

function ArtistsLoader() {
  return Array.from(Array(40)).map((_, i) => (
    <Stack key={i} gap={'xs'} align={'center'} pb={'md'}>
      <Skeleton
        radius={'50%'}
        pb={'100%'}
        w={'100%'}
        style={(theme) => ({ boxShadow: theme.shadows.sm })}
      />
      <Skeleton w={90} h={15} mb={4} />
    </Stack>
  ))
}

export default ArtistsLoader
