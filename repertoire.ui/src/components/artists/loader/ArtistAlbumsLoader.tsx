import { Group, SimpleGrid, Skeleton, Space, Stack } from '@mantine/core'

function ArtistAlbumsLoader() {
  return (
    <Stack gap={0}>
      <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
        <Skeleton w={60} h={15} />
        <Skeleton w={100} h={11} />
        <Space flex={1} />
        <Skeleton radius={'sm'} w={20} h={20} />
      </Group>
      <SimpleGrid
        cols={{ sm: 1, md: 2, xl: 3 }}
        spacing={0}
        verticalSpacing={0}
        style={{ overflow: 'auto', maxHeight: '55vh' }}
      >
        {Array.from(Array(3)).map((_, i) => (
          <Group key={i} align={'center'} px={'md'} py={'xs'} wrap={'nowrap'}>
            <Skeleton radius={'md'} w={40} h={40} />
            <Stack gap={4}>
              <Skeleton w={100} h={14} />
              <Skeleton w={50} h={8} />
            </Stack>
          </Group>
        ))}
      </SimpleGrid>
    </Stack>
  )
}

export default ArtistAlbumsLoader
