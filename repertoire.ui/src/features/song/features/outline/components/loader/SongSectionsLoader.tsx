import { Group, Skeleton, Stack } from '@mantine/core'

function SongSectionsLoader() {
  return (
    <Stack aria-label={'song-sections-loader'} gap={0}>
      {Array.from({ length: 4 }).map((_, i) => (
        <Group key={i} px={'lg'} py={'md'} gap={'xs'}>
          <Skeleton h={15} w={16} mr={'xs'} />
          <Stack flex={1} gap={'xs'}>
            <Skeleton w={'40%'} h={15} />
            <Skeleton w={'13%'} h={13} />
          </Stack>

          <Skeleton radius={'50%'} h={30} w={30} />
          <Skeleton radius={'50%'} h={30} w={30} />
        </Group>
      ))}
    </Stack>
  )
}

export default SongSectionsLoader
