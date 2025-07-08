import { Card, Group, Skeleton, Stack } from '@mantine/core'

function PlaylistSongsLoader() {
  return (
    <Card variant={'panel'} p={0} mx={'xs'} mb={'lg'}>
      <Stack gap={0}>
        <Group px={'md'} pt={'md'} pb={'xs'} gap={'xs'}>
          <Skeleton w={50} h={15} />
          <Skeleton w={100} h={12} />
        </Group>
        <Stack gap={0} px={'md'} py={'xs'}>
          {Array.from(Array(5)).map((_, i) => (
            <Group key={i} px={'md'} py={'xs'}>
              <Skeleton radius={'sm'} w={15} h={16} mr={'xxs'} />
              <Skeleton radius={'8px'} w={38} h={38} />
              <Skeleton maw={170} miw={100} h={15} />
            </Group>
          ))}
        </Stack>
      </Stack>
    </Card>
  )
}

export default PlaylistSongsLoader
