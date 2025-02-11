import { Divider, Group, SimpleGrid, Skeleton, Stack } from '@mantine/core'

function ArtistDrawerLoader() {
  return (
    <Stack gap={'xs'} data-testid={'artist-drawer-loader'}>
      <Skeleton radius={0} w={'100%'} h={330} />

      <Stack px={'md'} pb={'md'} gap={6}>
        <Skeleton w={200} h={25} />
        <Group ml={2} gap={'xxs'}>
          <Skeleton w={55} h={15} />
          <Skeleton w={45} h={15} />
        </Group>

        <Stack gap={2} my={6}>
          <Skeleton ml={2} w={50} h={12} />
          <Divider />
        </Stack>

        <SimpleGrid cols={2} px={'xs'}>
          {Array.from({ length: 4 }).map((_, i) => (
            <Group key={i} wrap={'nowrap'} gap={'xs'}>
              <Skeleton radius={'8px'} w={28} h={28} />
              <Stack gap={2}>
                <Skeleton w={90} h={15} />
                <Skeleton w={60} h={10} />
              </Stack>
            </Group>
          ))}
        </SimpleGrid>

        <Stack gap={2} my={6}>
          <Skeleton ml={2} w={38} h={12} />
          <Divider />
        </Stack>

        <SimpleGrid cols={2} px={'xs'}>
          {Array.from({ length: 6 }).map((_, i) => (
            <Group key={i} wrap={'nowrap'} gap={'xs'}>
              <Skeleton radius={'8px'} w={28} h={28} />
              <Stack gap={2}>
                <Skeleton w={80} h={15} />
                <Skeleton w={60} h={8} />
              </Stack>
            </Group>
          ))}
        </SimpleGrid>
      </Stack>
    </Stack>
  )
}

export default ArtistDrawerLoader
