import { Group, Skeleton, Space, Stack } from '@mantine/core'

function SongPartsLoader() {
  return (
    <Stack aria-label={'song-parts-loader'} gap={0}>
      {Array.from({ length: 4 }).map((_, i) => (
        <Group key={i} p={'md'} gap={'xs'}>
          <Skeleton h={26} w={26} />
          <Skeleton radius={'50%'} h={30} w={30} />
          <Skeleton radius={'50%'} h={30} w={30} />
          <Skeleton w={'35%'} h={15} />
          <Space flex={1} />
          <Skeleton h={26} w={26} />
          <Skeleton h={26} w={26} />
        </Group>
      ))}
    </Stack>
  )
}

export default SongPartsLoader
