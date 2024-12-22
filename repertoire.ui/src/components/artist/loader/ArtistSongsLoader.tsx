import { Group, Skeleton, Space, Stack } from '@mantine/core'

function ArtistSongsLoader() {
  return (
    <Stack gap={0}>
      <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
        <Skeleton w={60} h={15} />
        <Skeleton w={100} h={11} />
        <Space flex={1} />
        <Skeleton radius={'sm'} w={20} h={20} />
      </Group>
        {Array.from(Array(5)).map((_, i) => (
      <Stack gap={0}>
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
  )
}

export default ArtistSongsLoader
