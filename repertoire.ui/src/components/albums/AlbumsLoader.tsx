import { Skeleton, Stack } from '@mantine/core'

function AlbumsLoader() {
  return Array.from(Array(40)).map((_, i) => (
    <Stack key={i} gap={'xs'} align={'center'} pb={'md'}>
      <Skeleton
        radius={'lg'}
        pt={'100%'}
        w={'100%'}
        style={(theme) => ({ boxShadow: theme.shadows.sm })}
      />
      <Stack gap={0} align={'center'}>
        <Skeleton w={{ base: 125, xs: 100, sm: 90, md: 105, xl: 120 }} h={15} mb={4} />
        <Skeleton w={{ base: 68, xs: 50, sm: 43, md: 60, xl: 70  }} h={10} />
      </Stack>
    </Stack>
  ))
}

export default AlbumsLoader
