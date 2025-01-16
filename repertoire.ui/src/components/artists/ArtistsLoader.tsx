import { Group, Skeleton, Stack } from '@mantine/core'

function ArtistsLoader() {
  return (
    <Group gap={'xl'} data-testid={'artists-loader'}>
      {Array.from(Array(20)).map((_, i) => (
        <Stack key={i} gap={'xs'} align={'center'}>
          <Skeleton
            radius={'50%'}
            h={125}
            w={125}
            style={(theme) => ({ boxShadow: theme.shadows.sm })}
          />
          <Skeleton w={90} h={15} mb={4} />
        </Stack>
      ))}
    </Group>
  )
}

export default ArtistsLoader
