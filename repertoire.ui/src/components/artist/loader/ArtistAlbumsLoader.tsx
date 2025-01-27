import { Card, Group, SimpleGrid, Skeleton, Stack } from '@mantine/core'

function ArtistAlbumsLoader() {
  return (
    <Card variant={'panel'} data-testid={'albums-loader'} p={0} h={'100%'} mb={'xs'}>
      <Stack gap={0}>
        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Skeleton w={60} h={15} />
          <Skeleton w={100} h={11} />
        </Group>
        <SimpleGrid cols={{ sm: 1, md: 2, xl: 3 }} spacing={0} verticalSpacing={0}>
          {Array.from(Array(8)).map((_, i) => (
            <Group key={i} px={'md'} py={'xs'} wrap={'nowrap'}>
              <Skeleton radius={'md'} w={40} h={40} />
              <Stack gap={4}>
                <Skeleton w={100} h={14} />
                <Skeleton w={50} h={8} />
              </Stack>
            </Group>
          ))}
        </SimpleGrid>
      </Stack>
    </Card>
  )
}

export default ArtistAlbumsLoader
