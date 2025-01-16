import { Card, Group, Skeleton, Stack } from '@mantine/core'

function ArtistSongsLoader() {
  return (
    <Card variant={'panel'} data-testid={'songs-loader'} p={0} h={'100%'} mb={'xs'}>
      <Stack gap={0}>
        <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
          <Skeleton w={60} h={15} />
          <Skeleton w={100} h={11} />
        </Group>
        <Stack gap={0}>
          {Array.from(Array(8)).map((_, i) => (
            <Group key={i} align={'center'} px={'md'} py={'xs'}>
              <Skeleton radius={'md'} w={40} h={40} />
              <Stack gap={4}>
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
