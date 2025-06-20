import { Card, Group, Skeleton, Stack } from '@mantine/core'

function ArtistSongsLoader() {
  return (
    <Card variant={'panel'} data-testid={'songs-loader'} p={0}>
      <Stack gap={0}>
        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Skeleton w={60} h={15} />
          <Skeleton w={100} h={11} />
        </Group>
        <Stack gap={0} h={'100%'}>
          {Array.from(Array(8)).map((_, i) => (
            <Group key={i} px={'md'} py={'xs'}>
              <Skeleton radius={'md'} w={40} h={40} />
              <Stack gap={'xxs'}>
                <Skeleton w={150} h={14} />
                <Skeleton w={70} h={8} />
              </Stack>
            </Group>
          ))}
        </Stack>
      </Stack>
    </Card>
  )
}

export default ArtistSongsLoader
