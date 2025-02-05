import { Group, Skeleton, Stack } from '@mantine/core'

function AlbumsLoader() {
  return (
    <Group gap={'xl'} data-testid={'albums-loader'}>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'xs'} align={'center'}>
          <Skeleton
            radius={'lg'}
            h={150}
            w={150}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Stack gap={0} align={'center'}>
            <Skeleton w={100} h={15} mb={4} />
            <Skeleton w={60} h={10} />
          </Stack>
        </Stack>
      ))}
    </Group>
  )
}

export default AlbumsLoader
